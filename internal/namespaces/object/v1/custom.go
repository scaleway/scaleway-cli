package object

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
)

func GetCommands() *core.Commands {
	human.RegisterMarshalerFunc(BucketResponse{}, bucketResponseMarshalerFunc)
	human.RegisterMarshalerFunc(bucketInfo{}, bucketInfoMarshalerFunc)
	human.RegisterMarshalerFunc(BucketGetResult{}, bucketGetResultMarshalerFunc)

	return core.NewCommands(
		objectRoot(),
		objectConfig(),
		objectBucket(),
		bucketCreateCommand(),
		bucketDeleteCommand(),
		bucketGetCommand(),
		bucketListCommand(),
		bucketUpdateCommand(),
		configGetCommand(),
		configInstallCommand(),
	)
}

func objectRoot() *core.Command {
	return &core.Command{
		Short:     `Object-storage utils`,
		Namespace: "object",
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
