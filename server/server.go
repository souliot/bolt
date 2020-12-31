package server

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/souliot/bolt/backend"
	pb "github.com/souliot/bolt/mvccpb"
	"github.com/souliot/bolt/pkg/fileutil"
	"github.com/souliot/bolt/pkg/idutil"
	"github.com/souliot/bolt/server/lease"
)

var (
	defaultID = 1
)

type Server struct {
	appliedIndex uint64
	Cfg          ServerConfig
	bemu         sync.Mutex
	be           backend.Backend
	lessor       lease.Lessor
	kv           ConsistentWatchableKV
	consistIndex consistentIndex
	applyBase    applier
	readMu       sync.RWMutex
	// read routine notifies etcd server that it waits for reading by sending an empty struct to
	// readwaitC
	readwaitc chan struct{}
	// readNotifier is used to notify the read routine that it can process the request
	// when there is no error
	readNotifier *notifier

	done     chan struct{}
	reqIDGen *idutil.Generator
}

var _ server = new(Server)
var _ Lessor = new(Server)

func NewServer(cfg ServerConfig) (srv *Server, err error) {
	srv = new(Server)
	srv.Cfg = cfg
	if err = fileutil.TouchDirAll(cfg.SnapDir()); err != nil {
		plog.Fatalf("create snapshot directory error: %v", err)
	}
	// bepath := cfg.backendPath()
	// beExist := fileutil.Exist(bepath)
	srv.be = openBackend(cfg)
	minTTL := time.Duration(4000) * time.Millisecond
	srv.lessor = lease.NewLessor(srv.be, int64(math.Ceil(minTTL.Seconds())))
	srv.kv = New(srv.be, srv.lessor, &srv.consistIndex)
	srv.consistIndex.setConsistentIndex(srv.kv.ConsistentIndex())
	srv.applyBase = srv.newApplierBackend()
	srv.readwaitc = make(chan struct{}, 1)
	srv.done = make(chan struct{})
	srv.readNotifier = newNotifier()
	srv.reqIDGen = idutil.NewGenerator(uint16(defaultID), time.Now())
	return
}

func (s *Server) Status() (ds backend.Status) {
	return s.be.Status()
}

func (s *Server) KV() ConsistentWatchableKV { return s.kv }

func (s *Server) Range(ctx context.Context, r *pb.RangeRequest) (resp *pb.RangeResponse, err error) {
	defer func(start time.Time) {
		warnOfExpensiveReadOnlyRangeRequest(start, r, resp, err)
	}(time.Now())

	return s.applyBase.Range(nil, r)
}

func (s *Server) Put(ctx context.Context, r *pb.PutRequest) (resp *pb.PutResponse, err error) {
	return s.applyBase.Put(nil, r)
}

func (s *Server) DeleteRange(ctx context.Context, r *pb.DeleteRangeRequest) (resp *pb.DeleteRangeResponse, err error) {
	return s.applyBase.DeleteRange(nil, r)
}

func (s *Server) Txn(ctx context.Context, r *pb.TxnRequest) (resp *pb.TxnResponse, err error) {
	if isTxnReadonly(r) {
		defer func(start time.Time) {
			warnOfExpensiveReadOnlyTxnRequest(start, r, resp, err)
		}(time.Now())
		resp, err = s.applyBase.Txn(r)
		return
	}
	return s.applyBase.Txn(r)
}

func isTxnReadonly(r *pb.TxnRequest) bool {
	for _, u := range r.Success {
		if r := u.GetRequestRange(); r == nil {
			return false
		}
	}
	for _, u := range r.Failure {
		if r := u.GetRequestRange(); r == nil {
			return false
		}
	}
	return true
}

func (s *Server) Compact(ctx context.Context, r *pb.CompactionRequest) (resp *pb.CompactionResponse, err error) {
	resp, ch, err := s.applyBase.Compaction(r)
	if ch != nil {
		<-ch
		s.be.ForceCommit()
	}
	if err != nil {
		return nil, err
	}
	if resp.Header == nil {
		resp.Header = &pb.ResponseHeader{}
	}
	resp.Header.Revision = s.kv.Rev()
	return
}

func (s *Server) LeaseGrant(ctx context.Context, r *pb.LeaseGrantRequest) (resp *pb.LeaseGrantResponse, err error) {
	// no id given? choose one
	for r.ID == int64(lease.NoLease) {
		// only use positive int64 id's
		r.ID = int64(s.reqIDGen.Next() & ((1 << 63) - 1))
	}

	return s.applyBase.LeaseGrant(r)
}

func (s *Server) LeaseRevoke(ctx context.Context, r *pb.LeaseRevokeRequest) (resp *pb.LeaseRevokeResponse, err error) {
	return s.applyBase.LeaseRevoke(r)
}

func (s *Server) LeaseRenew(ctx context.Context, id lease.LeaseID) (int64, error) {
	ttl, err := s.lessor.Renew(id)
	if err == nil { // already requested to primary lessor(leader)
		return ttl, nil
	}
	if err != lease.ErrNotPrimary {
		return -1, err
	}
	return -1, ErrTimeout
}

func (s *Server) LeaseTimeToLive(ctx context.Context, r *pb.LeaseTimeToLiveRequest) (resp *pb.LeaseTimeToLiveResponse, err error) {
	// primary; timetolive directly from leader
	le := s.lessor.Lookup(lease.LeaseID(r.ID))
	if le == nil {
		return nil, lease.ErrLeaseNotFound
	}
	// TODO: fill out ResponseHeader
	resp = &pb.LeaseTimeToLiveResponse{Header: &pb.ResponseHeader{}, ID: r.ID, TTL: int64(le.Remaining().Seconds()), GrantedTTL: le.TTL()}
	if r.Keys {
		ks := le.Keys()
		kbs := make([][]byte, len(ks))
		for i := range ks {
			kbs[i] = []byte(ks[i])
		}
		resp.Keys = kbs
	}

	return
}

func (s *Server) LeaseLeases(ctx context.Context, r *pb.LeaseLeasesRequest) (resp *pb.LeaseLeasesResponse, err error) {
	ls := s.lessor.Leases()
	lss := make([]*pb.LeaseStatus, len(ls))
	for i := range ls {
		lss[i] = &pb.LeaseStatus{ID: int64(ls[i].ID)}
	}
	return &pb.LeaseLeasesResponse{Header: newHeader(s), Leases: lss}, nil
}

func (s *Server) Close() {
	s.kv.Close()
	s.lessor.Stop()
	s.be.Close()
}

func newBackend(cfg ServerConfig) backend.Backend {
	bcfg := backend.DefaultBackendConfig()
	bcfg.Path = cfg.backendPath()
	if cfg.QuotaBackendBytes > 0 && cfg.QuotaBackendBytes != DefaultQuotaBytes {
		// permit 10% excess over quota for disarm
		bcfg.MmapSize = uint64(cfg.QuotaBackendBytes + cfg.QuotaBackendBytes/10)
	}
	return backend.New(bcfg)
}

func openBackend(cfg ServerConfig) backend.Backend {
	fn := cfg.backendPath()
	beOpened := make(chan backend.Backend)
	go func() {
		beOpened <- newBackend(cfg)
	}()
	select {
	case be := <-beOpened:
		return be
	case <-time.After(10 * time.Second):
		plog.Warningf("another etcd process is using %q and holds the file lock, or loading backend file is taking >10 seconds", fn)
		plog.Warningf("waiting for it to exit before starting...")
	}
	return <-beOpened
}
