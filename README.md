# spanner-query-benchmark

A Google Cloud Spanner query planner benchmarking tool.
It also allows you to run the same benchmark against
multiple versions of the query planner.

See `benchmark.yaml` for an example configuration.

## Installation

```
$ go get -u github.com/rakyll/spanner-query-benchmark
```

## Usage

```
$ spanner-query-benchmark
Benchmark1
    (scanned)    (total)      (cpu)     (plan)
1:        965   38.187ms    3.015ms    798.5µs    965/965
2:        965   38.271ms   2.9755ms      715µs    965/965
Benchmark2
    (scanned)    (total)      (cpu)     (plan)
1:          3  35.5635ms   1.5415ms    1.074ms    3/3
2:          3   36.135ms    1.716ms   1.2895ms    3/3
Benchmark3
    (scanned)    (total)      (cpu)     (plan)
2:        100  36.5565ms   1.7675ms   1.2715ms    100/100
```

### Output explained...

| planner version | total rows scanned | total execution time | total CPU time | planning time | rows scanned/rows returned |
|-|-|-|-|-|-|
| 1: | 965  | 38.187ms | 3.015ms | 798.5µs | 965/965 |

## Notes

* spanner-query-benchmark runs the query planner and the query to
  display query execution stats.
* The "scanned rows" is the most useful metric, others would be
  noisy and inconsistent.
* It currently doesn't support timestamp-bound queries but it's in
  the roadmap. 

## Disclaimer

This is not an official Google product.