package instance

import "github.com/scaleway/scaleway-cli/internal/core"

// Builders

func snapshotCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName(oldOrganizationFieldName).Name = newOrganizationFieldName
	return c
}

func snapshotListBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName(oldOrganizationFieldName).Name = newOrganizationFieldName
	return c
}
