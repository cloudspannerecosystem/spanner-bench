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
	"time"

	"cloud.google.com/go/spanner"
	"github.com/rakyll/spannerbench"
)

func main() {
	spannerbench.Benchmark(
		"projects/YOUR_PROJECT/instances/YOUR_INSTANCE/databases/YOUR_DB",
		BenchmarkReadOnly,
		Benchmark,
	)
}

func BenchmarkReadOnly(b *spannerbench.B) {
	b.N(50) // Runs for 100 times.
	b.MaxStaleness(500 * time.Millisecond)
	b.RunReadOnly(func(tx *spanner.ReadOnlyTransaction) error {
		return nil
	})
}

func Benchmark(b *spannerbench.B) {
	b.N(50) // Runs for 100 times.
	b.Run(func(tx *spanner.ReadWriteTransaction) error {
		return nil
	})
}
