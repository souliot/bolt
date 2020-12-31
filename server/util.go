// Copyright 2016 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/souliot/bolt/backend"
	pb "github.com/souliot/bolt/mvccpb"
)

func UpdateConsistentIndex(be backend.Backend, index uint64) {
	tx := be.BatchTx()
	tx.Lock()
	defer tx.Unlock()

	var oldi uint64
	_, vs := tx.UnsafeRange(metaBucketName, consistentIndexKeyName, nil, 0)
	if len(vs) != 0 {
		oldi = binary.BigEndian.Uint64(vs[0])
	}

	if index <= oldi {
		return
	}

	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, index)
	tx.UnsafePut(metaBucketName, consistentIndexKeyName, bs)
}

func WriteKV(be backend.Backend, kv pb.KeyValue) {
	ibytes := newRevBytes()
	revToBytes(revision{main: kv.ModRevision}, ibytes)

	d, err := kv.Marshal()
	if err != nil {
		plog.Fatalf("cannot marshal event: %v", err)
	}

	be.BatchTx().Lock()
	be.BatchTx().UnsafePut(keyBucketName, ibytes, d)
	be.BatchTx().Unlock()
}

type notifier struct {
	c   chan struct{}
	err error
}

func newNotifier() *notifier {
	return &notifier{
		c: make(chan struct{}),
	}
}

func (nc *notifier) notify(err error) {
	nc.err = err
	close(nc.c)
}

func warnOfExpensiveReadOnlyTxnRequest(now time.Time, r *pb.TxnRequest, txnResponse *pb.TxnResponse, err error) {
	reqStringer := pb.NewLoggableTxnRequest(r)
	var resp string
	if !isNil(txnResponse) {
		var resps []string
		for _, r := range txnResponse.Responses {
			switch op := r.Response.(type) {
			case *pb.ResponseOp_ResponseRange:
				resps = append(resps, fmt.Sprintf("range_response_count:%d", len(op.ResponseRange.Kvs)))
			default:
				// only range responses should be in a read only txn request
			}
		}
		resp = fmt.Sprintf("responses:<%s> size:%d", strings.Join(resps, " "), proto.Size(txnResponse))
	}
	warnOfExpensiveGenericRequest(now, reqStringer, "read-only range ", resp, err)
}

func warnOfExpensiveReadOnlyRangeRequest(now time.Time, reqStringer fmt.Stringer, rangeResponse *pb.RangeResponse, err error) {
	var resp string
	if !isNil(rangeResponse) {
		resp = fmt.Sprintf("range_response_count:%d size:%d", len(rangeResponse.Kvs), proto.Size(rangeResponse))
	}
	warnOfExpensiveGenericRequest(now, reqStringer, "read-only range ", resp, err)
}

func warnOfExpensiveGenericRequest(now time.Time, reqStringer fmt.Stringer, prefix string, resp string, err error) {
	// TODO: add metrics
	d := time.Since(now)
	if d > warnApplyDuration {
		var result string
		if err != nil {
			result = fmt.Sprintf("error:%v", err)
		} else {
			result = resp
		}
		plog.Warningf("%srequest %q with result %q took too long (%v) to execute", prefix, reqStringer.String(), result, d)
	}
}

func isNil(msg proto.Message) bool {
	return msg == nil || reflect.ValueOf(msg).IsNil()
}
