package instance

import "github.com/scaleway/scaleway-cli/v2/core"

func placementGroupUpdateServersBuilder(c *core.Command) *core.Command {
	return deprecateV1OnlyCommand(c)
}

func placementGroupGetServersBuilder(c *core.Command) *core.Command {
	return deprecateV1OnlyCommand(c)
}

func placementGroupSetServersBuilder(c *core.Command) *core.Command {
	return deprecateV1OnlyCommand(c)
}

func placementGroupSetBuilder(c *core.Command) *core.Command {
	return deprecateV1OnlyCommand(c)
}

func deprecateV1OnlyCommand(c *core.Command) *core.Command {
	c.Deprecated = true
	c.DeprecationMessage = "it is not available in v2 of the API and will be removed when v1 reaches end of support."

	return c
}
