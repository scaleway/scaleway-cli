# Performance Benchmarks

This document describes the performance benchmarking system for the Scaleway CLI.

## Overview

The benchmarking system allows you to:
- **Measure performance** of CLI commands over time
- **Detect regressions** automatically before merging code
- **Track performance trends** across releases
- **Compare implementations** between branches

## Architecture

The system consists of two main components:

1. **Benchmark tests** (`*_benchmark_test.go`) - Go benchmark functions that measure command performance
2. **Benchstat tool** (`cmd/scw-benchstat`) - Wrapper that runs benchmarks and detects regressions using [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)

### Workflow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Scan project for *_benchmark_test.go files              │
└──────────────────┬──────────────────────────────────────────┘
                   ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. Run benchmarks: go test -bench=. -benchtime=1s -count=10│
└──────────────────┬──────────────────────────────────────────┘
                   ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. Save/Load baseline from testdata/benchmark.baseline     │
└──────────────────┬──────────────────────────────────────────┘
                   ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. Compare with benchstat (statistical analysis)           │
└──────────────────┬──────────────────────────────────────────┘
                   ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. Fail CI if regression > threshold (e.g., 1.5x slower)   │
└─────────────────────────────────────────────────────────────┘
```

## Quick Start

### Running Benchmarks

```bash
# Run all benchmarks for all namespaces
SCW_PROFILE=cli \
CLI_RUN_BENCHMARKS=true \
go run ./cmd/scw-benchstat

# Run benchmarks for a specific namespace
SCW_PROFILE=cli \
CLI_RUN_BENCHMARKS=true \
go run ./cmd/scw-benchstat --target-dirs=internal/namespaces/rdb/v1
```

### Creating a Baseline

The first time you run benchmarks, a baseline file is automatically created at:
```
internal/namespaces/<namespace>/<version>/testdata/benchmark.baseline
```

Example for RDB:
```
internal/namespaces/rdb/v1/testdata/benchmark.baseline
```

**Commit this baseline file** to track performance over time:

```bash
git add internal/namespaces/rdb/v1/testdata/benchmark.baseline
git commit -S -m "chore(rdb): add performance benchmark baseline"
```

### Updating the Baseline

When you intentionally make performance changes, update the baseline:

```bash
# Run benchmarks and save new baseline
CLI_RUN_BENCHMARKS=true go test \
  -bench=. \
  -benchtime=1s \
  -count=10 \
  ./internal/namespaces/rdb/v1 \
  > internal/namespaces/rdb/v1/testdata/benchmark.baseline

# Commit the updated baseline
git add internal/namespaces/rdb/v1/testdata/benchmark.baseline
git commit -S -m "chore(rdb): update benchmark baseline after optimization"
```

## Command Options

### Basic Usage

```bash
go run ./cmd/scw-benchstat [OPTIONS]
```

### Available Options

| Option | Default | Description |
|--------|---------|-------------|
| `--bench` | `.` | Benchmark pattern (Go regex) |
| `--benchtime` | `1s` | Duration per benchmark run |
| `--count` | `5` | Number of benchmark runs (for statistical accuracy) |
| `--benchmem` | `false` | Measure memory allocations |
| `--fail-metrics` | - | Comma-separated metrics to check: `time/op`, `B/op`, `allocs/op` |
| `--threshold` | `1.5` | Regression threshold (1.5 = fail if 50% slower) |
| `--install-benchstat` | `false` | Auto-install benchstat if not found |
| `--target-dirs` | (all) | Comma-separated directories to benchmark |
| `--verbose` | `false` | Enable verbose output |

### Examples

#### Run specific benchmarks only
```bash
SCW_PROFILE=cli CLI_RUN_BENCHMARKS=true \
go run ./cmd/scw-benchstat --bench="BenchmarkInstance.*"
```

#### Strict regression detection (20% threshold)
```bash
SCW_PROFILE=cli CLI_RUN_BENCHMARKS=true \
go run ./cmd/scw-benchstat \
  --threshold=1.2 \
  --fail-metrics=time/op,B/op
