package instance

import "github.com/scaleway/scaleway-cli/internal/core"

// Builders

func snapshotCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("organization").Name = "organization-id"
	return c
}

func snapshotListBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("organization").Name = "organization-id"
	return c
}
