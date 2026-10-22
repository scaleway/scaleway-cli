package semver_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/semver"
)

func TestLatestVersion(t *testing.T) {
	t.Run("Picks the highest version regardless of order", func(t *testing.T) {
		got, err := semver.LatestVersion([]string{"7.2.7", "8.0.2", "7.0.15", "8.6.6"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "8.6.6" {
			t.Fatalf("expected 8.6.6, got %s", got)
		}
	})

	t.Run("Returns error when no version is available", func(t *testing.T) {
		if _, err := semver.LatestVersion(nil); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
