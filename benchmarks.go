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

	"cloud.google.com/go/spanner"
	"github.com/rakyll/spannerbench/internal/histogram"
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
	elapsed, _, _ := b.runN(fn)

	histogram := histogram.NewHistogram(elapsed)
	fmt.Println(histogram)
}

func (b *benchmarks) makeReadOnly(bench Benchmark) func() (benchmarkResult, error) {
	// TODO(jbd): Add staleness options.
	stmts := parseSQL(bench.SQL)
	return func() (benchmarkResult, error) {
		ctx := context.Background()
		var result benchmarkResult

		tx := b.client.ReadOnlyTransaction()
		defer tx.Close()

		mode := sppb.ExecuteSqlRequest_PROFILE
		for _, stmt := range stmts {
			it := tx.QueryWithOptions(ctx, stmt, spanner.QueryOptions{
				Mode: &mode,
				Options: &sppb.ExecuteSqlRequest_QueryOptions{
					OptimizerVersion: bench.Optimizer,
				},
			})
			r, err := parseBenchmarkResult(it)
			if err != nil {
				return benchmarkResult{}, err
			}
			result.Elapsed += r.Elapsed
			result.CPUElapsed += r.CPUElapsed
			result.OptimizerElapsed += r.OptimizerElapsed
			it.Stop()
		}
		return result, nil
	}
}

func (b *benchmarks) makeReadWrite(bench Benchmark) func() (benchmarkResult, error) {
	ctx := context.Background()

	stmts := parseSQL(bench.SQL)
	return func() (benchmarkResult, error) {
		var result benchmarkResult

		mode := sppb.ExecuteSqlRequest_PROFILE
		_, err := b.client.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
			// TODO(jbd): Handle the errors.
			for _, stmt := range stmts {
				it := tx.QueryWithOptions(ctx, stmt, spanner.QueryOptions{
					Mode: &mode,
					Options: &sppb.ExecuteSqlRequest_QueryOptions{
						OptimizerVersion: bench.Optimizer,
					},
				})
				r, err := parseBenchmarkResult(it)
				if err != nil {
					return err
				}
				result.Elapsed += r.Elapsed
				result.CPUElapsed += r.CPUElapsed
				result.OptimizerElapsed += r.OptimizerElapsed
				it.Stop()
			}
			return nil
		})
		return result, err
	}
}

func (b *benchmarks) runN(f func() (benchmarkResult, error)) (elapsed, cpu, optimizer []int64) {
	var i, retries int

	for {
		if i == b.n {
			// TODO(jbd): Stop until result is stabilized.
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
		cpu = append(cpu, int64(result.CPUElapsed))
		optimizer = append(optimizer, int64(result.OptimizerElapsed))
		i++
	}
	return elapsed, cpu, optimizer
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

func parseBenchmarkResult(it *spanner.RowIterator) (benchmarkResult, error) {
	for { // Required to be able to read the stats.
		_, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return benchmarkResult{}, err
		}
	}

	var result benchmarkResult
	for k, v := range it.QueryStats {
		switch k {
		case "query_plan_creation_time":
			result.OptimizerElapsed = parseDuration(v.(string))
		case "cpu_time":
			result.CPUElapsed = parseDuration(v.(string))
		case "elapsed_time":
			result.Elapsed = parseDuration(v.(string))
		}
	}
	return result, nil
}
