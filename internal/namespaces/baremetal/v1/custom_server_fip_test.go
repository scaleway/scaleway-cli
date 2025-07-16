package baremetal_test

import (
	"errors"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/baremetal/v1"
	flexibleip "github.com/scaleway/scaleway-cli/v2/internal/namespaces/flexibleip/v1alpha1"
	baremetalSDK "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func Test_CreateFlexibleIPInteractive(t *testing.T) {
	promptResponse := []string{
		`" "`,
	}
	interactive.IsInteractive = true
	cmds := baremetal.GetCommands()
	cmds.Merge(flexibleip.GetCommands())
	t.Run("Simple Interactive", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			func(ctx *core.BeforeFuncCtx) error {
				api := baremetalSDK.NewAPI(ctx.Client)
				server, _ := api.GetOfferByName(&baremetalSDK.GetOfferByNameRequest{
					OfferName: offerNameNVME,
					Zone:      scw.Zone(zone),
				})
				if server.Stock != baremetalSDK.OfferStockAvailable {
					return errors.New("offer out of stock")
				}

				return nil
			},
			createServerAndWait(),
		),
		Cmd: "scw baremetal server add-flexible-ip {{ .Server.ID }} zone=" + zone,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			core.ExecAfterCmd("scw fip ip delete {{ .CmdResult.ID }} zone="+zone),
		),
		PromptResponseMocks: promptResponse,
	}))
}

func Test_CreateFlexibleIP(t *testing.T) {
	interactive.IsInteractive = false
	cmds := baremetal.GetCommands()
	cmds.Merge(flexibleip.GetCommands())
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			func(ctx *core.BeforeFuncCtx) error {
				api := baremetalSDK.NewAPI(ctx.Client)
				server, _ := api.GetOfferByName(&baremetalSDK.GetOfferByNameRequest{
					OfferName: offerNameNVME,
					Zone:      scw.Zone(zone),
				})
				if server.Stock != baremetalSDK.OfferStockAvailable {
					return errors.New("offer out of stock")
				}

				return nil
			},
			createServerAndWait(),
		),
		Cmd: "scw baremetal server add-flexible-ip {{ .Server.ID }} ip-type=IPv4 zone=" + zone,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			core.ExecAfterCmd("scw fip ip delete {{ .CmdResult.ID }} zone="+zone),
		),
	}))
}
