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

import "github.com/prometheus/client_golang/prometheus"

var (
	commitSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "etcd",
		Subsystem: "disk",
		Name:      "backend_commit_duration_seconds",
		Help:      "The latency distributions of commit called by backend.",

		// lowest bucket start of upper bound 0.001 sec (1 ms) with factor 2
		// highest bucket start of 0.001 sec * 2^13 == 8.192 sec
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 14),
	})

	rebalanceSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "etcd_debugging",
		Subsystem: "disk",
		Name:      "backend_commit_rebalance_duration_seconds",
		Help:      "The latency distributions of commit.rebalance called by bboltdb backend.",

		// lowest bucket start of upper bound 0.001 sec (1 ms) with factor 2
		// highest bucket start of 0.001 sec * 2^13 == 8.192 sec
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 14),
	})

	spillSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "etcd_debugging",
		Subsystem: "disk",
		Name:      "backend_commit_spill_duration_seconds",
		Help:      "The latency distributions of commit.spill called by bboltdb backend.",

		// lowest bucket start of upper bound 0.001 sec (1 ms) with factor 2
		// highest bucket start of 0.001 sec * 2^13 == 8.192 sec
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 14),
	})

	writeSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "etcd_debugging",
		Subsystem: "disk",
		Name:      "backend_commit_write_duration_seconds",
		Help:      "The latency distributions of commit.write called by bboltdb backend.",

		// lowest bucket start of upper bound 0.001 sec (1 ms) with factor 2
		// highest bucket start of 0.001 sec * 2^13 == 8.192 sec
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 14),
	})

	defragSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "etcd",
		Subsystem: "disk",
		Name:      "backend_defrag_duration_seconds",
		Help:      "The latency distribution of backend defragmentation.",

		// 100 MB usually takes 1 sec, so start with 10 MB of 100 ms
		// lowest bucket start of upper bound 0.1 sec (100 ms) with factor 2
		// highest bucket start of 0.1 sec * 2^12 == 409.6 sec
		Buckets: prometheus.ExponentialBuckets(.1, 2, 13),
	})

	snapshotTransferSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "etcd",
		Subsystem: "disk",
		Name:      "backend_snapshot_duration_seconds",
		Help:      "The latency distribution of backend snapshots.",

		// lowest bucket start of upper bound 0.01 sec (10 ms) with factor 2
		// highest bucket start of 0.01 sec * 2^16 == 655.36 sec
		Buckets: prometheus.ExponentialBuckets(.01, 2, 17),
	})
)

func init() {
	prometheus.MustRegister(commitSec)
	prometheus.MustRegister(rebalanceSec)
	prometheus.MustRegister(spillSec)
	prometheus.MustRegister(writeSec)
	prometheus.MustRegister(defragSec)
	prometheus.MustRegister(snapshotTransferSec)
}
