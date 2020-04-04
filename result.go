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
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

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
	fmt.Fprintf(buf, "%10d ", b.RowsScanned)
	fmt.Fprintf(buf, "%10s ", b.ElapsedTime)
	fmt.Fprintf(buf, "%10s ", b.CPUTime)
	fmt.Fprintf(buf, "%10s    ", b.QueryPlanTime)
	fmt.Fprintf(buf, "%d/%d", b.RowsScanned, b.RowsReturned)
	return buf.String()
}

func (b benchmarkResult) Diff(prev benchmarkResult) string {
	diffScanned := float64(prev.RowsScanned-b.RowsScanned) / float64(prev.RowsScanned) * -100
	return fmt.Sprintf("   %10s ", formatPercentage(diffScanned))
}

func formatPercentage(v float64) string {
	txt := fmt.Sprintf("%2.2f", v) + "%"
	if strings.Contains(txt, "0.00") {
		return "0.00%"
	}
	if v > 0 {
		txt = "+" + txt
		return color.RedString(txt)
	}
	return color.GreenString(txt)
}

func pad(v string, col int) string {
	if len(v) >= col {
		return v
	}
	return strings.Repeat(" ", col-len(v)) + v
}
