package instance

import "github.com/scaleway/scaleway-cli/internal/core"

// Builders

func ipCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("organization").Name = "organization-id"
	return c
}

func ipListBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("organization").Name = "organization-id"
	return c
}
