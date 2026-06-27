# Performance Testing

## Focus Areas

Performance testing at Titan goes beyond load testing. It includes profiling, benchmarking, regression detection, and capacity planning for every engine.

## Benchmarking

Go benchmarks run with `go test -bench=. -benchmem` for core library functions. Rust benchmarks use `cargo bench` with `criterion` for critical code paths. Python benchmarks use `pytest-benchmark` for ML inference and data processing. Benchmarks are run on every PR and compared against the `main` branch baseline with a regression budget of 10%.

## Profiling

### CPU Profiling
Go uses `pprof` via `net/http/pprof` with continuous profiling through Pyroscope. Rust uses `perf` and `flamegraph` for hot path identification. Python uses `cProfile` and `py-spy` for sampling profiler in production.

### Memory Profiling
Heap allocation analysis runs via `pprof` and `heaptrack`. Goroutine leak detection uses `pprof goroutine`. Object retention analysis uses `gcvis`.

### Concurrency
Race condition detection runs with `go test -race`. Deadlock detection uses `go tool trace`. Lock contention profiling uses `pprof mutex`.

## Regression Detection

A `benchstat` comparison runs in CI. If any benchmark shows more than 10% regression with statistical significance, the pipeline warns the developer. If more than 25%, it blocks the PR.

## Capacity Planning

Quarterly load tests run at 2x, 5x, and 10x projected peak traffic to identify scaling bottlenecks. Results inform resource requests, node group sizing, and provisioning budgets for the next quarter. All profiling data is collected by Pyroscope and visualized in Grafana alongside metrics and traces.