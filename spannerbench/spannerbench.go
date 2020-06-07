// Copyright 2020 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package spannerbench provides a benchmarking framework
// for Google Cloud Spanner.
package spannerbench

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/rakyll/spannerbench/internal/histogram"
	"google.golang.org/api/option"
)

// B represents a benchmark.
// Use Benchmark function to run benchmarks.
type B struct {
	client    *spanner.Client
	staleness *spanner.TimestampBound
	n         int

	elapsed []int64
}

// MaxStaleness sets the max staleness in reads
// It will be ignored for read-write transactions.
func (b *B) MaxStaleness(d time.Duration) {
	tb := spanner.MaxStaleness(d)
	b.staleness = &tb
}

// ExactStaleness represents the exact staleness
// in reads. It will be ignored for read-write transactions.
func (b *B) ExactStaleness(d time.Duration) {
	tb := spanner.ExactStaleness(d)
	b.staleness = &tb
}

// N sets the number of times a benchmarks will be run.
// If not set, default value (20) is used.
func (b *B) N(n int) {
	b.n = n
}

// TODO(jbd): Allow users to set concurrency.

// RunReadOnly runs readonly transaction benchmarks.
// It starts a read-only transaction and calls fn.
// The benchmark will be repeated for a number of times
// and results will be printed.
//
// Run is not safe for concurrent usage. Don't reuse this
// benchmark once you call RunReadOnly.
func (b *B) RunReadOnly(fn func(tx *spanner.ReadOnlyTransaction) error) {
	// TODO(jbd): Cleanup after running.
	var i, retries int

	n := b.numberOfRuns()
	for {
		if i == n {
			break
		}
		err := b.startAndRunReadOnly(fn)
		retries++
		if err != nil {
			if retries > 2*n {
				log.Fatalf("Query failed too many times: %v\n", err)
			}
			continue
		}
		i++
	}
	b.print()
}

func (b *B) startAndRunReadOnly(fn func(tx *spanner.ReadOnlyTransaction) error) error {
	start := time.Now()
	defer func() {
		dur := time.Now().Sub(start)
		b.elapsed = append(b.elapsed, int64(dur))
	}()

	tx := b.client.ReadOnlyTransaction()
	defer tx.Close()

	if b.staleness != nil {
		tx = tx.WithTimestampBound(*b.staleness)
	}
	return fn(tx)
}

// Run runs read-write transaction benchmarks.
// It starts a read-write transaction and calls fn.
// The benchmark will be repeated for a number of times
// and results will be printed.
//
// Run is not safe for concurrent usage. Don't reuse this
// benchmark once you call Run.
func (b *B) Run(fn func(tx *spanner.ReadWriteTransaction) error) {
	// TODO(jbd): Remove duplication by merging Run and RunReadOnly.
	var i, retries int

	n := b.numberOfRuns()
	for {
		if i == n {
			break
		}
		err := b.startAndRun(fn)
		retries++
		if err != nil {
			if retries > 2*n {
				log.Fatalf("Query failed too many times: %v\n", err)
			}
			continue
		}
		i++
	}
	b.print()
}

func (b *B) startAndRun(fn func(tx *spanner.ReadWriteTransaction) error) error {
	start := time.Now()
	defer func() {
		dur := time.Now().Sub(start)
		b.elapsed = append(b.elapsed, int64(dur))
	}()

	ctx := context.Background() // TODO(jbd): Consider adding context to the APIs.
	_, err := b.client.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		return fn(tx)
	})
	return err
}

func (b *B) numberOfRuns() int {
	if b.n == 0 {
		return defaultN
	}
	return b.n
}

func (b *B) print() {
	if histogram := histogram.NewHistogram(b.elapsed); histogram != nil {
		fmt.Println("Latency histogram:")
		fmt.Println(histogram)
	}
}

// Benchmark starts the benchmarks.
// Provide the full-identifier of the Google Cloud Spanner
// database as db.
func Benchmark(db string, fn ...func(b *B)) {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db, option.WithUserAgent(userAgent))
	if err != nil {
		log.Fatalf("Cannot create Spanner client: %v", err)
	}

	for _, f := range fn {
		name := funcName(f)
		fmt.Println(name)
		f(&B{
			client: client,
		}) // TODO(jbd): Fill the B.
	}
}

func funcName(fn func(b *B)) string {
	fullname := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	name := strings.Split(fullname, ".")
	if len(name) == 2 {
		return name[1]
	}
	return fullname // Anonymous functions.
}

const (
	userAgent = "spannerbench"
	defaultN  = 20
)
