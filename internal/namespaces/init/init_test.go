package init

//func TestInit(t *testing.T) {
//	t.Run("Simple", func(t *testing.T) {
//		ctx := context.Background()
//		// httpClient, cleanup := createHttpRecoder(t, cassetName)
//		// defer cleanup()
//		// ctx = accont.InjectHttpClient(httpClient)
//
//		//ctx = interactive.InjectMockResponseToContext(ctx, []string{
//		//	"yes",
//		//	"no",
//		//})
//
//		core.Test(&core.TestConfig{
//			Ctx: ctx,
//			PromptResponseMocks: []string{
//				"yes",
//				"no",
//			},
//			TmpHomeDir: true,
//			Cmd:        "scw -c {{ .TmpHomeDir }}/config.yml init",
//			BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
//				//homeDir := ctx.Meta["TmpHomeDir"].(string)
//				// Write fake config file
//				return nil
//			},
//			Check: core.TestCheckCombine(
//				core.TestCheckExitCode(0),
//				core.TestCheckGolden(),
//				func(t *testing.T, ctx *core.CheckFuncCtx) {
//					//homeDir := ctx.Meta["TmpHomeDir"].(string)
//					// Check config file is correct
//				}),
//		})(t)
//
//	})
//}
