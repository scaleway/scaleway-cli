//go:build !wasm

package tui

import (
	"context"
	"errors"
	"os"
	"os/exec"

	applesilicon "github.com/scaleway/scaleway-cli/v2/internal/namespaces/applesilicon/v1alpha1"
	instance "github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	k8s "github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
	mongodb "github.com/scaleway/scaleway-cli/v2/internal/namespaces/mongodb/v1"
	rdb "github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/tui"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func init() {
	// SSH connectors, reusing the `scw ... server ssh` implementation.
	tui.RegisterConnector("instance server", tui.Connector{
		Hint: "ssh",
		Connect: func(ctx context.Context, client *scw.Client, loc *tui.Locality, row tui.Row) (*exec.Cmd, func(), error) {
			zone, _ := loc.Get()
			cmd, err := instance.SSHServerCmd(ctx, client, zone, row.ID, "root", 22, "")

			return cmd, nil, err
		},
	})
	tui.RegisterConnector("apple-silicon server", tui.Connector{
		Hint: "ssh",
		Connect: func(ctx context.Context, client *scw.Client, loc *tui.Locality, row tui.Row) (*exec.Cmd, func(), error) {
			zone, _ := loc.Get()
			cmd, err := applesilicon.SSHServerCmd(client, zone, row.ID, "m1", 22, "")

			return cmd, nil, err
		},
	})
	tui.RegisterConnector("rdb instance", tui.Connector{
		Hint: "connect",
		Connect: func(ctx context.Context, client *scw.Client, loc *tui.Locality, row tui.Row) (*exec.Cmd, func(), error) {
			_, region := loc.Get()
			// The CLI command requires a username; the TUI picks the
			// instance's first admin user, and the client prompts for
			// the password on the terminal.
			user, err := rdb.ConnectUser(ctx, client, region, row.ID)
			if err != nil {
				return nil, nil, err
			}
			cmd, _, err := rdb.ConnectCmd(ctx, client, region, row.ID, user, false)

			return cmd, nil, err
		},
	})
	tui.RegisterConnector("mongodb instance", tui.Connector{
		Hint: "connect",
		Connect: func(ctx context.Context, client *scw.Client, loc *tui.Locality, row tui.Row) (*exec.Cmd, func(), error) {
			_, region := loc.Get()
			cmd, err := mongodb.ConnectCmd(ctx, client, region, row.ID)

			return cmd, nil, err
		},
	})
	tui.RegisterConnector("k8s cluster", tui.Connector{
		Hint:    "k9s",
		Connect: k9sConnect,
	})
}

// k9sConnect launches k9s on the cluster: it generates the kubeconfig
// (the core of `scw k8s kubeconfig get`) in a temp file that k9s reads.
func k9sConnect(
	ctx context.Context,
	client *scw.Client,
	loc *tui.Locality,
	row tui.Row,
) (*exec.Cmd, func(), error) {
	if _, err := exec.LookPath("k9s"); err != nil {
		return nil, nil, errors.New("k9s is not installed")
	}
	_, region := loc.Get()
	raw, err := k8s.KubeconfigRaw(ctx, region, row.ID)
	if err != nil {
		return nil, nil, err
	}
	f, err := os.CreateTemp("", "scw-k9s-*.yaml")
	if err != nil {
		return nil, nil, err
	}
	if _, err := f.WriteString(raw); err != nil {
		_ = f.Close()
		os.Remove(f.Name())

		return nil, nil, err
	}
	if err := f.Close(); err != nil {
		os.Remove(f.Name())

		return nil, nil, err
	}

	return exec.Command("k9s", "--kubeconfig", f.Name()), func() { os.Remove(f.Name()) }, nil
}
