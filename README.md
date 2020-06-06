# spannerbench

A Google Cloud Spanner transaction benchmarking tool.

See `examples/benchmark.yaml` for an example configuration.

## Downloads

* [Linux 64-bit](https://github.com/rakyll/spannerbench/raw/master/bin/spannerbench_linux_amd64)
* [MacOS 64-bit](https://github.com/rakyll/spannerbench/raw/master/bin/spannerbench_darwin_amd64)
* [Windows 64-bit](https://github.com/rakyll/spannerbench/raw/master/bin/spannerbench_windows_amd64)

## Usage

```
$ spannerbench -f examples/benhmark.yaml
Query stream where likes > 100
  Latency   : 7.99ms
  CPU time  : 7.14ms
  Optimizer : 2.22ms
Latency histogram:
  7.11ms    : ■■■■■■■■■■ (1)
  7.694ms   : ■■■■■■■■■■■■■■■■■■■■ (2)
  8.278ms   : ■■■■■■■■■■■■■■■■■■■■ (2)
  8.862ms   : ■■■■■■■■■■■■■■■■■■■■ (2)
  9.446ms   : ■■■■■■■■■■■■■■■■■■■■ (2)
  10.03ms   : ■■■■■■■■■■ (1)

Query stream where likes > 100 (read-only; optimizer=1)
  Latency   : 36.5ms
  CPU time  : 2.36ms
  Optimizer : 880µs
Latency histogram:
  36.37ms   : ■■■■■■■ (1)
  36.418ms  : ■■■■■■■ (1)
  36.466ms  : ■■■■■■■■■■■■■ (2)
  36.514ms  : ■■■■■■■ (1)
  36.562ms  : ■■■■■■■■■■■■■ (2)
  36.61ms   : ■■■■■■■■■■■■■■■■■■■■ (3)
```

## Notes

* The tool currently doesn't support timestamp-bound
  queries but it's in the roadmap.
* The tool currently doesn't allow to use a specified index when
  querying. It will be enabled if requested by the users.

## Disclaimer

This is not an official Google product.
