package instance

import (
	"context"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
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
