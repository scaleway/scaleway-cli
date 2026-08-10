package feedback

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
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
		Groups:    []string{"utility"},
		Short:     `Send feedback to the Scaleway CLI Team!`,
		Namespace: "feedback",
		ArgsType:  reflect.TypeOf(struct{}{}),
		ArgSpecs:  core.ArgSpecs{},
	}
}

func feedbackBugCommand() *core.Command {
	return &core.Command{
		Groups:               []string{"utility"},
		Short:                `Send a bug-report`,
		Long:                 `Send a bug-report to the Scaleway CLI team.`,
		Namespace:            "feedback",
		Resource:             `bug`,
		ArgsType:             reflect.TypeOf(struct{}{}),
		ArgSpecs:             core.ArgSpecs{},
		AllowAnonymousClient: true,

		Run: func(ctx context.Context, _ any) (i any, e error) {
			issue := issue{
				IssueTemplate: bug,
				BuildInfo:     core.ExtractBuildInfo(ctx),
			}
			err := issue.openInBrowser(ctx)
			if err != nil {
				return nil, err
			}

			return &core.SuccessResult{
				Message: "Successfully opened the page",
				Details: issue.getURL(),
			}, nil
		},
	}
}

func feedbackFeatureRequestCommand() *core.Command {
	return &core.Command{
		Groups:               []string{"utility"},
		Short:                `Send a feature request`,
		Long:                 `Send a feature request to the Scaleway CLI team.`,
		Namespace:            "feedback",
		Resource:             `feature`,
		ArgsType:             reflect.TypeOf(struct{}{}),
		ArgSpecs:             core.ArgSpecs{},
		AllowAnonymousClient: true,
		Run: func(ctx context.Context, _ any) (i any, e error) {
			issue := issue{
				IssueTemplate: feature,
				BuildInfo:     core.ExtractBuildInfo(ctx),
			}
			err := issue.openInBrowser(ctx)
			if err != nil {
				return nil, err
			}

			return &core.SuccessResult{
				Message: "Successfully opened the page",
				Details: issue.getURL(),
			}, nil
		},
	}
}
