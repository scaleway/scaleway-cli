package instance_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	instanceCLIV1 "github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	instanceCLIV2 "github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v2alpha1"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
)

const cloudInitScript = `#cloud-config
package_update: true
package_upgrade: true
packages:
  - sshfs
  - shellcheck`

func createTemplate(args string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"Template",
		fmt.Sprintf("scw instance template create name=%s %s", core.GetRandomName("tpl"), args),
	)
}

func deleteTemplate() core.AfterFunc {
	return func(ctx *core.AfterFuncCtx) error {
		return core.ExecAfterCmd("scw instance template delete {{ .Template.ID }}")(ctx)
	}
}

func Test_Template(t *testing.T) {
	cmds := instanceCLIV2.GetCommands()
	serverType1 := "BASIC2-A2C-4G"

	t.Run("Create", core.Test(&core.TestConfig{
		Commands: cmds,
		Cmd:      "scw instance template create tags.0=tmpl-tag server-tags.0=server-tag public-ipv4-count=2 public-ipv6-count=4 server-type=" + serverType1,
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				if tpl, ok := ctx.Result.(*instanceSDK.Template); !ok {
					t.Errorf("result was not *instance.Template but %T", ctx.Result)
				} else {
					assert.Len(t, tpl.Tags, 1)
					assert.Equal(t, "tmpl-tag", tpl.Tags[0])
					assert.Len(t, tpl.ServerTags, 1)
					assert.Equal(t, "server-tag", tpl.ServerTags[0])
					assert.Equal(t, uint32(2), tpl.PublicIPV4Count)
					assert.Equal(t, uint32(4), tpl.PublicIPV6Count)
					assert.Equal(t, serverType1, tpl.ServerType)
				}
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			return core.ExecAfterCmd("scw instance template delete {{ .CmdResult.ID }}")(ctx)
		},
	}))

	serverType2 := "POP2-HM-8C-64G"

	t.Run("Update", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: createTemplate(
			"tags.0=tmpl-tag public-ipv4-count=2 server-type=" + serverType1,
		),
		Cmd: "scw instance template update {{ .Template.ID }} tags.0=new tags.1=tags server-tags.0=server-tag" +
			" public-ipv4-count=0 public-ipv6-count=1 server-type=" + serverType2 +
			" volumes.0.name=cli-test volumes.0.size=20GB volumes.0.volume-type=sbs volumes.0.image-label=ubuntu_resolute",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				if tpl, ok := ctx.Result.(*instanceSDK.Template); !ok {
					t.Errorf("result was not *instance.Template but %T", ctx.Result)
				} else {
					assert.Len(t, tpl.Tags, 2)
					assert.Equal(t, "new", tpl.Tags[0])
					assert.Equal(t, "tags", tpl.Tags[1])
					assert.Len(t, tpl.ServerTags, 1)
					assert.Equal(t, "server-tag", tpl.ServerTags[0])
					assert.Equal(t, uint32(0), tpl.PublicIPV4Count)
					assert.Equal(t, uint32(1), tpl.PublicIPV6Count)
					assert.Equal(t, serverType2, tpl.ServerType)
					assert.Len(t, tpl.Volumes, 1)
					assert.Equal(t, "cli-test", tpl.Volumes[0].Name)
					assert.Equal(t, 20*scw.GB, *tpl.Volumes[0].Size)
					assert.Equal(
						t,
						instanceSDK.CreateServerRequestServerVolumeVolumeTypeSbs,
						tpl.Volumes[0].VolumeType,
					)
					assert.Equal(t, "ubuntu_resolute", *tpl.Volumes[0].ImageLabel)
				}
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: deleteTemplate(),
	}))

	t.Run("Check", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: createTemplate(
			"public-ipv4-count=5 volumes.0.volume-type=sbs volumes.0.size=300GB volumes.0.image-label=ubuntu_jammy volumes.0.name=ubuntu-root-vol server-type=" + serverType1,
		),
		Cmd: "scw instance template check {{ .Template.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: deleteTemplate(),
	}))

	serverType3 := "DEV1-L"
	cmds.Merge(instanceCLIV1.GetCommands())

	t.Run("Create Server", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: createTemplate(
			"public-ipv4-count=1 server-tags.0=cli-tpl server-type=" + serverType3 +
				" volumes.0.volume-type=l_ssd volumes.0.size=20GB volumes.0.image-label=ubuntu_resolute volumes.0.name=cli-tpl",
		),
		Cmd: "scw instance template create-server {{ .Template.ID }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				if server, ok := ctx.Result.(*instanceSDK.Server); !ok {
					t.Errorf("result was not *instance.Server but %T", ctx.Result)
				} else {
					assert.Len(t, server.Tags, 1)
					assert.Equal(t, "cli-tpl", server.Tags[0])
					assert.Len(t, server.PublicNetworkInterface.IPs, 1)
					assert.Equal(t, serverType3, server.ServerType)
					assert.Len(t, server.Volumes, 1)
					assert.Equal(
						t,
						instanceSDK.ServerVolumeVolumeTypeLSSD,
						server.Volumes[0].VolumeType,
					)
				}
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteTemplate(),
			func(ctx *core.AfterFuncCtx) error {
				server := ctx.CmdResult.(*instanceSDK.Server)
				if server.Status == instanceSDK.ServerStatusStarted {
					err := core.ExecAfterCmd("scw instance server stop -w " + server.ID)(ctx)
					if err != nil {
						return err
					}
				}

				return core.ExecAfterCmd(
					"scw instance server delete " + server.ID + " with-ip=true with-volumes=all",
				)(
					ctx,
				)
			},
		),
	}))
}

