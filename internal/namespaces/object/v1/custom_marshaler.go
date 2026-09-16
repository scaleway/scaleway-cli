package object

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type CustomS3ACLGrant struct {
	Grantee    *string
	Permission types.Permission
}

type BucketResponse struct {
	SuccessResult *core.SuccessResult
	BucketInfo    *bucketInfo
}

func bucketResponseMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	resp := i.(BucketResponse)

	messageStr, err := resp.SuccessResult.MarshalHuman()
	if err != nil {
		return "", err
	}
	bucketStr, err := bucketInfoMarshalerFunc(*resp.BucketInfo, opt)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{
		messageStr,
		bucketStr,
	}, "\n"), nil
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

func bucketInfoMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	// To avoid recursion of human.Marshal we create a dummy type
	type tmp bucketInfo
	info := tmp(i.(bucketInfo))

	opt.Sections = []*human.MarshalSection{
		{
			FieldName:   "Tags",
			HideIfEmpty: true,
		},
		{
			FieldName:   "ACL",
			HideIfEmpty: true,
		},
	}
	str, err := human.Marshal(info, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

type BucketGetResult struct {
	*bucketInfo
	Size      *scw.Size
	NbObjects *int64
	NbParts   *int64
}

func bucketGetResultMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp BucketGetResult
	result := tmp(i.(BucketGetResult))
	opt.Sections = []*human.MarshalSection{
		{
			FieldName:   "Tags",
			HideIfEmpty: true,
		},
		{
			FieldName:   "ACL",
			HideIfEmpty: true,
		},
	}
	str, err := human.Marshal(result, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

type bucketGetArgs struct {
	Region    scw.Region
	Name      string
	WithSize  bool `json:"with-size"`
	ProjectID string
}

func bucketMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp []types.Bucket
	result := tmp(i.([]types.Bucket))
	opt.Fields = []*human.MarshalFieldOpt{
		{
			FieldName: "Name",
			Label:     "Name",
		},
		{
			FieldName: "CreationDate",
			Label:     "Creation Date",
		},
	}
	str, err := human.Marshal(result, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

// CleanAndIndentJSON takes any struct or data, marshals it,
// strips all top-level and nested "null" fields, and returns formatted JSON.
func CleanAndIndentJSON(v any, prefix, indent string) ([]byte, error) {
	rawBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var genericData any
	if err := json.Unmarshal(rawBytes, &genericData); err != nil {
		return nil, err
	}

	cleanedData := removeNulls(genericData)

	cleanBytes, err := json.Marshal(cleanedData)
	if err != nil {
		return nil, err
	}

	var prettyBuf bytes.Buffer
	if err := json.Indent(&prettyBuf, cleanBytes, prefix, indent); err != nil {
		return nil, err
	}

	return prettyBuf.Bytes(), nil
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
