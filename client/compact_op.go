// Copyright (c) 2020 Souliot
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in a
// copies or substantial portions of the Software.

package client

import (
	pb "github.com/souliot/bolt/mvccpb"
)

// CompactOp represents a compact operation.
type CompactOp struct {
	revision int64
	physical bool
}

// CompactOption configures compact operation.
type CompactOption func(*CompactOp)

func (op *CompactOp) applyCompactOpts(opts []CompactOption) {
	for _, opt := range opts {
		opt(op)
	}
}

// OpCompact wraps slice CompactOption to create a CompactOp.
func OpCompact(rev int64, opts ...CompactOption) CompactOp {
	ret := CompactOp{revision: rev}
	ret.applyCompactOpts(opts)
	return ret
}

func (op CompactOp) toRequest() *pb.CompactionRequest {
	return &pb.CompactionRequest{Revision: op.revision, Physical: op.physical}
}

// WithCompactPhysical makes Compact wait until all compacted entries are
// removed from the etcd server's storage.
func WithCompactPhysical() CompactOption {
	return func(op *CompactOp) { op.physical = true }
}
