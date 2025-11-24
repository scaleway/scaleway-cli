package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	benchstatVersion = "v0.0.0-20251023143056-3684bd442cc8"
)

type config struct {
	bench       string
	benchtime   string
	count       int
	benchmem    bool
	failMetrics []string
	threshold   float64
	installTool bool
	targetDirs  []string
	verbose     bool
	update      bool
}

type benchResult struct {
	name        string
	timePerOp   float64
	bytesPerOp  int64
	allocsPerOp int64
}

func main() {
	cfg := parseFlags()

	if cfg.installTool {
		if err := installBenchstat(); err != nil {
			log.Fatalf("failed to install benchstat: %v", err)
		}
	}

	if !isBenchstatAvailable() {
		log.Fatalf("benchstat not found in PATH; install golang.org/x/perf/cmd/benchstat@%s in your environment or run with --install-benchstat", benchstatVersion)
	}

	if len(cfg.targetDirs) == 0 {
		cfg.targetDirs = findBenchmarkDirs()
	}

	if len(cfg.targetDirs) == 0 {
		log.Fatal("no benchmark directories found")
	}

	var hadError bool
	for _, dir := range cfg.targetDirs {
		if err := runBenchmarksForDir(cfg, dir); err != nil {
			log.Printf("❌ failed to run benchmarks for %s: %v", dir, err)
			hadError = true
		}
	}

	if hadError {
		os.Exit(1)
	}
}

func parseFlags() config {
	cfg := config{}

	flag.StringVar(&cfg.bench, "bench", ".", "benchmark pattern to run")
	flag.StringVar(&cfg.benchtime, "benchtime", "1s", "benchmark time")
	flag.IntVar(&cfg.count, "count", 5, "number of benchmark runs")
	flag.BoolVar(&cfg.benchmem, "benchmem", false, "include memory allocation stats")
	flag.Float64Var(&cfg.threshold, "threshold", 1.5, "performance regression threshold (e.g., 1.5 = 50% slower)")
	flag.BoolVar(&cfg.installTool, "install-benchstat", false, "install benchstat tool if not found")
	flag.BoolVar(&cfg.verbose, "verbose", false, "verbose output")
	flag.BoolVar(&cfg.update, "update", false, "update baseline files instead of comparing")

	var failMetricsStr string
	flag.StringVar(&failMetricsStr, "fail-metrics", "", "comma-separated list of metrics to check for regressions (default: time/op)")

	var targetDirsStr string
	flag.StringVar(&targetDirsStr, "target-dirs", "", "comma-separated list of directories to benchmark")

	flag.Parse()

	if failMetricsStr != "" {
		cfg.failMetrics = strings.Split(failMetricsStr, ",")
	} else {
		cfg.failMetrics = []string{"time/op"}
	}

	if targetDirsStr != "" {
		cfg.targetDirs = strings.Split(targetDirsStr, ",")
	}

	return cfg
}

func installBenchstat() error {
	fmt.Printf("Installing benchstat@%s...\n", benchstatVersion)
	cmd := exec.Command("go", "install", fmt.Sprintf("golang.org/x/perf/cmd/benchstat@%s", benchstatVersion))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isBenchstatAvailable() bool {
	_, err := exec.LookPath("benchstat")
	return err == nil
}

func findBenchmarkDirs() []string {
	var dirs []string

	err := filepath.WalkDir("internal/namespaces", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), "_benchmark_test.go") {
			dir := filepath.Dir(path)
			dirs = append(dirs, dir)
		}

		return nil
	})

	if err != nil {
		log.Printf("error scanning for benchmark directories: %v", err)
	}

	return dirs
}

func runBenchmarksForDir(cfg config, dir string) error {
	fmt.Printf(">>> Running benchmarks for %s\n", dir)

	baselineFile := filepath.Join(dir, "testdata", "benchmark.baseline")

	// Run benchmarks
	newResults, err := runBenchmarks(cfg, dir)
	if err != nil {
		return fmt.Errorf("failed to run benchmarks: %w", err)
	}

	// Update mode: always overwrite baseline
	if cfg.update {
		if err := saveBaseline(baselineFile, newResults); err != nil {
			return fmt.Errorf("failed to update baseline: %w", err)
		}
		fmt.Printf("✅ Baseline updated: %s\n", baselineFile)
		return nil
	}

	// Check if baseline exists
	if _, err := os.Stat(baselineFile); os.IsNotExist(err) {
		fmt.Printf("No baseline found at %s. Creating new baseline.\n", baselineFile)
		if err := saveBaseline(baselineFile, newResults); err != nil {
			return fmt.Errorf("failed to save baseline: %w", err)
		}
		fmt.Printf("Baseline saved to %s\n", baselineFile)
		return nil
	}

	// Compare with baseline
	return compareWithBaseline(cfg, baselineFile, newResults)
}

