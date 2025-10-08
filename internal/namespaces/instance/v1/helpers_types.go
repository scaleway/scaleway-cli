package instance

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	block "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	product_catalog "github.com/scaleway/scaleway-sdk-go/api/product_catalog/v2alpha1"
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

func volumeIsFromSBS(api *block.API, zone scw.Zone, volumeID string) bool {
	_, err := api.GetVolume(&block.GetVolumeRequest{
		Zone:     zone,
		VolumeID: volumeID,
	})

	return err == nil
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

	eosDate, err := getEndOfServiceDate(ctx, client, server.Zone, server.CommercialType)
	if err != nil {
		return warning
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
	Your Instance will reach End of Service by %s. You can check the exact date on the Scaleway console. We recommend that you migrate your Instance before that.
	Here are the %d best options for %q, ordered by relevance: [%s]
	You can check the full list of compatible server types:
		- on the Scaleway console
		- using the CLI command 'scw instance server get-compatible-types %s zone=%s'`,
		eosDate,
		len(mostRelevantTypes),
		server.CommercialType,
		strings.Join(mostRelevantTypes, ", "),
		server.ID,
		server.Zone,
	)

	return append(warning, details)
}

func getEndOfServiceDate(
	ctx context.Context,
	client *scw.Client,
	zone scw.Zone,
	commercialType string,
) (string, error) {
	api := product_catalog.NewPublicCatalogAPI(client)

	products, err := api.ListPublicCatalogProducts(
		&product_catalog.PublicCatalogAPIListPublicCatalogProductsRequest{
			ProductTypes: []product_catalog.ListPublicCatalogProductsRequestProductType{
				product_catalog.ListPublicCatalogProductsRequestProductTypeInstance,
			},
		},
		scw.WithAllPages(),
		scw.WithContext(ctx),
	)
	if err != nil {
		return "", fmt.Errorf("could not list product catalog entries: %w", err)
	}

	for _, product := range products.Products {
		if product.Properties != nil && product.Properties.Instance != nil &&
			product.Properties.Instance.OfferID == commercialType {
			return product.EndOfLifeAt.Format(time.DateOnly), nil
		}
	}

	return "", fmt.Errorf("could not find product catalog entry for %q in %s", commercialType, zone)
}
