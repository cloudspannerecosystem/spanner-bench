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
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v2"
)

const userAgent = "spannerbench/0.1"

var (
	config string
	n      int // number of iterations for each
)

func main() {
	ctx := context.Background()
	flag.StringVar(&config, "f", "benchmark.yaml", "")
	flag.IntVar(&n, "n", 50, "")
	flag.Usage = func() {
		fmt.Println(usageText)
	}
	flag.Parse()

	data, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatalf("Failed to read the config file: %v", err)
	}

	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		log.Fatalf("Cannot parse the config file: %v", err)
	}

	client, err := spanner.NewClient(ctx, c.Database, option.WithUserAgent(userAgent))
	if err != nil {
		log.Fatalf("Cannot create Spanner client: %v", err)
	}

	b := benchmarks{
		client:     client,
		n:          n,
		benchmarks: c.Benchmarks,
	}
	b.start()
}

const usageText = `spannerbench [options...]

Options:
-f   Config file to read from, by default "benchmark.yaml". 
-n   Number of times to run a query, by default 20.`
