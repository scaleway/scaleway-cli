package instance

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func completeServerType(
	ctx context.Context,
	prefix string,
	createReq any,
) core.AutocompleteSuggestions {
	req := createReq.(*instanceCreateServerRequest)
	resp, err := instance.NewAPI(core.ExtractClient(ctx)).
		ListServersTypes(&instance.ListServersTypesRequest{
			Zone: req.Zone,
		}, scw.WithAllPages())
	if err != nil {
		return nil
	}

	suggestions := make([]string, 0, len(resp.Servers))

	for serverType := range resp.Servers {
		if strings.HasPrefix(serverType, prefix) {
			suggestions = append(suggestions, serverType)
		}
	}

	return suggestions
}

func commercialTypeIsWindowsServer(commercialType string) bool {
	return strings.HasSuffix(commercialType, "-WIN")
}

func SizeValue(s *scw.Size) scw.Size {
	if s != nil {
		return *s
	}

	return 0
}

func warningServerTypeDeprecated(
	ctx context.Context,
	client *scw.Client,
	server *instance.Server,
) []string {
	warning := []string{
		terminal.Style(
			fmt.Sprintf(
				"Warning: server type %q will soon reach EndOfService",
				server.CommercialType,
			),
			color.Bold,
			color.FgYellow,
		),
	}

	compatibleTypes, err := instance.NewAPI(client).
		GetServerCompatibleTypes(&instance.GetServerCompatibleTypesRequest{
			Zone:     server.Zone,
			ServerID: server.ID,
		}, scw.WithContext(ctx))
	if err != nil {
		return warning
	}

	mostRelevantTypes := compatibleTypes.CompatibleTypes[:5]
	details := fmt.Sprintf(`
	Your Instance will soon reach End of Service. You can check the exact date on the Scaleway console. We recommend that you migrate your Instance before that.
	Here are the %d best options for %q, ordered by relevance: [%s]
	You can check the full list of compatible server types:
		- on the Scaleway console
		- using the CLI command 'scw instance server get-compatible-types %s zone=%s'`,
		len(mostRelevantTypes),
		server.CommercialType,
		strings.Join(mostRelevantTypes, ", "),
		server.ID,
		server.Zone,
	)

	return append(warning, details)
}