//
// User Data
//

func setUserData(userDataKey, userDataContent string) core.BeforeFunc {
	return core.ExecBeforeCmd(
		fmt.Sprintf(
			"scw instance template set-user-data {{ .Template.ID }} key=%s content=%s",
			userDataKey,
			userDataContent,
		),
	)
}

func Test_TemplateUserData(t *testing.T) {
	serverType := "GP1-M"

	t.Run("Set User Data", core.Test(&core.TestConfig{
		Commands:   instanceCLIV2.GetCommands(),
		BeforeFunc: createTemplate("server-type=" + serverType),
		Cmd:        "scw instance template set-user-data {{ .Template.ID }} key=example content=put-user-data-here",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				api := instanceSDK.NewAPI(ctx.Client)
				zone, _ := ctx.Client.GetDefaultZone()
				templateID := ctx.Meta["Template"].(*instanceSDK.Template).ID

				tmpl, err := api.GetTemplateUserData(&instanceSDK.GetTemplateUserDataRequest{
					Zone:       zone,
					TemplateID: templateID,
					Key:        "example",
				})
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, []byte("put-user-data-here"), tmpl.Content)
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: deleteTemplate(),
	}))

	t.Run("Get User Data", core.Test(&core.TestConfig{
		Commands: instanceCLIV2.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createTemplate("server-type="+serverType),
			createTmpCloudInitFile(cloudInitScript),
			setUserData("file-key", "@{{ .filePath }}"),
		),
		Cmd: "scw instance template get-user-data {{ .Template.ID }} key=file-key",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				if ud, ok := ctx.Result.(*instanceSDK.UserData); !ok {
					t.Errorf("result was not *instance.UserData but %T", ctx.Result)
				} else {
					assert.Equal(t, cloudInitScript, string(ud.Content))
				}
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteTemplate(),
			closeTemporaryFile(),
		),
	}))

	t.Run("List User Data Keys", core.Test(&core.TestConfig{
		Commands: instanceCLIV2.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createTemplate("server-type="+serverType),
			setUserData("key1", "content1"),
			setUserData("key2", "content2"),
			setUserData("key3", "content3"),
		),
		Cmd: "scw instance template list-user-data-keys {{ .Template.ID }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				if keys, ok := ctx.Result.([]string); !ok {
					t.Errorf("result was not []string but %T", ctx.Result)
				} else {
					assert.Len(t, keys, 3)
					assert.Equal(t, "key1", keys[0])
					assert.Equal(t, "key2", keys[1])
					assert.Equal(t, "key3", keys[2])
				}
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: deleteTemplate(),
	}))

	t.Run("Delete User Data", core.Test(&core.TestConfig{
		Commands: instanceCLIV2.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createTemplate("server-type="+serverType),
			setUserData("key1", "content1"),
			setUserData("key2", "content2"),
			setUserData("key3", "content3"),
		),
		Cmd: "scw instance template delete-user-data {{ .Template.ID }} key=key2",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				api := instanceSDK.NewAPI(ctx.Client)
				zone, _ := ctx.Client.GetDefaultZone()
				templateID := ctx.Meta["Template"].(*instanceSDK.Template).ID

				keys, err := api.ListTemplateUserDataKeys(
					&instanceSDK.ListTemplateUserDataKeysRequest{
						Zone:       zone,
						TemplateID: templateID,
					},
					scw.WithAllPages(),
				)
				if err != nil {
					t.Fatal(err)
				}

				assert.Len(t, keys.Keys, 2)
				assert.Equal(t, "key1", keys.Keys[0])
				assert.Equal(t, "key3", keys.Keys[1])
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: deleteTemplate(),
	}))
}

