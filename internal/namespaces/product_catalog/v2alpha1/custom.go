package product_catalog

import "github.com/scaleway/scaleway-cli/v2/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	listCmd := cmds.MustFind("product-catalog", "product", "list")
	listCmd.Override(productListBuilder)
	listCmd.Override(core.ListViewBuilder("Sku", "ServiceCategory", "ProductCategory", "Product", "Variant", "Description", "Status", "EndOfLifeAt", "EndOfGrowthAt", "EndOfSaleAt", "Badges"))

	return cmds
}

func productListBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("product-types.{index}").AutoCompleteFunc = autocompleteProductType

	return c
}
