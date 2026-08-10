package product_catalog

import "github.com/scaleway/scaleway-cli/v2/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("product-catalog", "product", "list").Override(productListBuilder)

	return cmds
}

func productListBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("product-types.{index}").AutoCompleteFunc = autocompleteProductType

	return c
}
