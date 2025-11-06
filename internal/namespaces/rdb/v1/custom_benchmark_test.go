package rdb_test

import (
	"bytes"
	"context"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
	rdbSDK "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Benchmarks for RDB commands.
//
// Baseline stored in testdata/benchmark.baseline (like golden files).
//
// To compare performance:
//
//	benchstat testdata/benchmark.baseline <(CLI_RUN_BENCHMARKS=true go test -bench=. -benchtime=100x .)
//
// To update baseline:
//
//	CLI_RUN_BENCHMARKS=true go test -bench=. -benchtime=100x . > testdata/benchmark.baseline

const (
	defaultCmdTimeout    = 30 * time.Second
	instanceReadyTimeout = 3 * time.Minute
)

func setupBenchmark(b *testing.B) (*scw.Client, core.TestMetadata, func(args []string) any) {
	b.Helper()

	clientOpts := []scw.ClientOption{
		scw.WithDefaultRegion(scw.RegionFrPar),
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithUserAgent("cli-benchmark-test"),
		scw.WithEnv(),
	}

	config, err := scw.LoadConfig()
	if err == nil {
		activeProfile, err := config.GetActiveProfile()
		if err == nil {
			envProfile := scw.LoadEnvProfile()
			profile := scw.MergeProfiles(activeProfile, envProfile)
			clientOpts = append(clientOpts, scw.WithProfile(profile))
		}
	}

	client, err := scw.NewClient(clientOpts...)
	if err != nil {
		b.Fatalf(
			"Failed to create Scaleway client: %v\nMake sure you have configured your credentials with 'scw config'",
			err,
		)
	}

	meta := core.TestMetadata{
		"t": b,
	}

	executeCmd := func(args []string) any {
		stdoutBuffer := &bytes.Buffer{}
		stderrBuffer := &bytes.Buffer{}
		_, result, err := core.Bootstrap(&core.BootstrapConfig{
			Args:             args,
			Commands:         rdb.GetCommands().Copy(),
			BuildInfo:        nil,
			Stdout:           stdoutBuffer,
			Stderr:           stderrBuffer,
			Client:           client,
			DisableTelemetry: true,
			DisableAliases:   true,
			OverrideEnv:      map[string]string{},
			Ctx:              context.Background(),
		})
		if err != nil {
			b.Errorf("error executing cmd (%s): %v\nstdout: %s\nstderr: %s",
				args, err, stdoutBuffer.String(), stderrBuffer.String())
		}

		return result
	}

	return client, meta, executeCmd
}

func cleanupWithRetry(b *testing.B, name string, resourceID string, cleanupFn func() error) {
	b.Helper()

	if err := cleanupFn(); err != nil {
		b.Logf("cleanup failed (%s=%s): %v; retrying...", name, resourceID, err)
		time.Sleep(2 * time.Second)
		if err2 := cleanupFn(); err2 != nil {
			b.Errorf("final cleanup failure (%s=%s): %v", name, resourceID, err2)
		}
	}
}

type benchmarkStats struct {
	timings []time.Duration
	enabled bool
}

func newBenchmarkStats() *benchmarkStats {
	return &benchmarkStats{
		enabled: os.Getenv("CLI_BENCH_TRACE") == "true",
		timings: make([]time.Duration, 0, 1000),
	}
}

func (s *benchmarkStats) record(d time.Duration) {
	s.timings = append(s.timings, d)
}

func (s *benchmarkStats) getMean() time.Duration {
	if len(s.timings) == 0 {
		return 0
	}

	var sum time.Duration
	for _, t := range s.timings {
		sum += t
	}

	return sum / time.Duration(len(s.timings))
}

func (s *benchmarkStats) report(b *testing.B) {
	b.Helper()

	if !s.enabled || len(s.timings) == 0 {
		return
	}

	sort.Slice(s.timings, func(i, j int) bool {
		return s.timings[i] < s.timings[j]
	})

	minVal := s.timings[0]
	maxVal := s.timings[len(s.timings)-1]
	median := s.timings[len(s.timings)/2]
	p95 := s.timings[int(float64(len(s.timings))*0.95)]
	mean := s.getMean()

	b.Logf("Distribution (n=%d): min=%v median=%v mean=%v p95=%v max=%v",
		len(s.timings), minVal, median, mean, p95, maxVal)
}

func BenchmarkInstanceGet(b *testing.B) {
	if os.Getenv("CLI_RUN_BENCHMARKS") != "true" {
		b.Skip("Skipping benchmark. Set CLI_RUN_BENCHMARKS=true to run.")
	}

	client, meta, executeCmd := setupBenchmark(b)

	ctx := &core.BeforeFuncCtx{
		Client:     client,
		ExecuteCmd: executeCmd,
		Meta:       meta,
	}
	err := createInstanceDirect(engine)(ctx)
	if err != nil {
		b.Fatalf("Failed to create instance: %v", err)
	}

	instance := meta["Instance"].(rdb.CreateInstanceResult).Instance

	b.Cleanup(func() {
		afterCtx := &core.AfterFuncCtx{
			Client:     client,
			ExecuteCmd: executeCmd,
			Meta:       meta,
		}
		cleanupWithRetry(b, "instance", instance.ID, func() error {
			return deleteInstanceDirect()(afterCtx)
		})
	})

	stats := newBenchmarkStats()
	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		start := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), defaultCmdTimeout)
		done := make(chan any, 1)

		go func() {
			done <- executeCmd([]string{"scw", "rdb", "instance", "get", instance.ID})
		}()

		select {
		case <-done:
			stats.record(time.Since(start))
		case <-ctx.Done():
			cancel()
			b.Fatalf("command timeout after %v", defaultCmdTimeout)
		}
		cancel()
	}

	b.StopTimer()
	stats.report(b)
}

