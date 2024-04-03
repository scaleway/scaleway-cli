package secret

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	secret "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
)

type customAccessSecretVersionRequest struct {
	secret.AccessSecretVersionRequest
	Field *string
	Raw   bool
}

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()
	cmds.MustFind("secret", "version", "create").Override(secretVersionCreateData)
	cmds.MustFind("secret", "version", "access").Override(secretVersionAccessCommand)
	return cmds
}

func secretVersionCreateData(c *core.Command) *core.Command {
	*c.ArgSpecs.GetByName("data") = core.ArgSpec{
		Name:        "data",
		Short:       "Content of the secret version. Base64 is handled by the SDK",
		Required:    true,
		CanLoadFile: true,
	}
	return c
}

func secretVersionAccessCommand(c *core.Command) *core.Command {
	c.ArgsType = reflect.TypeOf(customAccessSecretVersionRequest{})

	c.ArgSpecs.AddBefore("region", &core.ArgSpec{
		Name:  "field",
		Short: "Return only the JSON field of the given name",
	})

	c.ArgSpecs.AddBefore("region", &core.ArgSpec{
		Name:  "raw",
		Short: "Return only the raw payload",
	})

	c.Run = func(ctx context.Context, args interface{}) (i interface{}, e error) {
		client := core.ExtractClient(ctx)
		api := secret.NewAPI(client)

		request := args.(*customAccessSecretVersionRequest)

		response, err := api.AccessSecretVersion(&secret.AccessSecretVersionRequest{
			Region:   request.Region,
			SecretID: request.SecretID,
			Revision: request.Revision,
		})
		if err != nil {
			return nil, err
		}

		if request.Field != nil {
			response.Data, err = getSecretVersionField(response.Data, *request.Field)
			if err != nil {
				return nil, err
			}
		}

		if request.Raw {
			return core.RawResult(response.Data), nil
		}

		return response, nil
	}

	return c
}

func getSecretVersionField(data []byte, field string) ([]byte, error) {
	var rawFields interface{}
	if err := json.Unmarshal(data, &rawFields); err != nil {
		return nil, fmt.Errorf("cannot unmarshal JSON data")
	}

	rawField, ok := rawFields.(map[string]interface{})[field]
	if !ok {
		return nil, fmt.Errorf("JSON field is not present")
	}

	switch field := rawField.(type) {
	case string:
		return []byte(field), nil
	default:
		return nil, fmt.Errorf("JSON field type is not valid")
	}
}
