package feedback

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
)

func GetCommands() *core.Commands {
	return core.NewCommands(
		feedbackRoot(),
		feedbackBugCommand(),
		feedbackFeatureRequestCommand(),
	)
}

func feedbackRoot() *core.Command {
	return &core.Command{
		Short:     `Send feedback to the Scaleway CLI Team!`,
		Namespace: "feedback",
		ArgsType:  reflect.TypeOf(struct{}{}),
		ArgSpecs:  core.ArgSpecs{},
	}
}

func feedbackBugCommand() *core.Command {
	return &core.Command{
		Short:     `Send a bug-report`,
		Long:      `Send a bug-report to the Scaleway CLI team.`,
		Namespace: "feedback",
		Resource:  `bug`,
		ArgsType:  reflect.TypeOf(struct{}{}),
		ArgSpecs:  core.ArgSpecs{},

		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			buildInfoStr, err := human.Marshal(core.ExtractBuildInfo(ctx), nil)
			if err != nil {
				return nil, err
			}
			issue := issue{
				IssueTemplate: bug,
				BuildInfoStr:  buildInfoStr,
			}
			err = issue.openInBrowser(ctx)
			if err != nil {
				log.Fatal(err)
			}
			return &core.SuccessResult{
				Message:  "Successfully opened the page",
				Details: fmt.Sprintf(" (%s)", issue.getURL()),
			}, nil
		},

	}
}

func feedbackFeatureRequestCommand() *core.Command {
	return &core.Command{
		Short:     `Send a feature request`,
		Long:      `Send a feature request to the Scaleway CLI team.`,
		Namespace: "feedback",
		Resource:  `feature`,

		ArgsType:  reflect.TypeOf(struct{}{}),
		ArgSpecs:  core.ArgSpecs{},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			buildInfoStr, err := human.Marshal(core.ExtractBuildInfo(ctx), nil)
			if err != nil {
				return nil, err
			}
			issue := issue{
				IssueTemplate: feature,
				BuildInfoStr: buildInfoStr,
			}
			err = issue.openInBrowser(ctx)
			if err != nil {
				log.Fatal(err)
			}
			return &core.SuccessResult{
				Message:  "Successfully opened the page",
				Details: fmt.Sprintf(" (%s)", issue.getURL()),
			}, nil
		},
	}
}

