package client

import (
	"context"

	pb "github.com/souliot/bolt/mvccpb"
)

type (
	LeaseRevokeResponse pb.LeaseRevokeResponse
	LeaseID             int64
)

// LeaseGrantResponse wraps the protobuf message LeaseGrantResponse.
type LeaseGrantResponse struct {
	*pb.ResponseHeader
	ID    LeaseID
	TTL   int64
	Error string
}

// LeaseKeepAliveResponse wraps the protobuf message LeaseKeepAliveResponse.
type LeaseKeepAliveResponse struct {
	*pb.ResponseHeader
	ID  LeaseID
	TTL int64
}

// LeaseTimeToLiveResponse wraps the protobuf message LeaseTimeToLiveResponse.
type LeaseTimeToLiveResponse struct {
	*pb.ResponseHeader
	ID LeaseID `json:"id"`

	// TTL is the remaining TTL in seconds for the lease; the lease will expire in under TTL+1 seconds. Expired lease will return -1.
	TTL int64 `json:"ttl"`

	// GrantedTTL is the initial granted time in seconds upon lease creation/renewal.
	GrantedTTL int64 `json:"granted-ttl"`

	// Keys is the list of keys attached to this lease.
	Keys [][]byte `json:"keys"`
}

// LeaseStatus represents a lease status.
type LeaseStatus struct {
	ID LeaseID `json:"id"`
	// TODO: TTL int64
}

// LeaseLeasesResponse wraps the protobuf message LeaseLeasesResponse.
type LeaseLeasesResponse struct {
	*pb.ResponseHeader
	Leases []LeaseStatus `json:"leases"`
}

type Lease interface {
	// Grant creates a new lease.
	Grant(ctx context.Context, ttl int64) (*LeaseGrantResponse, error)

	// Revoke revokes the given lease.
	Revoke(ctx context.Context, id LeaseID) (*LeaseRevokeResponse, error)

	// TimeToLive retrieves the lease information of the given lease ID.
	TimeToLive(ctx context.Context, id LeaseID, opts ...LeaseOption) (*LeaseTimeToLiveResponse, error)

	// Leases retrieves all leases.
	Leases(ctx context.Context) (*LeaseLeasesResponse, error)

	// KeepAlive keeps the given lease alive forever. If the keepalive response
	// posted to the channel is not consumed immediately, the lease client will
	// continue sending keep alive requests to the etcd server at least every
	// second until latest response is consumed.
	//
	// The returned "LeaseKeepAliveResponse" channel closes if underlying keep
	// alive stream is interrupted in some way the client cannot handle itself;
	// given context "ctx" is canceled or timed out. "LeaseKeepAliveResponse"
	// from this closed channel is nil.
	//
	// If client keep alive loop halts with an unexpected error (e.g. "etcdserver:
	// no leader") or canceled by the caller (e.g. context.Canceled), the error
	// is returned. Otherwise, it retries.
	//
	// TODO(v4.0): post errors to last keep alive message before closing
	// (see https://github.com/coreos/etcd/pull/7866)
	KeepAlive(ctx context.Context, id LeaseID) (<-chan *LeaseKeepAliveResponse, error)

	// KeepAliveOnce renews the lease once. The response corresponds to the
	// first message from calling KeepAlive. If the response has a recoverable
	// error, KeepAliveOnce will retry the RPC with a new keep alive message.
	//
	// In most of the cases, Keepalive should be used instead of KeepAliveOnce.
	KeepAliveOnce(ctx context.Context, id LeaseID) (*LeaseKeepAliveResponse, error)

	// Close releases all resources Lease keeps for efficient communication
	// with the etcd server.
	Close() error
}
