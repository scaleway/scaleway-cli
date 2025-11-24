package rdb_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/env"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
	rdbSDK "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Benchmarks for RDB commands.
//
// Baseline stored in testdata/benchmark.baseline (like golden files).
//
// Install benchstat (required for comparison):
//
//	go install golang.org/x/perf/cmd/benchstat@latest
//
// To compare performance:
//
//	benchstat testdata/benchmark.baseline <(CLI_RUN_BENCHMARKS=true go test -bench=. -benchtime=100x .)
//
// To update baseline:
//
//	CLI_RUN_BENCHMARKS=true go test -bench=. -benchtime=100x . > testdata/benchmark.baseline
//
// Or use the automated tool (installs benchstat automatically):
//
//	go run ./cmd/scw-benchstat --install-benchstat --bench=. --count=10

const (
	defaultCmdTimeout    = 30 * time.Second
	instanceReadyTimeout = 3 * time.Minute
)

var (
	sharedInstance *rdbSDK.Instance
)

// TestMain ensures shared instance cleanup
func TestMain(m *testing.M) {
	code := m.Run()
	cleanupSharedInstance()
	os.Exit(code)
}

func setupBenchmark(b *testing.B) (*scw.Client, core.TestMetadata, func(args []string) any) {
	b.Helper()
	return testhelpers.SetupBenchmark(b, rdb.GetCommands())
}

func cleanupWithRetry(b *testing.B, name string, resourceID string, cleanupFn func() error) {
	b.Helper()

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err := cleanupFn()
		if err == nil {
			return
		}

		// Check if it's a 409 Conflict using typed error
		var respErr *scw.ResponseError
		isConflict := errors.As(err, &respErr) && respErr.StatusCode == http.StatusConflict

		// Fallback: check error message for transient state keywords
		if !isConflict {
			errMsg := err.Error()
			isConflict = strings.Contains(errMsg, "transient state") || strings.Contains(errMsg, "backuping")
		}

		if isConflict && i < maxRetries-1 {
			waitTime := time.Duration(2*(i+1)) * time.Second
			b.Logf("cleanup conflict for %s=%s (attempt %d/%d), waiting %v: %v", name, resourceID, i+1, maxRetries, waitTime, err)
			time.Sleep(waitTime)
			continue
		}

		b.Errorf("cleanup failure (%s=%s) after %d attempts: %v", name, resourceID, i+1, err)
		return
	}
}

type benchmarkStats struct {
	timings []time.Duration
	enabled bool
}

