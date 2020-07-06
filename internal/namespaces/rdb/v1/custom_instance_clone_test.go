package rdb

//func Test_CloneInstance(t *testing.T) {
//	// Simple use cases
//	t.Run("Simple", func(t *testing.T) {
//		t.Run("Default", core.Test(&core.TestConfig{
//			Commands: GetCommands(),
//			Cmd:      "scw rdb instance clone -w",
//			Check: core.TestCheckCombine(
//				core.TestCheckGolden(),
//				core.TestCheckExitCode(0),
//			),
//			AfterFunc: func(ctx *core.AfterFuncCtx) error {
//				_, err := baremetal.NewAPI(ctx.Client).DeleteServer(&baremetal.DeleteServerRequest{
//					ServerID: ctx.CmdResult.(*baremetal.Server).ID,
//				})
//				if err != nil {
//					return err
//				}
//				return nil
//			},
//			DefaultZone: scw.ZoneFrPar2,
//		}))
//	}
//
//}
