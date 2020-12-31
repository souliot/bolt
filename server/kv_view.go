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

package server

import (
	"github.com/souliot/bolt/server/lease"
)

type readView struct{ kv KV }

func (rv *readView) FirstRev() int64 {
	tr := rv.kv.Read()
	defer tr.End()
	return tr.FirstRev()
}

func (rv *readView) Rev() int64 {
	tr := rv.kv.Read()
	defer tr.End()
	return tr.Rev()
}

func (rv *readView) Range(key, end []byte, ro RangeOptions) (r *RangeResult, err error) {
	tr := rv.kv.Read()
	defer tr.End()
	return tr.Range(key, end, ro)
}

type writeView struct{ kv KV }

func (wv *writeView) DeleteRange(key, end []byte) (n, rev int64) {
	tw := wv.kv.Write()
	defer tw.End()
	return tw.DeleteRange(key, end)
}

func (wv *writeView) Put(key, value []byte, lease lease.LeaseID) (rev int64) {
	tw := wv.kv.Write()
	defer tw.End()
	return tw.Put(key, value, lease)
}
