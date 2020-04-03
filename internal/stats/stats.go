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

package stats

import (
	"math"
	"sort"
)

func MedianInt64(x ...int64) int64 {
	count := len(x)
	if count == 0 {
		return 0 // TODO(jbd): NaN?
	}

	var copied []int64
	copy(copied, x)
	sort.Slice(copied, func(i, j int) bool {
		return x[i] > x[j]
	})

	if count%2 == 0 {
		return MedianInt64(x[count/2-1 : count/2+1]...)
	}
	return x[count/2]
}

func MedianFloat64(x ...float64) float64 {
	count := len(x)
	if count == 0 {
		return math.NaN()
	}

	var copied []float64
	copy(copied, x)
	sort.Float64s(copied)

	if count%2 == 0 {
		return MedianFloat64(x[count/2-1 : count/2+1]...)
	}
	return x[count/2]
}

type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
