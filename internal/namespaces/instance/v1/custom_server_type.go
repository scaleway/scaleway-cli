package instance

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func serverTypeGetBuilder(c *core.Command) *core.Command {
	deprecateV1OnlyCommand(c)
	c.DeprecationMessage += " Availability can be fetched using \"scw instance server-type list\"."

	return c
}
