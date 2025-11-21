// Package env contains a list of environment variables used to modify the behaviour of the CLI.
package env

const (
	// RunBenchmarks if set to "true" will enable benchmark execution
	RunBenchmarks = "CLI_RUN_BENCHMARKS"

	// BenchTrace if set to "true" will enable detailed benchmark statistics (min, median, mean, p95, max)
	BenchTrace = "CLI_BENCH_TRACE"

	// UpdateCassettes if set to "true" will trigger the cassettes to be recorded
	UpdateCassettes = "CLI_UPDATE_CASSETTES"

	// UpdateGoldens if set to "true" will update golden files during tests
	UpdateGoldens = "CLI_UPDATE_GOLDENS"
)
