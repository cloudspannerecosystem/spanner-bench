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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/fatih/color"
	"google.golang.org/api/iterator"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"gopkg.in/yaml.v2"
)

var (
	config string
	n      int // number of iterations for each
	// TODO(jbd): Allow concurrent runs.
)

func main() {
	ctx := context.Background()
	flag.StringVar(&config, "f", "benchmark.yaml", "")
	flag.IntVar(&n, "n", 20, "")
	flag.Parse()

	data, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatalf("Failed to read the config file: %v", err)
	}

	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		log.Fatalf("Cannot parse the config file: %v", err)
	}

	client, err := spanner.NewClient(ctx, c.Database)
	if err != nil {
		log.Fatalf("Cannot create Spanner client: %v", err)
	}

	b := benchmarks{
		client:  client,
		n:       n,
		queries: c.Queries,
	}
	b.start()
}

type benchmarks struct {
	client  *spanner.Client
	n       int
	queries []Query
}

func (b *benchmarks) start() {
	for _, q := range b.queries {
		b.run(q)
	}
}

func (b *benchmarks) run(q Query) {
	fmt.Println(q.Name)

	var results []benchmarkResult
	stmt := spanner.NewStatement(q.SQL)
	for _, opt := range q.Optimizers {
		results = append(results, b.queryN(opt, stmt))
	}
	b.print(q, results...)
}

func (b *benchmarks) queryN(v string, stmt spanner.Statement) benchmarkResult {
	var i int

	var total benchmarkResult
	for {
		if i == b.n {
			break
		}
		result, err := b.query(v, stmt)
		if err != nil {
			// TODO(jbd): Error if too many retries.
			continue
		}
		total.RowsScanned += result.RowsScanned
		total.RowsReturned += result.RowsReturned
		total.CPUTime += result.CPUTime
		total.QueryPlanTime += result.QueryPlanTime
		total.ElapsedTime += result.ElapsedTime
		i++
	}
	// TODO(jbd): Use the 50th percentile instead.
	return benchmarkResult{
		Optimizer:     v,
		RowsScanned:   total.RowsReturned / int64(i),
		RowsReturned:  total.RowsReturned / int64(i),
		QueryPlanTime: total.QueryPlanTime / time.Duration(i),
		CPUTime:       total.CPUTime / time.Duration(i),
		ElapsedTime:   total.ElapsedTime / time.Duration(i),
	}
}

func (b *benchmarks) query(v string, stmt spanner.Statement) (benchmarkResult, error) {
	mode := sppb.ExecuteSqlRequest_PROFILE
	it := b.client.Single().QueryWithOptions(context.Background(), stmt, spanner.QueryOptions{
		Mode: &mode,
		Options: &sppb.ExecuteSqlRequest_QueryOptions{
			OptimizerVersion: v,
		},
	})
	defer it.Stop()
	for {
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
		case "rows_scanned":
			result.RowsScanned = parseInt64(v.(string))
		case "rows_returned":
			result.RowsReturned = parseInt64(v.(string))
		case "query_plan_creation_time":
			result.QueryPlanTime = parseDuration(v.(string))
		case "cpu_time":
			result.CPUTime = parseDuration(v.(string))
		case "elapsed_time":
			result.ElapsedTime = parseDuration(v.(string))
		}
	}
	result.Optimizer = v
	return result, nil
}

func parseInt64(v string) int64 {
	parsed, _ := strconv.ParseInt(v, 10, 64)
	return parsed
}

func parseDuration(v string) time.Duration {
	parts := strings.Split(v, " ")
	if len(parts) < 1 {
		return time.Duration(0)
	}
	dur, _ := time.ParseDuration(parts[0] + "ms")
	return dur
}

func (b *benchmarks) print(q Query, r ...benchmarkResult) {
	// Print the header.
	fmt.Printf("   %10s %10s %10s \n", "(total)", "(cpu)", "(plan)")
	for i, result := range r {
		fmt.Println(result)
		if i > 0 {
			// Compare with the previous.
			fmt.Println(result.Diff(r[i-1]))
		}
	}
}

type benchmarkResult struct {
	Optimizer     string
	ElapsedTime   time.Duration
	CPUTime       time.Duration
	QueryPlanTime time.Duration
	RowsScanned   int64
	RowsReturned  int64
}

func (b benchmarkResult) String() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s: ", b.Optimizer)
	fmt.Fprintf(buf, "%10s ", b.ElapsedTime)
	fmt.Fprintf(buf, "%10s ", b.CPUTime)
	fmt.Fprintf(buf, "%10s    ", b.QueryPlanTime)
	fmt.Fprintf(buf, "%d/%d", b.RowsScanned, b.RowsReturned)
	return buf.String()
}

func (b benchmarkResult) Diff(prev benchmarkResult) string {
	buf := &strings.Builder{}
	diffElapsed := float64(prev.ElapsedTime-b.ElapsedTime) / float64(prev.ElapsedTime) * -100
	diffCPU := float64(prev.CPUTime-b.CPUTime) / float64(prev.CPUTime) * -100
	diffQuery := float64(prev.QueryPlanTime-b.QueryPlanTime) / float64(prev.QueryPlanTime) * -100

	fmt.Fprintf(buf, "  %s", formatPercentage(diffElapsed))
	fmt.Fprintf(buf, "  %s", formatPercentage(diffCPU))
	fmt.Fprintf(buf, "  %s", formatPercentage(diffQuery))
	return buf.String()
}

func formatPercentage(v float64) string {
	const col = 10
	txt := fmt.Sprintf("%2.2f", v) + "%"
	if v == 0 {
		return pad(txt, 10)
	}
	if v > 0 {
		txt = "+" + txt
		return color.RedString(pad(txt, 10))
	}
	return color.GreenString(pad(txt, 10))
}

func pad(v string, col int) string {
	if len(v) >= col {
		return v
	}
	return strings.Repeat(" ", col-len(v)) + v
}
