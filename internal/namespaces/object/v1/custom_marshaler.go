package object

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// S3Args holds the fields shared by every bucket-related arg struct.
// Required here for the WASM build.
type S3Args struct {
	Region         scw.Region
	ProjectID      string
	S3Endpoint     string
	S3UsePathStyle string
}

type CustomS3ACLGrant struct {
	Grantee    *string
	Permission types.Permission
}

type BucketResponse struct {
	SuccessResult *core.SuccessResult
	BucketInfo    *bucketInfo
}

type bucketInfo struct {
	ID               string
	Region           scw.Region
	APIEndpoint      string
	BucketEndpoint   string
	EnableVersioning bool
	Tags             []types.Tag
	ACL              []CustomS3ACLGrant
	Owner            string
}

type BucketGetResult struct {
	*bucketInfo
	Size      *scw.Size
	NbObjects *int64
	NbParts   *int64
}

type bucketGetArgs struct {
	S3Args
	Name     string
	WithSize bool `json:"with-size"`
}

var bucketInfoLayout = humanv2.Layout{
	Sections: []humanv2.Section{
		{Path: "Tags", Title: "Tags", Omittable: true},
		{Path: "ACL", Title: "Acl", Omittable: true},
	},
}

var bucketGetResultLayout = humanv2.Layout{Sections: bucketInfoLayout.Sections}

func bucketResponseView(v BucketResponse) (string, error) {
	messageStr, err := humanv2.Format(v.SuccessResult)
	if err != nil {
		return "", err
	}
	bucketStr, err := humanv2.Render(*v.BucketInfo, bucketInfoLayout)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{messageStr, bucketStr}, "\n"), nil
}

// CleanAndIndentJSON takes any struct or data, marshals it,
// strips all top-level and nested "null" fields, and returns formatted JSON.
var awsMetadataMarshalers = json.WithMarshalers(
	json.MarshalFunc(func(middleware.Metadata) ([]byte, error) {
		return []byte("{}"), nil
	}),
)

func CleanAndIndentJSON(v any, prefix, indent string) ([]byte, error) {
	rawBytes, err := json.Marshal(v, awsMetadataMarshalers, json.Deterministic(true))
	if err != nil {
		return nil, err
	}

	var genericData any
	if err := json.Unmarshal(rawBytes, &genericData); err != nil {
		return nil, err
	}

	cleanedData := removeNulls(genericData)

	return json.Marshal(
		cleanedData,
		json.Deterministic(true),
		jsontext.WithIndentPrefix(prefix),
		jsontext.WithIndent(indent),
	)
}

// removeNulls recursively purges nil entries from maps and slices
func removeNulls(v any) any {
	switch val := v.(type) {
	case map[string]any:
		cleanedMap := make(map[string]any)

		for k, item := range val {
			if item == nil {
				continue // Skip "null" fields
			} else if item == "" {
				continue // Skip empty strings
			}

			cleanedItem := removeNulls(item)
			if cleanedItem == nil {
				continue // Skip empty structs
			}

			cleanedMap[k] = removeNulls(item)
		}

		if len(cleanedMap) == 0 {
			return nil // Skip empty maps
		}

		return cleanedMap

	case []any:
		cleanedSlice := make([]any, 0, len(val))

		for _, item := range val {
			if item != nil {
				cleanedSlice = append(cleanedSlice, removeNulls(item))
			}
		}

		if len(cleanedSlice) == 0 {
			return nil // Skip empty slices
		}

		return cleanedSlice

	default:
		return v
	}
}
