package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

var (
	project  string
	instance string
	database string

	concurrency  int
	number       int
	file         string
	transaction  string // ro, rw
	maxStaleness string // TODO(jbd): Add exact staleness.
)

func main() {
	flag.StringVar(&project, "project", "", "")
	flag.StringVar(&instance, "instance", "", "")
	flag.StringVar(&database, "database", "", "")

	flag.IntVar(&concurrency, "c", 1, "")
	flag.IntVar(&number, "n", 10, "")
	flag.StringVar(&file, "f", "", "")
	flag.StringVar(&transaction, "t", "rw", "")
	flag.Parse()

	ctx := context.Background()
	db := "projects/" + project + "/instances/" + instance + "/databases/" + database
	client, err := spanner.NewClientWithConfig(ctx, db, spanner.ClientConfig{
		NumChannels: concurrency,
	})
	if err != nil {
		log.Fatal(err)
	}

	dml, err := parseFile(file)
	if err != nil {
		log.Fatalf("Cannot read and parse the DML file: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			durs, err := benchmarkN(ctx, client, dml)
			if err != nil {
				log.Println(err) // TODO(jbd): Add more context.
			} else {
				log.Println(durs)
			}
		}()
	}
	wg.Wait()
}

func benchmarkN(ctx context.Context, client *spanner.Client, statements []spanner.Statement) ([]time.Duration, error) {
	durs := make([]time.Duration, number)
	for i := 0; i < number; i++ {
		// TODO(jbd): Try running at least number times, at most 2*number times.
		dur, err := benchmark(ctx, client, statements)
		if err != nil {
			log.Println(err) // TODO(jbd): Return error if errored too many times.
		}
		durs[i] = dur
	}
	return durs, nil
}

func benchmark(ctx context.Context, client *spanner.Client, statements []spanner.Statement) (time.Duration, error) {
	start := time.Now()

	var err error
	switch transaction {
	case "ro":
		err = benchmarkReadOnly(ctx, client, statements)
	case "rw":
		err = benchmarkReadWrite(ctx, client, statements)
	}
	return time.Now().Sub(start), err
}

func benchmarkReadOnly(ctx context.Context, client *spanner.Client, statements []spanner.Statement) error {
	return nil
}
func benchmarkReadWrite(ctx context.Context, client *spanner.Client, statements []spanner.Statement) error {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, rw *spanner.ReadWriteTransaction) error {
		for _, stmt := range statements {
			it := rw.Query(ctx, stmt)
			_, err := it.Next()
			if err != iterator.Done {
				return err
			}
		}
		return nil
	})
	return err
}

func parseFile(filename string) ([]spanner.Statement, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(f)

	var statements []spanner.Statement
	for {
		l, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		// TODO(jbd): Filter out empty lines.
		stmt := spanner.NewStatement(string(l))
		statements = append(statements, stmt)
	}
	return statements, nil
}
