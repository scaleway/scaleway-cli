package rdb_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/connect"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
	secretSDK "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/require"
)

func Test_ListUser(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd:       "scw rdb user list instance-id={{ .Instance.ID }}",
		Check:     core.TestCheckGolden(),
		AfterFunc: deleteInstance(),
	}))
}

func Test_CreateUser(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: fmt.Sprintf(
			"scw rdb user create instance-id={{ $.Instance.Instance.ID }} name=%s password=%s",
			name,
			password,
		),
		Check:     core.TestCheckGolden(),
		AfterFunc: deleteInstance(),
	}))

	t.Run("With password generator", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb user create instance-id={{ $.Instance.Instance.ID }} name=" + name,
		// do not check the golden as the password generated locally and on CI will necessarily be different
		Check:     core.TestCheckExitCode(0),
		AfterFunc: deleteInstance(),
	}))
}

func Test_UpdateUser(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
			core.ExecBeforeCmd(
				fmt.Sprintf(
					"scw rdb user create instance-id={{ $.Instance.Instance.ID }} name=%s password=%s",
					name,
					password,
				),
			),
		),
		Cmd: fmt.Sprintf(
			"scw rdb user update instance-id={{ $.Instance.Instance.ID }} name=%s password=Newp1ssw0rd! is-admin=true",
			name,
		),
		Check:     core.TestCheckGolden(),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Without password", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
			core.ExecBeforeCmd(
				fmt.Sprintf(
					"scw rdb user create instance-id={{ $.Instance.Instance.ID }} name=%s password=%s",
					name,
					password,
				),
			),
		),
		Cmd: fmt.Sprintf(
			"scw rdb user update instance-id={{ $.Instance.Instance.ID }} name=%s is-admin=true",
			name,
		),
		// an update without --password must not rotate the password:
		// the response carries an empty Password
		Check:     core.TestCheckGolden(),
		AfterFunc: deleteInstance(),
	}))
}

func Test_UpdateUserCliConnect(t *testing.T) {
	// The instance is created with cli-connect so the credentials secret already
	// exists: the update must add a new version on the existing secret.
	createCommand := strings.Replace(baseCommand, "--wait", "cli-connect=true --wait", 1)

	t.Run("With existing secret", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			core.ExecStoreBeforeCmd(
				"Instance",
				fmt.Sprintf(createCommand, name, engine, user, password),
			),
		),
		Cmd: fmt.Sprintf(
			"scw rdb user update instance-id={{ .Instance.ID }} name=%s cli-connect=true password=Newp1ssw0rd!",
			user,
		),
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				region, _ := ctx.Client.GetDefaultRegion()
				secretAPI := secretSDK.NewAPI(ctx.Client)
				secretName := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.ID
				path := connect.SecretPath
				list, err := secretAPI.ListSecrets(&secretSDK.ListSecretsRequest{
					Region: region,
					Name:   &secretName,
					Path:   &path,
				}, scw.WithAllPages())
				require.NoError(t, err)
				require.Len(t, list.Secrets, 1)

				secret := list.Secrets[0]
				require.Equal(t, []string{string(connect.RDB), string(region)}, secret.Tags)

				version, err := secretAPI.AccessSecretVersion(
					&secretSDK.AccessSecretVersionRequest{
						Region:   region,
						SecretID: secret.ID,
						Revision: "latest",
					},
				)
				require.NoError(t, err)

				var creds connect.Credentials
				require.NoError(t, json.Unmarshal(version.Data, &creds))
				require.Equal(t, "Newp1ssw0rd!", creds[user])
			},
		),
		AfterFunc: deleteInstance(),
	}))
}
