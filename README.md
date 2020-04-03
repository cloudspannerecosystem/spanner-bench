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
1:        965  37.7455ms   2.7275ms    812.5µs    965/965
2:        965   37.565ms    2.982ms    677.5µs    965/965
Benchmark2
    (scanned)    (total)      (cpu)     (plan)
1:        965   35.991ms   1.5865ms    1.063ms    965/3
2:        965   36.391ms     1.58ms    1.225ms    965/3
Benchmark3
    (scanned)    (total)      (cpu)     (plan)
2:        100  36.3065ms   1.7225ms    1.221ms    100/100
```

### Output explained...

| planner version | total rows scanned | total execution time | total CPU time | planning time | rows scanned/rows returned |
|-|-|-|-|-|-|
| 1: | 965  | 38.187ms | 3.015ms | 798.5µs | 965/965 |

## Notes

* The tool runs the query planner and the query to
  display query execution stats.
* The "scanned rows" is the most useful metric, others might be
  noisy and inconsistent.
* The tool currently doesn't support timestamp-bound queries but it's in
  the roadmap.
* The tool currently doesn't allow to use a specified index when
  querying. It will be enabled if requested by the users.

## Disclaimer

This is not an official Google product.