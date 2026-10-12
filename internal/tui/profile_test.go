package tui

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/platform/terminal"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	profileTestOrgA  = "11111111-aaaa-bbbb-cccc-111111111111"
	profileTestOrgB  = "22222222-aaaa-bbbb-cccc-222222222222"
	profileTestProjA = "11111111-aaaa-bbbb-cccc-333333333333"
	profileTestProjB = "22222222-aaaa-bbbb-cccc-444444444444"
)

const profileTestConfig = `
access_key: "SCWAAAAAAAAAAAAAAAAA"
secret_key: "11111111-aaaa-bbbb-cccc-111111111111"
default_organization_id: "` + profileTestOrgA + `"
default_project_id: "` + profileTestProjA + `"
default_zone: "fr-par-1"
profiles:
  prod:
    access_key: "SCWBBBBBBBBBBBBBBBBB"
    secret_key: "22222222-aaaa-bbbb-cccc-222222222222"
    default_organization_id: "` + profileTestOrgB + `"
    default_project_id: "` + profileTestProjB + `"
    default_zone: "fr-par-1"
`

// TestAppProfileSwitch verifies the full ctrl+p flow: the picker lists
// every profile from the config file (default first), and selecting one
// re-creates the API client for it and shows it in the header.
func TestAppProfileSwitch(t *testing.T) {
	configPath := profileTestConfigPath(t)

	platform := terminal.NewPlatform("scw-tui-test")
	meta := &core.Meta{ConfigPathFlag: configPath, Platform: platform}
	client, err := platform.CreateClient(nil, configPath, scw.DefaultProfileName)
	if err != nil {
		t.Fatalf("CreateClient: %v", err)
	}
	res := &Resource{Name: "instance server", Title: "instance server"}
	a := NewApp(
		context.Background(), client, meta,
		map[string]*Resource{"instance server": res}, nil, nil, NewLocality("fr-par-1"),
	)
	a.width, a.height = 100, 30
	a.sizeCurrent()

	if a.profile != scw.DefaultProfileName {
		t.Fatalf("initial profile = %q, want %q", a.profile, scw.DefaultProfileName)
	}

	// ctrl+p opens the profile picker: default first, then prod.
	a.Update(tea.KeyMsg{Type: tea.KeyCtrlP})
	pm, ok := a.current().(*ProfileMenu)
	if !ok {
		t.Fatalf("ctrl+p: expected the profile picker, got %T", a.current())
	}
	if len(pm.items) != 2 || pm.items[0] != scw.DefaultProfileName || pm.items[1] != "prod" {
		t.Fatalf("picker items = %v, want [default prod]", pm.items)
	}
	if pm.cursor != 0 {
		t.Fatalf("cursor = %d, want 0 (start on the current profile)", pm.cursor)
	}

	// Move to "prod" and select it. The returned cmd carries the
	// setProfileMsg; the refresh/load-names cmds it yields are not run
	// here (they would hit the API).
	a.Update(tea.KeyMsg{Type: tea.KeyDown})
	_, cmd := a.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		a.Update(cmd())
	}

	if a.profile != "prod" {
		t.Fatalf("profile = %q, want prod", a.profile)
	}
	if p, _ := a.client.GetDefaultProjectID(); p != profileTestProjB {
		t.Fatalf("client project = %q, want %q (client not re-created)", p, profileTestProjB)
	}
	if a.project != profileTestProjB {
		t.Fatalf("app project = %q, want %q", a.project, profileTestProjB)
	}
	if a.flash != "profile set to prod" {
		t.Fatalf("flash = %q", a.flash)
	}
	if len(a.stack) != 1 {
		t.Fatalf("stack len = %d, want 1 (picker closed)", len(a.stack))
	}
	// The header leads with the active profile, top-left.
	plain := ansiRe.ReplaceAllString(a.viewHeader(), "")
	if !strings.HasPrefix(plain, "profile: prod") {
		t.Fatalf("header should lead with the profile, got:\n%s", plain)
	}
}

// profileTestConfigPath writes the two-profile test config to a temp
// file and returns its path.
func profileTestConfigPath(t *testing.T) string {
	t.Helper()
	unsetScwEnv(t)
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(configPath, []byte(profileTestConfig), 0o600); err != nil {
		t.Fatal(err)
	}

	return configPath
}

// unsetScwEnv removes the SCW_* variables for the test duration: the
// SDK reads them (t.Setenv cannot unset — an empty value would still
// be read and override the test config).
func unsetScwEnv(t *testing.T) {
	t.Helper()
	unset := func(keys ...string) {
		for _, key := range keys {
			old, exists := os.LookupEnv(key)
			if err := os.Unsetenv(key); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if exists {
					os.Setenv(key, old) //nolint:usetesting // restoring the pre-test value
				}
			})
		}
	}
	unset(
		scw.ScwAccessKeyEnv, "SCALEWAY_ACCESS_KEY",
		scw.ScwSecretKeyEnv, "SCW_TOKEN", "SCALEWAY_TOKEN",
		scw.ScwActiveProfileEnv,
		scw.ScwDefaultOrganizationIDEnv, "SCW_ORGANIZATION", "SCALEWAY_ORGANIZATION",
		scw.ScwDefaultProjectIDEnv,
		"SCW_DEFAULT_REGION", "SCW_DEFAULT_ZONE", "SCW_REGION", "SCALEWAY_REGION",
		"SCW_API_URL", "SCW_USER_AGENT",
	)
}
