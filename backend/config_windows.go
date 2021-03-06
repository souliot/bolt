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

// +build windows

package backend

import bolt "go.etcd.io/bbolt"

var boltOpenOptions *bolt.Options = nil

// setting mmap size != 0 on windows will allocate the entire
// mmap size for the file, instead of growing it. So, force 0.

func (bcfg *BackendConfig) mmapSize() int { return 0 }
