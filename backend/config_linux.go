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
	"syscall"

	bolt "go.etcd.io/bbolt"
)

// syscall.MAP_POPULATE on linux 2.6.23+ does sequential read-ahead
// which can speed up entire-database read with boltdb. We want to
// enable MAP_POPULATE for faster key-value store recovery in storage
// package. If your kernel version is lower than 2.6.23
// (https://github.com/torvalds/linux/releases/tag/v2.6.23), mmap might
// silently ignore this flag. Please update your kernel to prevent this.
var boltOpenOptions = &bolt.Options{
	MmapFlags:      syscall.MAP_POPULATE,
	NoFreelistSync: true,
}

func (bcfg *BackendConfig) mmapSize() int { return int(bcfg.MmapSize) }
