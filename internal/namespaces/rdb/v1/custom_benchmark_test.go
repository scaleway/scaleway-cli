package rdb_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
)

// Benchmarks for the RDB commands listed in issue #4429.
//
// They replay the cassettes recorded by the tests of this package, so they
// measure the CLI cost of a command, not the latency of the RDB API.
// See docs/developer.md.

func Benchmark_GetInstance(b *testing.B) {
	core.Benchmark(&core.BenchmarkConfig{
		Commands: rdb.GetCommands(),
		Cassette: "test-get-instance-simple",
		Args: []string{
			"scw", "rdb", "instance", "get", "ee17644d-8d1c-4fc1-86db-3ad64dd8e2dd",
		},
	})(b)
}

func Benchmark_GetBackup(b *testing.B) {
	core.Benchmark(&core.BenchmarkConfig{
		Commands: rdb.GetCommands(),
		Cassette: "test-list-backup-simple",
		Args: []string{
			"scw", "rdb", "backup", "get", "8c2e48c7-31c1-4e15-9556-08040b8701a9",
		},
	})(b)
}

func Benchmark_ListBackup(b *testing.B) {
	core.Benchmark(&core.BenchmarkConfig{
		Commands: rdb.GetCommands(),
		Cassette: "test-list-backup-simple",
		Args: []string{
			"scw", "rdb", "backup", "list",
			"instance-id=319e7311-4804-4035-9b70-c762c694a3e0",
		},
	})(b)
}

func Benchmark_ListDatabase(b *testing.B) {
	core.Benchmark(&core.BenchmarkConfig{
		Commands: rdb.GetCommands(),
		Cassette: "test-user-get-url-with-custom-database",
		Args: []string{
			"scw", "rdb", "database", "list",
			"instance-id=d3c58633-54f4-4529-a4a8-3ad2bd6b3b09",
			"name=custom-db",
			// The cassette was recorded with this argument explicitly set.
			"skip-size-retrieval=false",
		},
	})(b)
}
