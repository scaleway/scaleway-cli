package core

import (
	"bytes"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// BenchmarkConfig configures a command benchmark. See Benchmark.
type BenchmarkConfig struct {
	// Array of command to load (see main.go)
	Commands *Commands

	// Name of the cassette to replay, without its extension, as stored in the
	// testdata folder of the benchmarked package. The cassette is usually one
	// already recorded by the test suite.
	Cassette string

	// The command line to benchmark, as a program would receive it.
	Args []string

	// DefaultRegion to use with scw client (default: scw.RegionFrPar)
	DefaultRegion scw.Region

	// DefaultZone to use with scw client (default: scw.ZoneFrPar1)
	DefaultZone scw.Zone
}

// Benchmark measures how long a CLI command takes, replaying a recorded
// cassette instead of calling the API. It therefore measures what the CLI
// itself costs, argument parsing, SDK marshaling and output rendering, and
// leaves out the API latency, which the CLI has no control over.
//
// Benchmarks run offline: they need no credentials and create no resource.
func Benchmark(config *BenchmarkConfig) func(b *testing.B) {
	return func(b *testing.B) {
		b.Helper()

		// Waiters must not sleep while replaying a cassette.
		DefaultRetryInterval = new(0 * time.Second)

		cassettePath := filepath.Join(".", "testdata", config.Cassette+".cassette")

		httpClient, stop, err := NewReplayHTTPClient(cassettePath)
		if err != nil {
			b.Fatalf("cannot replay cassette %s: %v", cassettePath, err)
		}

		b.Cleanup(func() {
			if err := stop(); err != nil {
				b.Errorf("cannot stop the cassette recorder: %v", err)
			}
		})

		client, err := newBenchmarkClient(config, httpClient)
		if err != nil {
			b.Fatalf("cannot create the benchmark client: %v", err)
		}

		homeDir := b.TempDir()
		overrideEnv := map[string]string{
			"HOME":             homeDir,
			scw.ScwCacheDirEnv: homeDir,
		}

		buildInfo := &BuildInfo{
			Version:   version.Must(version.NewSemver("v0.0.0+bench")),
			BuildDate: "unknown",
			GoVersion: "unknown",
			GitBranch: "unknown",
			GitCommit: "unknown",
			GoArch:    "unknown",
			GoOS:      "unknown",
		}

		ctx := b.Context()

		b.ReportAllocs()

		for b.Loop() {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			// Commands are copied because Bootstrap mutates them.
			_, _, err := Bootstrap(ctx, &BootstrapConfig{
				Args:             config.Args,
				Commands:         config.Commands.Copy(),
				BuildInfo:        buildInfo,
				Stdout:           stdout,
				Stderr:           stderr,
				Client:           client,
				DisableTelemetry: true,
				DisableAliases:   true,
				OverrideEnv:      overrideEnv,
			})
			if err != nil {
				b.Fatalf("command %s failed: %v\nstderr: %s", config.Args, err, stderr)
			}
		}
	}
}

// newBenchmarkClient creates a client with the same dummy credentials as the
// test client, so that it matches the requests stored in the cassettes.
func newBenchmarkClient(config *BenchmarkConfig, httpClient *http.Client) (*scw.Client, error) {
	region := config.DefaultRegion
	if region == "" {
		region = scw.RegionFrPar
	}

	zone := config.DefaultZone
	if zone == "" {
		zone = scw.ZoneFrPar1
	}

	return scw.NewClient(
		scw.WithDefaultRegion(region),
		scw.WithDefaultZone(zone),
		scw.WithAuth("SCWXXXXXXXXXXXXXXXXX", "11111111-1111-1111-1111-111111111111"),
		scw.WithDefaultOrganizationID("11111111-1111-1111-1111-111111111111"),
		scw.WithDefaultProjectID("11111111-1111-1111-1111-111111111111"),
		scw.WithUserAgent("cli-benchmark"),
		scw.WithHTTPClient(httpClient),
	)
}