```

#### Quick run with fewer iterations
```bash
SCW_PROFILE=cli CLI_RUN_BENCHMARKS=true \
go run ./cmd/scw-benchstat --count=3 --benchtime=500ms
```

#### Precise measurement with more runs
```bash
SCW_PROFILE=cli CLI_RUN_BENCHMARKS=true \
go run ./cmd/scw-benchstat --count=20 --benchtime=3s
```

## Benchmark Metrics

Each benchmark reports three key metrics:

| Metric | Unit | Description |
|--------|------|-------------|
| **time/op** | ns/op | Nanoseconds per operation (execution time) |
| **B/op** | bytes | Bytes allocated per operation (memory usage) |
| **allocs/op** | count | Number of memory allocations per operation |

### Example Output

```
BenchmarkInstanceGet-11     	       5	 235361983 ns/op	  379590 B/op	    4369 allocs/op
BenchmarkBackupGet-11       	      15	  70244775 ns/op	  272052 B/op	    2845 allocs/op
BenchmarkBackupList-11      	      12	  92052913 ns/op	  284125 B/op	    2994 allocs/op
BenchmarkDatabaseList-11    	       9	 164681597 ns/op	  299008 B/op	    3152 allocs/op
```

**Reading the output:**
- **Column 1**: Benchmark name + number of parallel workers
- **Column 2**: Number of iterations executed
- **Column 3**: Average time per operation (in nanoseconds)
- **Column 4**: Average memory allocated per operation
- **Column 5**: Average number of allocations per operation

## Interpreting Results

### Comparison Output (benchstat)

When comparing with a baseline, `benchstat` shows statistical differences:

```
name               old time/op    new time/op    delta
InstanceGet-11      235ms ± 2%     220ms ± 3%   -6.38%  (p=0.000 n=10+10)
BackupGet-11       70.2ms ± 5%    72.1ms ± 4%   +2.70%  (p=0.043 n=10+10)

name               old alloc/op   new alloc/op   delta
InstanceGet-11      380kB ± 1%     383kB ± 1%   +0.79%  (p=0.000 n=10+10)
```

**Understanding the columns:**
- `±` : Standard deviation (variance in measurements)
- `delta` : Percentage change (negative = improvement, positive = regression)
- `p=0.000` : Statistical significance (p < 0.05 = significant change)
- `n=10+10` : Number of samples in each measurement

### Regression Detection

The tool fails if any metric exceeds the threshold:

```bash
❌ Performance regressions detected (threshold: 1.5x):
  - BenchmarkInstanceGet: time/op 2.1x slower (235ms → 493ms)
  - BenchmarkBackupList: B/op 1.8x more memory (284KB → 512KB)
```


## Writing Benchmarks

### Basic Structure

Create a file named `custom_benchmark_test.go` in your namespace:

```go
package mynamespace_test

import (
    "os"
    "testing"
    "time"
)

func BenchmarkMyCommand(b *testing.B) {
    // Skip unless explicitly enabled
    if os.Getenv("CLI_RUN_BENCHMARKS") != "true" {
        b.Skip("Skipping benchmark. Set CLI_RUN_BENCHMARKS=true to run.")
    }

    // Setup: create resources, clients, etc.
    client, meta, executeCmd := setupBenchmark(b)
    
    // Reset timer to exclude setup time
    b.ResetTimer()
    b.ReportAllocs() // Track memory allocations
    
    // Benchmark loop (Go adjusts b.N automatically)
    for range b.N {
        executeCmd([]string{"scw", "my-namespace", "my-command", "arg1"})
    }
    
    b.StopTimer()
}
```

### Best Practices

#### 1. **Resource Management**

Reuse expensive resources across benchmarks:

```go
var (
    sharedInstance   *MyResource
    sharedInstanceMu sync.Mutex
)

func getOrCreateSharedInstance(b *testing.B) *MyResource {
    b.Helper()
    sharedInstanceMu.Lock()
    defer sharedInstanceMu.Unlock()
    
    if sharedInstance != nil {
        b.Log("Reusing existing shared instance")
        return sharedInstance
    }
    
    b.Log("Creating shared instance...")
    sharedInstance = createExpensiveResource()
    return sharedInstance
}

