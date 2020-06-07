# spannerbench

A Google Cloud Spanner transaction benchmarking framework.

See `examples/helloworld` for an example.

## Usage

```
$ go run examples/bench/main.go
BenchmarkReadOnly
Latency histogram:
  216ns     : ■ (1)
  2.963µs   : ■■■■■■■■■■■■■■■■■■■■ (44)
  5.71µs    : ■ (3)
  8.457µs   : ■ (1)
  11.204µs  :
  13.951µs  : ■ (1)

Benchmark
Latency histogram:
  82.870914ms: ■ (1)
  290.643631ms: ■■■■■■■■■■■■■■■■■■■■ (46)
  498.416348ms: ■ (1)
  706.189065ms: ■ (1)
  913.961782ms:
  1.121734499s: ■ (1)
```

## Notes

* The framework only reports the client-perceived latency at the moment.
* The benchmarks are run sequentially, concurrency support is in the
  roadmap but is not implemented yet.

## Disclaimer

This is not an official Google product.
