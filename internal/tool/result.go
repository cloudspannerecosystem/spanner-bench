// Copyright 2020 Google LLC
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
	Elapsed          time.Duration
	CPUElapsed       time.Duration
	OptimizerElapsed time.Duration
}

func (b benchmarkResult) String() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%v %v %v", b.Elapsed, b.CPUElapsed, b.OptimizerElapsed)
	return buf.String()
}
