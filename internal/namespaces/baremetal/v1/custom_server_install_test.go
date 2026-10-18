package baremetal_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/baremetal/v1"
	iam "github.com/scaleway/scaleway-cli/v2/internal/namespaces/iam/v1alpha1"
)

func Test_InstallServer(t *testing.T) {
	// All test below should succeed to create an instance.
	t.Run("Simple", func(t *testing.T) {
		// baremetal api requires that the key must be at least 1024 bits long. Regardless of the algorithm
		sshKey := `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQCbJuYSOQc01zjHsMyn4OUsW61cqRvttKt3StJgbvt2WBuGpwi1/5RtSoMQpudYlZpdeivFb21S8QRas8zcOc+6WqgWa2nj/8yA+cauRlV6CMWY+hOTkkg39xaekstuQ+WR2/AP7O/9hjVx5735+9ZNIxxHsFjVYdBEuk9gEX+1Rw== foobar@foobar`
		osID := `b9a016fc-947f-4bbc-bf95-343c8524535a`
		cmds := baremetal.GetCommands()
		cmds.Merge(iam.GetCommands())

		t.Run("With ID", core.Test(&core.TestConfig{
			BeforeFunc: core.BeforeFuncCombine(
				selectOffer(
					zone,
					offerFilterAvailable,
					offerFilterNVMe,
					offerFilterCompatibleOS(osID),
				),
				addSSH(sshKey),
				createServerAndWait(),
			),
			Commands: cmds,
			Cmd:      "scw baremetal server install {{ .Server.ID }} zone={{ .Server.Zone }} hostname=test-install-server ssh-key-ids.0={{ .Key.ID }} os-id=" + osID + " -w",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
			AfterFunc: core.AfterFuncCombine(
				deleteSSH(),
				deleteServer(),
			),
		}))

		t.Run("All SSH keys", core.Test(&core.TestConfig{
			Commands: cmds,
			BeforeFunc: core.BeforeFuncCombine(
				selectOffer(
					zone,
					offerFilterAvailable,
					offerFilterNVMe,
					offerFilterCompatibleOS(osID),
				),
				addSSH(sshKey),
				createServerAndWait(),
			),
			Cmd: "scw baremetal server install {{ .Server.ID }} zone={{ .Server.Zone }} hostname=test-install-server all-ssh-keys=true os-id=" + osID + " -w",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
			AfterFunc: core.AfterFuncCombine(
				deleteSSH(),
				deleteServer(),
			),
		}))
	})
}
