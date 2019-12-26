package instance

import (
	"context"
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

func applyCustomErrors(c *core.Commands) {
	c.MustFind("instance", "security-group", "delete").ConvertErrorFunc = customInstanceSecurityGroupDeleteConvertErrorFunc
}

func customInstanceSecurityGroupDeleteConvertErrorFunc(ctx context.Context, argsI interface{}, currentError error) error {
	e := strings.ToLower(currentError.Error())

	switch {
	case strings.HasSuffix(e, "group is in use. you cannot delete it."):
		req := argsI.(*instance.DeleteSecurityGroupRequest)
		api := instance.NewAPI(core.ExtractClient(ctx))

		newError := &core.CliError{
			Err: fmt.Errorf("cannot delete in use security-group"),
		}

		// Get security-group.
		sg, err := api.GetSecurityGroup(&instance.GetSecurityGroupRequest{
			SecurityGroupID: req.SecurityGroupID,
		})
		if err != nil {
			// Ignore API error and return a minimal error.
			return newError
		}

		// Create detail message.
		details := "Attach all these instances to another security-group before deleting this one:"
		for _, s := range sg.SecurityGroup.Servers {
			details += "\n  - scw instance server update server-id=" + s.ID + " security-group.id=$NEW_SECURITY_GROUP_ID"
		}

		newError.Details = details
		return newError

	default:
		return currentError
	}
}