//
// Cloud Init
//

func setCloudInit() core.BeforeFunc {
	return core.ExecBeforeCmd(
		"scw instance template set-cloud-init {{ .Template.ID }} content=@{{ .filePath }}",
	)
}

func Test_TemplateCloudInit(t *testing.T) {
	serverType := "DEV1-M"

	t.Run("Set Cloud Init", core.Test(&core.TestConfig{
		Commands: instanceCLIV2.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createTemplate("server-type="+serverType),
			createTmpCloudInitFile(cloudInitScript),
		),
		Cmd: `scw instance template set-cloud-init {{ .Template.ID }} content=@{{ .filePath }}`,
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				api := instanceSDK.NewAPI(ctx.Client)
				zone, _ := ctx.Client.GetDefaultZone()
				templateID := ctx.Meta["Template"].(*instanceSDK.Template).ID

				tmpl, err := api.GetTemplateCloudInit(&instanceSDK.GetTemplateCloudInitRequest{
					Zone:       zone,
					TemplateID: templateID,
				})
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, cloudInitScript, string(tmpl.Content))
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			closeTemporaryFile(),
			deleteTemplate(),
		),
	}))

	t.Run("Get Cloud Init", core.Test(&core.TestConfig{
		Commands: instanceCLIV2.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createTemplate("server-type="+serverType),
			createTmpCloudInitFile(cloudInitScript),
			setCloudInit(),
		),
		Cmd: "scw instance template get-cloud-init {{ .Template.ID }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				if ud, ok := ctx.Result.(*instanceSDK.UserData); !ok {
					t.Errorf("result was not *instance.UserData but %T", ctx.Result)
				} else {
					assert.Equal(t, cloudInitScript, string(ud.Content))
				}
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			closeTemporaryFile(),
			deleteTemplate(),
		),
	}))
}

func createTmpCloudInitFile(cloudInit string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		file, err := os.CreateTemp(ctx.T.TempDir(), "test")
		if err != nil {
			ctx.T.Fatalf("%s", err)
		}
		_, err = file.WriteString(cloudInit)
		if err != nil {
			ctx.T.Fatalf("%s", err)
		}

		ctx.Meta["filePath"] = file.Name()
		ctx.Meta["File"] = file

		return nil
	}
}

func closeTemporaryFile() core.AfterFunc {
	return func(ctx *core.AfterFuncCtx) error {
		file := ctx.Meta["File"].(*os.File)

		// We need to close this file explicitly because it is not closed by the os.Remove call on windows
		// https://github.com/golang/go/issues/50510
		err := file.Close()
		if err != nil {
			return err
		}

		err = os.RemoveAll(file.Name())
		if err != nil {
			return err
		}

		return nil
	}
}
