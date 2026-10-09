package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/config"
)

// TestTuiConfigRoundTrip verifies the tui section of cli.yaml is loaded
// and survives a Save (scw alias / scw init rewrite the file from the
// template, which must not drop user settings).
func TestTuiConfigRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cli.yaml")
	content := `output: human
tui:
    keybindings:
        refresh: R
    commands:
        mydb: rdb instance
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Tui == nil {
		t.Fatal("tui section not loaded")
	}
	if cfg.Tui.KeyBindings["refresh"] != "R" {
		t.Fatalf("keybindings = %v", cfg.Tui.KeyBindings)
	}
	if cfg.Tui.Commands["mydb"] != "rdb instance" {
		t.Fatalf("commands = %v", cfg.Tui.Commands)
	}

	if err := cfg.Save(); err != nil {
		t.Fatal(err)
	}
	saved, err := config.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if saved.Tui == nil ||
		saved.Tui.KeyBindings["refresh"] != "R" ||
		saved.Tui.Commands["mydb"] != "rdb instance" {
		t.Fatalf("tui section lost on Save")
	}
}