func TestMain(m *testing.M) {
    code := m.Run()
    cleanupSharedInstance()
    os.Exit(code)
}
```

#### 2. **Proper Timing**

Exclude setup and cleanup from measurements:

```go
func BenchmarkMyCommand(b *testing.B) {
    // Setup (not timed)
    resource := createResource()
    b.Cleanup(func() { deleteResource(resource) })
    
    b.ResetTimer()  // ⏱️ Start timing here
    
    for range b.N {
        executeCmd(...)
    }
    
    b.StopTimer()  // ⏹️ Stop timing before cleanup
}
```

#### 3. **Avoid Interference**

Serialize operations that cannot run in parallel:

```go
var operationMu sync.Mutex

func BenchmarkOperation(b *testing.B) {
    operationMu.Lock()
    defer operationMu.Unlock()
    
    // Only one benchmark can run this at a time
    performExclusiveOperation()
}
```

## Advanced Usage

### Custom Statistics

Enable detailed statistics with `CLI_BENCH_TRACE`:

```bash
CLI_BENCH_TRACE=true CLI_RUN_BENCHMARKS=true go test -bench=. -v
```

Output:
```
Distribution (n=100): min=200ms median=235ms mean=238ms p95=280ms max=320ms
```

### CPU Profiling

Profile a slow benchmark:

```bash
CLI_RUN_BENCHMARKS=true go test -bench=BenchmarkSlowCommand \
  -cpuprofile=cpu.prof

go tool pprof -http=:8080 cpu.prof
```

### Memory Profiling

Analyze memory allocations:

```bash
CLI_RUN_BENCHMARKS=true go test -bench=BenchmarkMemoryHeavy \
  -memprofile=mem.prof

go tool pprof -http=:8080 mem.prof
```

### Comparing Branches

```bash
# Run on main branch
git checkout main
CLI_RUN_BENCHMARKS=true go test -bench=. -count=10 > /tmp/main.txt

# Run on feature branch
git checkout feature/my-optimization
CLI_RUN_BENCHMARKS=true go test -bench=. -count=10 > /tmp/feature.txt

# Compare
benchstat /tmp/main.txt /tmp/feature.txt
```

## Troubleshooting

### "benchstat not found"

Install benchstat:
```bash
go install golang.org/x/perf/cmd/benchstat@latest
```

Or use auto-install:
```bash
go run ./cmd/scw-benchstat --install-benchstat
```

### "signal: killed"

The benchmark process was killed (timeout or OOM). Try:
- Reduce `--count` or `--benchtime`
- Run specific benchmarks only with `--bench`
- Check system resources (RAM, CPU)

### "409 Conflict" errors

Resources are in a transient state. The cleanup retry mechanism should handle this automatically. If it persists:
- Increase retry attempts in `cleanupWithRetry`
- Add more wait time between operations

### Inconsistent Results

Run with more iterations for better statistical accuracy:
```bash
go run ./cmd/scw-benchstat --count=20 --benchtime=3s
```

## FAQ

**Q: How often should I update the baseline?**  
A: Update it after intentional performance changes or major refactors. Don't update for every PR unless you've specifically optimized performance.

**Q: What's a good threshold value?**  
A: Start with 1.5 (50% regression). Adjust based on your tolerance:
- 1.2 (20%) = strict
- 1.5 (50%) = balanced
- 2.0 (100%) = lenient

**Q: Should I commit baseline files?**  
A: Yes! Baselines should be tracked in git to enable comparison across branches and time.

**Q: How do I run benchmarks locally?**  
A: Use the same command as CI with your local credentials:
```bash
SCW_PROFILE=cli CLI_RUN_BENCHMARKS=true go run ./cmd/scw-benchstat
```

**Q: Can I benchmark other namespaces (instance, vpc, etc.)?**  
A: Yes! Create a `custom_benchmark_test.go` file in any namespace directory following the same pattern as RDB.

## References

- [Go Benchmarking Documentation](https://pkg.go.dev/testing#hdr-Benchmarks)
- [benchstat Tool](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)
- [Performance Profiling in Go](https://go.dev/blog/pprof)

