package testhelpers

import (
	"bytes"
	"context"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// SetupBenchmark initializes a Scaleway client and test metadata for benchmarks.
// It loads credentials from the active profile and environment variables.
func SetupBenchmark(
	b *testing.B,
	commands *core.Commands,
) (*scw.Client, core.TestMetadata, func(args []string) any) {
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
			Commands:         commands.Copy(),
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
