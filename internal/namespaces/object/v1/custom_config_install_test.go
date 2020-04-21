package object

import (
	"testing"
)

func Test_ConfigInstall(t *testing.T) {
	t.Run("rclone", func(t *testing.T) {

		//t.Run("simple", core.Test(&core.TestConfig{
		//	Commands: GetCommands(),
		//	Cmd:      "scw object config install type=rclone",
		//	Check: core.TestCheckCombine(
		//		// no golden tests since it's os specific
		//		func(t *testing.T, ctx *core.CheckFuncCtx) {
		//			testIfKubeconfigInFile(t, path.Join(os.TempDir(), "cli-test"), "-"+ctx.Meta["Cluster"].(*k8s.Cluster).ID, ctx.Meta["Kubeconfig"].(*k8s.Kubeconfig))
		//		},
		//		core.TestCheckExitCode(0),
		//	),
		//}))
	})

	t.Run("mc", func(t *testing.T) {

	})

	t.Run("s3cmd", func(t *testing.T) {

	})
}
