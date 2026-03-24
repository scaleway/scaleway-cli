package product_catalog

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	product_catalog "github.com/scaleway/scaleway-sdk-go/api/product_catalog/v2alpha1"
)

// Caching ListPublicCatalogProductsRequestProductType values for shell completion
var completeProductTypeCache []product_catalog.ListPublicCatalogProductsRequestProductType

func autocompleteProductType(ctx context.Context, _ string, _ any) core.AutocompleteSuggestions {
	suggestions := core.AutocompleteSuggestions(nil)

	if len(completeProductTypeCache) == 0 {
		var productTypes product_catalog.ListPublicCatalogProductsRequestProductType
		completeProductTypeCache = productTypes.Values()
	}

	for _, productType := range completeProductTypeCache {
		suggestions = append(suggestions, string(productType))
	}

	return suggestions
}