func BenchmarkBackupGet(b *testing.B) {
	if os.Getenv("CLI_RUN_BENCHMARKS") != "true" {
		b.Skip("Skipping benchmark. Set CLI_RUN_BENCHMARKS=true to run.")
	}

	client, meta, executeCmd := setupBenchmark(b)

	ctx := &core.BeforeFuncCtx{
		Client:     client,
		ExecuteCmd: executeCmd,
		Meta:       meta,
	}
	err := createInstanceDirect(engine)(ctx)
	if err != nil {
		b.Fatalf("Failed to create instance: %v", err)
	}

	instance := meta["Instance"].(rdb.CreateInstanceResult).Instance

	if err := waitForInstanceReady(executeCmd, instance.ID, instanceReadyTimeout); err != nil {
		b.Fatalf("Instance not ready: %v", err)
	}

	err = createBackupDirect("Backup")(ctx)
	if err != nil {
		b.Fatalf("Failed to create backup: %v", err)
	}

	backup := meta["Backup"].(*rdbSDK.DatabaseBackup)

	b.Cleanup(func() {
		afterCtx := &core.AfterFuncCtx{
			Client:     client,
			ExecuteCmd: executeCmd,
			Meta:       meta,
		}
		cleanupWithRetry(b, "backup", backup.ID, func() error {
			return deleteBackupDirect("Backup")(afterCtx)
		})
		cleanupWithRetry(b, "instance", instance.ID, func() error {
			return deleteInstanceDirect()(afterCtx)
		})
	})

	stats := newBenchmarkStats()
	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		start := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), defaultCmdTimeout)
		done := make(chan any, 1)

		go func() {
			done <- executeCmd([]string{"scw", "rdb", "backup", "get", backup.ID})
		}()

		select {
		case <-done:
			stats.record(time.Since(start))
		case <-ctx.Done():
			cancel()
			b.Fatalf("command timeout after %v", defaultCmdTimeout)
		}
		cancel()
	}

	b.StopTimer()
	stats.report(b)
}

