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

type CompareTarget int
type CompareResult int

const (
	CompareVersion CompareTarget = iota
	CompareCreated
	CompareModified
	CompareValue
)

type Cmp pb.Compare

func Compare(cmp Cmp, result string, v interface{}) Cmp {
	var r pb.Compare_CompareResult

	switch result {
	case "=":
		r = pb.Compare_EQUAL
	case "!=":
		r = pb.Compare_NOT_EQUAL
	case ">":
		r = pb.Compare_GREATER
	case "<":
		r = pb.Compare_LESS
	default:
		panic("Unknown result op")
	}

	cmp.Result = r
	switch cmp.Target {
	case pb.Compare_VALUE:
		val, ok := v.(string)
		if !ok {
			panic("bad compare value")
		}
		cmp.TargetUnion = &pb.Compare_Value{Value: []byte(val)}
	case pb.Compare_VERSION:
		cmp.TargetUnion = &pb.Compare_Version{Version: mustInt64(v)}
	case pb.Compare_CREATE:
		cmp.TargetUnion = &pb.Compare_CreateRevision{CreateRevision: mustInt64(v)}
	case pb.Compare_MOD:
		cmp.TargetUnion = &pb.Compare_ModRevision{ModRevision: mustInt64(v)}
	case pb.Compare_LEASE:
		cmp.TargetUnion = &pb.Compare_Lease{Lease: mustInt64orLeaseID(v)}
	default:
		panic("Unknown compare type")
	}
	return cmp
}

func Value(key string) Cmp {
	return Cmp{Key: []byte(key), Target: pb.Compare_VALUE}
}

func Version(key string) Cmp {
	return Cmp{Key: []byte(key), Target: pb.Compare_VERSION}
}

func CreateRevision(key string) Cmp {
	return Cmp{Key: []byte(key), Target: pb.Compare_CREATE}
}

func ModRevision(key string) Cmp {
	return Cmp{Key: []byte(key), Target: pb.Compare_MOD}
}

// LeaseValue compares a key's LeaseID to a value of your choosing. The empty
// LeaseID is 0, otherwise known as `NoLease`.
func LeaseValue(key string) Cmp {
	return Cmp{Key: []byte(key), Target: pb.Compare_LEASE}
}

// KeyBytes returns the byte slice holding with the comparison key.
func (cmp *Cmp) KeyBytes() []byte { return cmp.Key }

// WithKeyBytes sets the byte slice for the comparison key.
func (cmp *Cmp) WithKeyBytes(key []byte) { cmp.Key = key }

// ValueBytes returns the byte slice holding the comparison value, if any.
func (cmp *Cmp) ValueBytes() []byte {
	if tu, ok := cmp.TargetUnion.(*pb.Compare_Value); ok {
		return tu.Value
	}
	return nil
}

// WithValueBytes sets the byte slice for the comparison's value.
func (cmp *Cmp) WithValueBytes(v []byte) { cmp.TargetUnion.(*pb.Compare_Value).Value = v }

// WithRange sets the comparison to scan the range [key, end).
func (cmp Cmp) WithRange(end string) Cmp {
	cmp.RangeEnd = []byte(end)
	return cmp
}

// WithPrefix sets the comparison to scan all keys prefixed by the key.
func (cmp Cmp) WithPrefix() Cmp {
	cmp.RangeEnd = getPrefix(cmp.Key)
	return cmp
}

// mustInt64 panics if val isn't an int or int64. It returns an int64 otherwise.
func mustInt64(val interface{}) int64 {
	if v, ok := val.(int64); ok {
		return v
	}
	if v, ok := val.(int); ok {
		return int64(v)
	}
	panic("bad value")
}

// mustInt64orLeaseID panics if val isn't a LeaseID, int or int64. It returns an
// int64 otherwise.
func mustInt64orLeaseID(val interface{}) int64 {
	if v, ok := val.(LeaseID); ok {
		return int64(v)
	}
	return mustInt64(val)
}