func newBenchmarkStats() *benchmarkStats {
	return &benchmarkStats{
		enabled: os.Getenv(env.BenchTrace) == "true",
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

func getOrCreateSharedInstance(b *testing.B, client *scw.Client, executeCmd func([]string) any, meta core.TestMetadata) *rdbSDK.Instance {
	b.Helper()

	if sharedInstance != nil {
		b.Log("Reusing existing shared RDB instance")
		return sharedInstance
	}

	b.Log("Creating shared RDB instance for all benchmarks...")
	ctx := &core.BeforeFuncCtx{
		Client:     client,
		ExecuteCmd: executeCmd,
		Meta:       meta,
	}

	if err := createInstanceDirect(engine)(ctx); err != nil {
		b.Fatalf("Failed to create shared instance: %v", err)
	}

	instance := meta["Instance"].(rdb.CreateInstanceResult).Instance
	sharedInstance = instance

	b.Logf("Shared RDB instance created: %s", instance.ID)

	if err := waitForInstanceReady(executeCmd, instance.ID, instanceReadyTimeout); err != nil {
		b.Fatalf("Shared instance not ready: %v", err)
	}

	b.Log("Shared instance is ready")
	return sharedInstance
}

func cleanupSharedInstance() {
	if sharedInstance == nil {
		return
	}

	fmt.Printf("Cleaning up shared RDB instance: %s\n", sharedInstance.ID)

	client, err := scw.NewClient(
		scw.WithDefaultRegion(scw.RegionFrPar),
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithEnv(),
	)
	if err != nil {
		fmt.Printf("Error creating client for cleanup: %v\n", err)
		return
	}

	config, err := scw.LoadConfig()
	if err == nil {
		activeProfile, err := config.GetActiveProfile()
		if err == nil {
			envProfile := scw.LoadEnvProfile()
			profile := scw.MergeProfiles(activeProfile, envProfile)
			client, _ = scw.NewClient(
				scw.WithDefaultRegion(scw.RegionFrPar),
				scw.WithDefaultZone(scw.ZoneFrPar1),
				scw.WithProfile(profile),
				scw.WithEnv(),
			)
		}
	}

	executeCmd := func(args []string) any {
		stdoutBuffer := &bytes.Buffer{}
		stderrBuffer := &bytes.Buffer{}
		_, result, _ := core.Bootstrap(&core.BootstrapConfig{
			Args:             args,
			Commands:         rdb.GetCommands().Copy(),
			BuildInfo:        &core.BuildInfo{},
			Stdout:           stdoutBuffer,
			Stderr:           stderrBuffer,
			Client:           client,
			DisableTelemetry: true,
			DisableAliases:   true,
			OverrideEnv:      map[string]string{},
			Ctx:              context.Background(),
		})
		return result
	}

	meta := core.TestMetadata{
		"Instance": rdb.CreateInstanceResult{Instance: sharedInstance},
	}

	afterCtx := &core.AfterFuncCtx{
		Client:     client,
		ExecuteCmd: executeCmd,
		Meta:       meta,
	}

	if err := deleteInstanceDirect()(afterCtx); err != nil {
		fmt.Printf("Error deleting shared instance: %v\n", err)
		time.Sleep(2 * time.Second)
		if err2 := deleteInstanceDirect()(afterCtx); err2 != nil {
			fmt.Printf("Final cleanup failure: %v\n", err2)
		}
	}

	sharedInstance = nil
}

func BenchmarkInstanceGet(b *testing.B) {
	if os.Getenv(env.RunBenchmarks) != "true" {
		b.Skip("Skipping benchmark. Set CLI_RUN_BENCHMARKS=true to run.")
	}

	client, meta, executeCmd := setupBenchmark(b)
	instance := getOrCreateSharedInstance(b, client, executeCmd, meta)

	stats := newBenchmarkStats()
	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		start := time.Now()
		executeCmd([]string{"scw", "rdb", "instance", "get", instance.ID})
		stats.record(time.Since(start))
	}

	b.StopTimer()
	stats.report(b)
}

func BenchmarkBackupGet(b *testing.B) {
	if os.Getenv(env.RunBenchmarks) != "true" {
		b.Skip("Skipping benchmark. Set CLI_RUN_BENCHMARKS=true to run.")
	}

	client, meta, executeCmd := setupBenchmark(b)
	instance := getOrCreateSharedInstance(b, client, executeCmd, meta)

	ctx := &core.BeforeFuncCtx{
		Client:     client,
		ExecuteCmd: executeCmd,
		Meta:       meta,
	}

	meta["Instance"] = rdb.CreateInstanceResult{Instance: instance}

	if err := waitForInstanceReady(executeCmd, instance.ID, instanceReadyTimeout); err != nil {
		b.Fatalf("Instance not ready before backup: %v", err)
	}

	if err := createBackupDirect("Backup")(ctx); err != nil {
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
	})

	stats := newBenchmarkStats()
	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		start := time.Now()
		executeCmd([]string{"scw", "rdb", "backup", "get", backup.ID})
		stats.record(time.Since(start))
	}

	b.StopTimer()
	stats.report(b)
}

func BenchmarkBackupList(b *testing.B) {
	if os.Getenv(env.RunBenchmarks) != "true" {
		b.Skip("Skipping benchmark. Set CLI_RUN_BENCHMARKS=true to run.")
	}

	client, meta, executeCmd := setupBenchmark(b)
	instance := getOrCreateSharedInstance(b, client, executeCmd, meta)

	ctx := &core.BeforeFuncCtx{
		Client:     client,
		ExecuteCmd: executeCmd,
		Meta:       meta,
	}

	meta["Instance"] = rdb.CreateInstanceResult{Instance: instance}

	if err := waitForInstanceReady(executeCmd, instance.ID, instanceReadyTimeout); err != nil {
		b.Fatalf("Instance not ready before backup 1: %v", err)
	}

	if err := createBackupDirect("Backup1")(ctx); err != nil {
		b.Fatalf("Failed to create backup 1: %v", err)
	}

	if err := waitForInstanceReady(executeCmd, instance.ID, instanceReadyTimeout); err != nil {
		b.Fatalf("Instance not ready before backup 2: %v", err)
	}

	if err := createBackupDirect("Backup2")(ctx); err != nil {
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
	})

	stats := newBenchmarkStats()
	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		start := time.Now()
		executeCmd([]string{"scw", "rdb", "backup", "list", "instance-id=" + instance.ID})
		stats.record(time.Since(start))
	}

	b.StopTimer()
	stats.report(b)
}

func BenchmarkDatabaseList(b *testing.B) {
	if os.Getenv(env.RunBenchmarks) != "true" {
		b.Skip("Skipping benchmark. Set CLI_RUN_BENCHMARKS=true to run.")
	}

	client, meta, executeCmd := setupBenchmark(b)
	instance := getOrCreateSharedInstance(b, client, executeCmd, meta)

	stats := newBenchmarkStats()
	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		start := time.Now()
		executeCmd([]string{"scw", "rdb", "database", "list", "instance-id=" + instance.ID})
		stats.record(time.Since(start))
	}

	b.StopTimer()
	stats.report(b)
}
