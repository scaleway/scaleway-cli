package object

import (
	"context"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// newTestClient builds a scw.Client that does not load any environment
// variables or config file, so the S3 options it exposes are fully under
// the test's control.
func newTestClient(t *testing.T, opts ...scw.ClientOption) *scw.Client {
	t.Helper()
	allOpts := append([]scw.ClientOption{scw.WithoutAuth()}, opts...)
	client, err := scw.NewClient(allOpts...)
	if err != nil {
		t.Fatalf("could not create scw client: %v", err)
	}
	return client
}

// ctxWithClient injects a scw.Client into a context the same way the CLI
// bootstrap does, so helpers relying on core.ExtractClient can use it.
func ctxWithClient(client *scw.Client) context.Context {
	return core.InjectMeta(context.Background(), &core.Meta{Client: client})
}

// clearS3Env unsets every environment variable that getS3Endpoint,
// getS3UsePathStyle and getBucketEndpoint rely on.
func clearS3Env(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		scw.ScwS3EndpointEnv,
		scw.ScwS3UsePathStyleEnv,
		scw.AwsEndpointURLS3,
		"SCW_S3_ENDPOINT",
	} {
		t.Setenv(key, "")
	}
}

func Test_getS3Endpoint(t *testing.T) {
	const region = "fr-par"

	cases := []struct {
		name           string
		customEndpoint string
		scwEnv         string
		awsEnv         string
		profileEp      string
		expected       string
	}{
		{
			name:           "CLI arg takes priority over everything",
			customEndpoint: "https://cli.example.com",
			scwEnv:         "https://scwenv.example.com",
			awsEnv:         "https://aws.example.com",
			profileEp:      "https://profile.example.com",
			expected:       "https://cli.example.com",
		},
		{
			name:      "SCW env takes priority over AWS conf and profile",
			scwEnv:    "https://scwenv.example.com",
			awsEnv:    "https://aws.example.com",
			profileEp: "https://profile.example.com",
			expected:  "https://scwenv.example.com",
		},
		{
			name:      "AWS conf takes priority over profile",
			awsEnv:    "https://aws.example.com",
			profileEp: "https://profile.example.com",
			expected:  "https://aws.example.com",
		},
		{
			name:      "profile value is used when no arg/env",
			profileEp: "https://profile.example.com",
			expected:  "https://profile.example.com",
		},
		{
			name:     "default value built from region",
			expected: "https://s3." + region + ".scw.cloud",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			clearS3Env(t)
			if c.scwEnv != "" {
				t.Setenv(scw.ScwS3EndpointEnv, c.scwEnv)
			}
			if c.awsEnv != "" {
				t.Setenv(scw.AwsEndpointURLS3, c.awsEnv)
			}

			opts := []scw.ClientOption{}
			if c.profileEp != "" {
				opts = append(opts, scw.WithS3Endpoint(c.profileEp))
			}
			client := newTestClient(t, opts...)
			ctx := ctxWithClient(client)

			got := getS3Endpoint(ctx, region, c.customEndpoint)
			if got != c.expected {
				t.Errorf("getS3Endpoint: expected %q, got %q", c.expected, got)
			}
		})
	}
}

func Test_getS3UsePathStyle(t *testing.T) {
	cases := []struct {
		name         string
		cliArg       string
		scwEnv       string
		profileValue bool
		expected     bool
	}{
		{
			name:         "CLI arg true wins over env and profile",
			cliArg:       "true",
			scwEnv:       "false",
			profileValue: true,
			expected:     true,
		},
		{
			name:         "CLI arg false wins over env and profile",
			cliArg:       "false",
			scwEnv:       "true",
			profileValue: true,
			expected:     false,
		},
		{
			name:         "SCW env true when CLI arg empty",
			cliArg:       "",
			scwEnv:       "true",
			profileValue: false,
			expected:     true,
		},
		{
			// SCW env var set to a valid boolean overrides the profile.
			name:         "SCW env false overrides profile true",
			cliArg:       "",
			scwEnv:       "false",
			profileValue: true,
			expected:     false,
		},
		{
			name:         "profile value used when CLI arg and env empty",
			cliArg:       "",
			profileValue: true,
			expected:     true,
		},
		{
			name:         "default false when nothing is set",
			cliArg:       "",
			profileValue: false,
			expected:     false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			clearS3Env(t)
			if c.scwEnv != "" {
				t.Setenv(scw.ScwS3UsePathStyleEnv, c.scwEnv)
			}

			opts := []scw.ClientOption{}
			if c.profileValue {
				opts = append(opts, scw.WithS3UsePathStyle(true))
			}
			client := newTestClient(t, opts...)
			ctx := ctxWithClient(client)

			got := getS3UsePathStyle(ctx, c.cliArg)
			if got != c.expected {
				t.Errorf("getS3UsePathStyle: expected %v, got %v", c.expected, got)
			}
		})
	}
}

func Test_getBucketEndpoint(t *testing.T) {
	const (
		region = "nl-ams"
		name   = "my-bucket"
	)

	cases := []struct {
		name         string
		custom       string
		scwEnv       string
		profileEp    string
		usePathStyle bool
		expected     string
	}{
		{
			name:         "custom endpoint, virtual-hosted style",
			custom:       "https://s3.example.com",
			usePathStyle: false,
			expected:     "https://my-bucket.s3.example.com",
		},
		{
			name:         "custom endpoint, path style",
			custom:       "https://s3.example.com",
			usePathStyle: true,
			expected:     "https://s3.example.com/my-bucket",
		},
		{
			name:         "custom endpoint with path, path style",
			custom:       "https://s3.example.com/foo",
			usePathStyle: true,
			expected:     "https://s3.example.com/foo/my-bucket",
		},
		{
			name:         "default endpoint, virtual-hosted style",
			usePathStyle: false,
			expected:     "https://my-bucket.s3.nl-ams.scw.cloud",
		},
		{
			name:         "default endpoint, path style",
			usePathStyle: true,
			expected:     "https://s3.nl-ams.scw.cloud/my-bucket",
		},
		{
			name:         "env var endpoint, virtual-hosted style",
			scwEnv:       "https://env.example.com",
			usePathStyle: false,
			expected:     "https://my-bucket.env.example.com",
		},
		{
			name:         "profile endpoint, path style",
			profileEp:    "https://profile.example.com",
			usePathStyle: true,
			expected:     "https://profile.example.com/my-bucket",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			clearS3Env(t)
			if c.scwEnv != "" {
				t.Setenv("SCW_S3_ENDPOINT", c.scwEnv)
			}

			opts := []scw.ClientOption{}
			if c.profileEp != "" {
				opts = append(opts, scw.WithS3Endpoint(c.profileEp))
			}
			client := newTestClient(t, opts...)
			ctx := ctxWithClient(client)

			got, err := getBucketEndpoint(ctx, name, region, c.custom, c.usePathStyle)
			if err != nil {
				t.Fatalf("getBucketEndpoint returned an unexpected error: %v", err)
			}
			if got != c.expected {
				t.Errorf("getBucketEndpoint: expected %q, got %q", c.expected, got)
			}
		})
	}
}

func Test_getBucketEndpoint_invalidURL(t *testing.T) {
	clearS3Env(t)
	t.Setenv("SCW_S3_ENDPOINT", "://not-a-valid-url")
	client := newTestClient(t)
	ctx := ctxWithClient(client)

	_, err := getBucketEndpoint(ctx, "my-bucket", "fr-par", "", false)
	if err == nil {
		t.Fatal("expected an error for an invalid endpoint URL, got nil")
	}
}
