package product_catalog_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/product_catalog/v2alpha1"
)

func Test_ListProduct(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: product_catalog.GetCommands(),
		Cmd:      "scw product-catalog product list",
		Check:    core.TestCheckGolden(),
	}))
}
