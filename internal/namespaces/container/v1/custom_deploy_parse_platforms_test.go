package container_test

import (
	"testing"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	container "github.com/scaleway/scaleway-cli/v2/internal/namespaces/container/v1"
	"github.com/stretchr/testify/assert"
)

func Test_parsePlatforms(t *testing.T) {
	t.Run("empty string returns nil (no platform constraint)", func(t *testing.T) {
		assert.Nil(t, container.ParsePlatforms(""))
	})

	t.Run("linux/amd64", func(t *testing.T) {
		got := container.ParsePlatforms("linux/amd64")
		assert.Equal(t, []ocispec.Platform{
			{OS: "linux", Architecture: "amd64"},
		}, got)
	})

	t.Run("linux/arm64", func(t *testing.T) {
		got := container.ParsePlatforms("linux/arm64")
		assert.Equal(t, []ocispec.Platform{
			{OS: "linux", Architecture: "arm64"},
		}, got)
	})

	t.Run("linux/arm/v7 keeps the variant", func(t *testing.T) {
		got := container.ParsePlatforms("linux/arm/v7")
		assert.Equal(t, []ocispec.Platform{
			{OS: "linux", Architecture: "arm", Variant: "v7"},
		}, got)
	})

	t.Run("invalid platform returns nil instead of erroring", func(t *testing.T) {
		// parsePlatforms is best-effort: a malformed value must not break the
		// build, it should fall back to the default platform.
		assert.Nil(t, container.ParsePlatforms("not-a-valid-platform"))
	})
}
