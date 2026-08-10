package secret

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	secret "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
)

func secretVersionCreateBuilder(c *core.Command) *core.Command {
	*c.ArgSpecs.GetByName("data") = core.ArgSpec{
		Name:        "data",
		Short:       "Content of the secret version.",
		Required:    true,
		CanLoadFile: true,
	}

	c.Examples = append(c.Examples, &core.Example{
		Short: "Create a json secret version",
		Raw:   "scw secret version create 11111111-1111-1111-111111111111 data={\"key\":\"value\"}",
	})

	return c
}

func secretVersionAccessBuilder(c *core.Command) *core.Command {
	c.ArgsType = reflect.TypeOf(customAccessSecretVersionRequest{})

	c.ArgSpecs.AddBefore("region", &core.ArgSpec{
		Name:  "field",
		Short: "Return only the JSON field of the given name",
	})

	c.ArgSpecs.AddBefore("region", &core.ArgSpec{
		Name:  "raw",
		Short: "Return only the raw payload",
	})

	c.Examples = append(c.Examples, &core.Example{
		Short: "Get a raw json value from a secret version",
		Raw:   "scw secret version access 11111111-1111-1111-111111111111 revision=1 field=key raw=true",
	})

	c.Run = func(ctx context.Context, args any) (i any, e error) {
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
	var rawFields any
	if err := json.Unmarshal(data, &rawFields); err != nil {
		return nil, errors.New("cannot unmarshal JSON data")
	}

	rawField, ok := rawFields.(map[string]any)[field]
	if !ok {
		return nil, errors.New("JSON field is not present")
	}

	switch field := rawField.(type) {
	case string:
		return []byte(field), nil
	default:
		return nil, errors.New("JSON field type is not valid")
	}
}
