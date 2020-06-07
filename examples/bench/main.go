package main

import (
	"time"

	"cloud.google.com/go/spanner"
	"github.com/rakyll/spannerbench/spannerbench"
)

func main() {
	spannerbench.Benchmark(
		"projects/computelabs/instances/hello/databases/db",
		BenchmarkQuery,
	)
}

func BenchmarkLikesQuery(b *spannerbench.B) {
	b.N(100) // Runs for 100 times.
	b.MaxStaleness(500 * time.Millisecond)
	b.RunReadOnly(func(tx *spanner.ReadOnlyTransaction) error {
		return nil
	})
}

func BenchmarkQuery(b *spannerbench.B) {
	b.N(5) // Runs for 100 times.
	b.Run(func(tx *spanner.ReadWriteTransaction) error {
		return nil
	})
}
