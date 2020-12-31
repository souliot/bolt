package client

import (
	"context"

	pb "github.com/souliot/bolt/mvccpb"
)

type Event pb.Event
type WatchChan <-chan WatchResponse
type Watcher interface {
	// Watch watches on a key or prefix. The watched events will be returned
	// through the returned channel. If revisions waiting to be sent over the
	// watch are compacted, then the watch will be canceled by the server, the
	// client will post a compacted error watch response, and the channel will close.
	// If the context "ctx" is canceled or timed out, returned "WatchChan" is closed,
	// and "WatchResponse" from this closed channel has zero events and nil "Err()".
	// The context "ctx" MUST be canceled, as soon as watcher is no longer being used,
	// to release the associated resources.
	//
	// If the context is "context.Background/TODO", returned "WatchChan" will
	// not be closed and block until event is triggered, except when server
	// returns a non-recoverable error (e.g. ErrCompacted).
	// For example, when context passed with "WithRequireLeader" and the
	// connected server has no leader (e.g. due to network partition),
	// error "etcdserver: no leader" (ErrNoLeader) will be returned,
	// and then "WatchChan" is closed with non-nil "Err()".
	// In order to prevent a watch stream being stuck in a partitioned node,
	// make sure to wrap context with "WithRequireLeader".
	//
	// Otherwise, as long as the context has not been canceled or timed out,
	// watch will retry on other recoverable errors forever until reconnected.
	Watch(ctx context.Context, key string, opts ...OpOption) WatchChan

	// RequestProgress requests a progress notify response be sent in all watch channels.
	RequestProgress(ctx context.Context) error

	// Close closes the watcher and cancels all watch requests.
	Close() error
}

type WatchResponse struct {
	Header pb.ResponseHeader
	Events []*Event

	// CompactRevision is the minimum revision the watcher may receive.
	CompactRevision int64

	// Canceled is used to indicate watch failure.
	// If the watch failed and the stream was about to close, before the channel is closed,
	// the channel sends a final response that has Canceled set to true with a non-nil Err().
	Canceled bool

	// Created is used to indicate the creation of the watcher.
	Created bool

	closeErr error

	// cancelReason is a reason of canceling watch
	cancelReason string
}
