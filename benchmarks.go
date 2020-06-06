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

package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/rakyll/spannerbench/internal/stats"
	"google.golang.org/api/iterator"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
)

type benchmarks struct {
	client     *spanner.Client
	n          int
	benchmarks []Benchmark
}

func (b *benchmarks) start() {
	for _, bench := range b.benchmarks {
		b.run(bench)
	}
}

func (b *benchmarks) run(bench Benchmark) {
	fmt.Println(bench.Name)

	var fn func() (benchmarkResult, error)
	if bench.ReadOnly {
		fn = b.makeReadOnly(bench)
	} else {
		fn = b.makeReadWrite(bench)
	}
	result := b.runN(fn)

	fmt.Println(result) // TODO(jbd): Allow other types of output.
}

func (b *benchmarks) makeReadOnly(bench Benchmark) func() (benchmarkResult, error) {
	// TODO(jbd): Add staleness options.
	stmts := parseSQL(bench.SQL)

	return func() (benchmarkResult, error) {
		ctx := context.Background()
		start := time.Now()

		tx := b.client.ReadOnlyTransaction()
		defer tx.Close()

		for _, stmt := range stmts {
			it := tx.QueryWithOptions(ctx, stmt, spanner.QueryOptions{
				Options: &sppb.ExecuteSqlRequest_QueryOptions{
					OptimizerVersion: bench.Optimizer,
				},
			})
			// TODO(jbd): Use server-side elapsed time, CPU time, etc.
			// Should we loop over the results?
			_, err := it.Next()
			if err != iterator.Done {
				return benchmarkResult{}, err
			}
		}

		return benchmarkResult{
			Elapsed: time.Now().Sub(start),
		}, nil
	}
}

func (b *benchmarks) makeReadWrite(bench Benchmark) func() (benchmarkResult, error) {
	ctx := context.Background()
	start := time.Now()

	stmts := parseSQL(bench.SQL)

	return func() (benchmarkResult, error) {
		_, err := b.client.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
			for _, stmt := range stmts {
				_ = tx.QueryWithOptions(ctx, stmt, spanner.QueryOptions{
					Options: &sppb.ExecuteSqlRequest_QueryOptions{
						OptimizerVersion: bench.Optimizer,
					},
				})
			}
			return nil
		})
		if err != nil {
			return benchmarkResult{}, err
		}
		return benchmarkResult{
			Elapsed: time.Now().Sub(start),
		}, nil
	}
}

func (b *benchmarks) runN(f func() (benchmarkResult, error)) benchmarkResult {
	var i, retries int
	var elapsed []int64

	for {
		if i == b.n {
			break
		}
		result, err := f()
		retries++
		if err != nil {
			if retries > 2*b.n {
				log.Fatalf("Query failed too many times: %v\n", err)
			}
			continue
		}
		elapsed = append(elapsed, int64(result.Elapsed))
		i++
	}

	return benchmarkResult{
		Elapsed: time.Duration(stats.MedianInt64(elapsed...)),
	}
}

func parseSQL(sql string) []spanner.Statement {
	var statements []spanner.Statement

	stmts := strings.Split(sql, ";")
	for _, stmt := range stmts {
		if stmt == "" {
			continue
		}
		statements = append(statements, spanner.NewStatement(stmt))
	}
	return statements
}