func runBenchmarks(cfg config, dir string) (string, error) {
	args := []string{"test", "-bench=" + cfg.bench, "-benchtime=" + cfg.benchtime, "-count=" + strconv.Itoa(cfg.count)}

	if cfg.benchmem {
		args = append(args, "-benchmem")
	}

	args = append(args, "./"+dir)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Env = append(os.Environ(), "CLI_RUN_BENCHMARKS=true")

	if cfg.verbose {
		fmt.Printf("Running: go %s\n", strings.Join(args, " "))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("benchmark execution failed: %w\nOutput: %s", err, output)
	}

	return string(output), nil
}

func saveBaseline(filename, content string) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filename, []byte(content), 0644)
}

func compareWithBaseline(cfg config, baselineFile, newResults string) error {
	// Create temporary file for new results
	tmpFile, err := os.CreateTemp("", "benchmark-new-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(newResults); err != nil {
		return fmt.Errorf("failed to write new results: %w", err)
	}
	tmpFile.Close()

	// Run benchstat comparison
	cmd := exec.Command("benchstat", "-format=csv", baselineFile, tmpFile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to compare with benchstat for %s: %w\nOutput: %s", filepath.Dir(baselineFile), err, output)
	}

	// Parse CSV output and check for regressions
	return checkForRegressions(cfg, string(output))
}

func checkForRegressions(cfg config, csvOutput string) error {
	reader := csv.NewReader(strings.NewReader(csvOutput))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse benchstat CSV output: %w", err)
	}

	if len(records) < 2 {
		fmt.Println("No benchmark comparisons found")
		return nil
	}

	// Find column indices
	header := records[0]
	nameIdx := findColumnIndex(header, "name")
	oldTimeIdx := findColumnIndex(header, "old time/op")
	newTimeIdx := findColumnIndex(header, "new time/op")
	oldBytesIdx := findColumnIndex(header, "old B/op")
	newBytesIdx := findColumnIndex(header, "new B/op")
	oldAllocsIdx := findColumnIndex(header, "old allocs/op")
	newAllocsIdx := findColumnIndex(header, "new allocs/op")

	if nameIdx == -1 {
		return fmt.Errorf("could not find 'name' column in benchstat output")
	}

	var regressions []string

	for i, record := range records[1:] {
		if len(record) <= nameIdx {
			continue
		}

		benchName := record[nameIdx]

		// Check time/op regression
		if contains(cfg.failMetrics, "time/op") && oldTimeIdx != -1 && newTimeIdx != -1 {
			if regression := checkMetricRegression(record, oldTimeIdx, newTimeIdx, cfg.threshold); regression != "" {
				regressions = append(regressions, fmt.Sprintf("%s: time/op %s", benchName, regression))
			}
		}

		// Check B/op regression
		if contains(cfg.failMetrics, "B/op") && oldBytesIdx != -1 && newBytesIdx != -1 {
			if regression := checkMetricRegression(record, oldBytesIdx, newBytesIdx, cfg.threshold); regression != "" {
				regressions = append(regressions, fmt.Sprintf("%s: B/op %s", benchName, regression))
			}
		}

		// Check allocs/op regression
		if contains(cfg.failMetrics, "allocs/op") && oldAllocsIdx != -1 && newAllocsIdx != -1 {
			if regression := checkMetricRegression(record, oldAllocsIdx, newAllocsIdx, cfg.threshold); regression != "" {
				regressions = append(regressions, fmt.Sprintf("%s: allocs/op %s", benchName, regression))
			}
		}

		if cfg.verbose && i < 5 { // Show first few comparisons
			fmt.Printf("  %s: compared\n", benchName)
		}
	}

	// Print full benchstat output
	fmt.Println("Benchmark comparison results:")
	fmt.Println(csvOutput)

	if len(regressions) > 0 {
		fmt.Printf("\n❌ Performance regressions detected (threshold: %.1fx):\n", cfg.threshold)
		for _, regression := range regressions {
			fmt.Printf("  - %s\n", regression)
		}
		return fmt.Errorf("performance regressions detected")
	}

	fmt.Printf("✅ No significant performance regressions detected (threshold: %.1fx)\n", cfg.threshold)
	return nil
}

func findColumnIndex(header []string, columnName string) int {
	for i, col := range header {
		if strings.Contains(strings.ToLower(col), strings.ToLower(columnName)) {
			return i
		}
	}
	return -1
}

func checkMetricRegression(record []string, oldIdx, newIdx int, threshold float64) string {
	if oldIdx >= len(record) || newIdx >= len(record) {
		return ""
	}

	oldVal, err1 := parseMetricValue(record[oldIdx])
	newVal, err2 := parseMetricValue(record[newIdx])

	if err1 != nil || err2 != nil || oldVal == 0 {
		return ""
	}

	ratio := newVal / oldVal
	if ratio > threshold {
		return fmt.Sprintf("%.2fx slower (%.2f → %.2f)", ratio, oldVal, newVal)
	}

	return ""
}

func parseMetricValue(s string) (float64, error) {
	// Remove common suffixes and parse
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "ns", "")
	s = strings.ReplaceAll(s, "B", "")
	s = strings.TrimSpace(s)

	if s == "" || s == "-" {
		return 0, fmt.Errorf("empty value")
	}

	return strconv.ParseFloat(s, 64)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
