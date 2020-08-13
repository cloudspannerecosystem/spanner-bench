![spannerbench logo](https://i.imgur.com/x6Z8yEd.png)

A Google Cloud Spanner transaction benchmarking framework.

[![CircleCI](https://circleci.com/gh/cloudspannerecosystem/spanner-bench.svg?style=svg)](https://circleci.com/gh/cloudspannerecosystem/spanner-bench)


See `examples/helloworld` for an example.

## Usage

```
$ go run examples/helloworld/main.go
BenchmarkReadOnly
Latency histogram:
  65.772419ms : ■ (1)
  310.325496ms: ■■■■■■■■■■■■■■■■■■■■ (37)
  554.878574ms: ■■■ (6)
  799.431651ms: ■■ (3)
  1.043984729s: ■ (1)
  1.288537807s: ■ (2)

Benchmark
Latency histogram:
  101.510159ms: ■ (1)
  355.964311ms: ■■■■■■■■■■■■■■■■■■■■ (46)
  610.418464ms: ■ (2)
  864.872616ms:
  1.119326769s:
  1.373780922s: ■ (1)
```

## Notes

* The framework only reports the client-perceived latency at the moment.
* The benchmarks are run sequentially, concurrency support is in the
  roadmap but is not implemented yet.
* Note timestamp bound support is work in progress.
* Framework can warm-up the sessions before starting to benchmark.
  This improvement is in the roadmap.

## Disclaimer

This is not an official Google product.
