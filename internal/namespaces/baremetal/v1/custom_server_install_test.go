package baremetal_test

import (
	"errors"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/baremetal/v1"
	iam "github.com/scaleway/scaleway-cli/v2/internal/namespaces/iam/v1alpha1"
	baremetalSDK "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func Test_InstallServer(t *testing.T) {
	// All test below should succeed to create an instance.
	t.Run("Simple", func(t *testing.T) {
		// baremetal api requires that the key must be at least 1024 bits long. Regardless of the algorithm
		sshKey := `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQCbJuYSOQc01zjHsMyn4OUsW61cqRvttKt3StJgbvt2WBuGpwi1/5RtSoMQpudYlZpdeivFb21S8QRas8zcOc+6WqgWa2nj/8yA+cauRlV6CMWY+hOTkkg39xaekstuQ+WR2/AP7O/9hjVx5735+9ZNIxxHsFjVYdBEuk9gEX+1Rw== foobar@foobar`
		osID := `03b7f4ba-a6a1-4305-984e-b54fafbf1681` // Ubuntu 20.04 LTS (Focal)
		cmds := baremetal.GetCommands()
		cmds.Merge(iam.GetCommands())

		t.Run("With ID", core.Test(&core.TestConfig{
			BeforeFunc: core.BeforeFuncCombine(
				func(ctx *core.BeforeFuncCtx) error {
					api := baremetalSDK.NewAPI(ctx.Client)
					server, _ := api.GetOfferByName(&baremetalSDK.GetOfferByNameRequest{
						OfferName: offerNameNVME,
						Zone:      scw.Zone(zone),
					})
					if server.Stock != baremetalSDK.OfferStockAvailable {
						err := errors.New("offer out of stock")

						return err
					}

					return nil
				},
				addSSH("key", sshKey),
				createServerAndWait(),
			),
			Commands: cmds,
			Cmd:      "scw baremetal server install {{ .Server.ID }} zone=" + zone + " hostname=test-install-server ssh-key-ids.0={{ .key.ID }} os-id=" + osID + " -w",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
			AfterFunc: core.AfterFuncCombine(
				deleteSSH("key"),
				deleteServer("Server"),
			),
		}))

		t.Run("All SSH keys", core.Test(&core.TestConfig{
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
				addSSH("key", sshKey),
				createServerAndWait(),
			),
			Cmd: "scw baremetal server install {{ .Server.ID }} zone=" + zone + " hostname=test-install-server all-ssh-keys=true os-id=" + osID + " -w",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
			AfterFunc: core.AfterFuncCombine(
				deleteSSH("key"),
				deleteServer("Server"),
			),
		}))
	})
}
