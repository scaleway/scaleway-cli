//go:build darwin || linux || windows

package object

import "github.com/scaleway/scaleway-cli/v2/core"

func objectRoot() *core.Command {
	return &core.Command{
		Short:     `Object-storage utils`,
		Namespace: "object",
		Groups:    []string{"storage"},
	}
}

func objectConfig() *core.Command {
	return &core.Command{
		Short:     `Manage configuration files for popular S3 tools`,
		Long:      `Configuration generation for S3 tools.`,
		Namespace: "object",
		Resource:  `config`,
	}
}

func objectBucket() *core.Command {
	return &core.Command{
		Short:     `Manage S3 buckets`,
		Long:      `Manage S3 buckets creation, deletion and updates to properties like tags, ACL and versioning.`,
		Namespace: "object",
		Resource:  `bucket`,
	}
}

func objectBucketLifecycle() *core.Command {
	return &core.Command{
		Short:     `Manage S3 buckets' lifecycle configuration`,
		Long:      `Manage S3 buckets lifecycle rules creation and deletion.`,
		Namespace: "object",
		Resource:  `bucket-lifecycle`,
  }
}

func objectBucketCors() *core.Command {
	return &core.Command{
		Short:     `Manage S3 bucket CORS`,
		Long:      `Manage S3 bucket CORS rules creation and deletion.`,
		Namespace: `object`,
		Resource:  `bucket-cors`,
	}
}

func objectBucketPolicy() *core.Command {
	return &core.Command{
		Short:     `Manage S3 bucket policies`,
		Long:      `Manage S3 bucket policies creation and deletion.`,
		Namespace: "object",
		Resource:  `bucket-policy`,
	}
}