func BenchmarkBackupList(b *testing.B) {
	if os.Getenv("CLI_RUN_BENCHMARKS") != "true" {
		b.Skip("Skipping benchmark. Set CLI_RUN_BENCHMARKS=true to run.")
	}

	client, meta, executeCmd := setupBenchmark(b)

	ctx := &core.BeforeFuncCtx{
		Client:     client,
		ExecuteCmd: executeCmd,
		Meta:       meta,
	}
	err := createInstanceDirect(engine)(ctx)
	if err != nil {
		b.Fatalf("Failed to create instance: %v", err)
	}

	instance := meta["Instance"].(rdb.CreateInstanceResult).Instance

	if err := waitForInstanceReady(executeCmd, instance.ID, instanceReadyTimeout); err != nil {
		b.Fatalf("Instance not ready: %v", err)
	}

	err = createBackupDirect("Backup1")(ctx)
	if err != nil {
		b.Fatalf("Failed to create backup 1: %v", err)
	}
	err = createBackupDirect("Backup2")(ctx)
	if err != nil {
		b.Fatalf("Failed to create backup 2: %v", err)
	}

	backup1 := meta["Backup1"].(*rdbSDK.DatabaseBackup)
	backup2 := meta["Backup2"].(*rdbSDK.DatabaseBackup)

	b.Cleanup(func() {
		afterCtx := &core.AfterFuncCtx{
			Client:     client,
			ExecuteCmd: executeCmd,
			Meta:       meta,
		}
		cleanupWithRetry(b, "backup1", backup1.ID, func() error {
			return deleteBackupDirect("Backup1")(afterCtx)
		})
		cleanupWithRetry(b, "backup2", backup2.ID, func() error {
			return deleteBackupDirect("Backup2")(afterCtx)
		})
		cleanupWithRetry(b, "instance", instance.ID, func() error {
			return deleteInstanceDirect()(afterCtx)
		})
	})

	stats := newBenchmarkStats()
	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		start := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), defaultCmdTimeout)
		done := make(chan any, 1)

		go func() {
			done <- executeCmd([]string{"scw", "rdb", "backup", "list", "instance-id=" + instance.ID})
		}()

		select {
		case <-done:
			stats.record(time.Since(start))
		case <-ctx.Done():
			cancel()
			b.Fatalf("command timeout after %v", defaultCmdTimeout)
		}
		cancel()
	}

	b.StopTimer()
	stats.report(b)
}

func BenchmarkDatabaseList(b *testing.B) {
	if os.Getenv("CLI_RUN_BENCHMARKS") != "true" {
		b.Skip("Skipping benchmark. Set CLI_RUN_BENCHMARKS=true to run.")
	}

	client, meta, executeCmd := setupBenchmark(b)

	ctx := &core.BeforeFuncCtx{
		Client:     client,
		ExecuteCmd: executeCmd,
		Meta:       meta,
	}
	err := createInstanceDirect(engine)(ctx)
	if err != nil {
		b.Fatalf("Failed to create instance: %v", err)
	}

	instance := meta["Instance"].(rdb.CreateInstanceResult).Instance

	b.Cleanup(func() {
		afterCtx := &core.AfterFuncCtx{
			Client:     client,
			ExecuteCmd: executeCmd,
			Meta:       meta,
		}
		cleanupWithRetry(b, "instance", instance.ID, func() error {
			return deleteInstanceDirect()(afterCtx)
		})
	})

	stats := newBenchmarkStats()
	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		start := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), defaultCmdTimeout)
		done := make(chan any, 1)

		go func() {
			done <- executeCmd([]string{"scw", "rdb", "database", "list", "instance-id=" + instance.ID})
		}()

		select {
		case <-done:
			stats.record(time.Since(start))
		case <-ctx.Done():
			cancel()
			b.Fatalf("command timeout after %v", defaultCmdTimeout)
		}
		cancel()
	}

	b.StopTimer()
	stats.report(b)
}
