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

package backend

import (
	"crypto/rand"
	"os"
	"testing"
	"time"
)

func BenchmarkBackendPut(b *testing.B) {
	backend, tmppath := NewTmpBackend(100*time.Millisecond, 10000)
	defer backend.Close()
	defer os.Remove(tmppath)

	// prepare keys
	keys := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = make([]byte, 64)
		rand.Read(keys[i])
	}
	value := make([]byte, 128)
	rand.Read(value)

	batchTx := backend.BatchTx()

	batchTx.Lock()
	batchTx.UnsafeCreateBucket([]byte("test"))
	batchTx.Unlock()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		batchTx.Lock()
		batchTx.UnsafePut([]byte("test"), keys[i], value)
		batchTx.Unlock()
	}
}
