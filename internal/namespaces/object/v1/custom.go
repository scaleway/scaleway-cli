package object

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
)

func GetCommands() *core.Commands {
	human.RegisterMarshalerFunc(BucketResponse{}, bucketResponseMarshalerFunc)
	human.RegisterMarshalerFunc(bucketInfo{}, bucketInfoMarshalerFunc)
	human.RegisterMarshalerFunc(BucketGetResult{}, bucketGetResultMarshalerFunc)
	human.RegisterMarshalerFunc(s3.ListBucketsOutput{}.Buckets, bucketMarshalerFunc)

	cmds := core.NewCommands()

	if cmdObjectRoot := objectRoot(); cmdObjectRoot != nil {
		cmds.Add(cmdObjectRoot)
	}
	if cmdObjectConfig := objectConfig(); cmdObjectConfig != nil {
		cmds.Add(cmdObjectConfig)
	}
	if cmdObjectBucket := objectBucket(); cmdObjectBucket != nil {
		cmds.Add(cmdObjectBucket)
	}
	if cmdBucketCreate := bucketCreateCommand(); cmdBucketCreate != nil {
		cmds.Add(cmdBucketCreate)
	}
	if cmdBucketDelete := bucketDeleteCommand(); cmdBucketDelete != nil {
		cmds.Add(cmdBucketDelete)
	}
	if cmdBucketGet := bucketGetCommand(); cmdBucketGet != nil {
		cmds.Add(cmdBucketGet)
	}
	if cmdBucketList := bucketListCommand(); cmdBucketList != nil {
		cmds.Add(cmdBucketList)
	}
	if cmdBucketUpdate := bucketUpdateCommand(); cmdBucketUpdate != nil {
		cmds.Add(cmdBucketUpdate)
	}
	if cmdConfigGet := configGetCommand(); cmdConfigGet != nil {
		cmds.Add(cmdConfigGet)
	}
	if cmdConfigInstall := configInstallCommand(); cmdConfigInstall != nil {
		cmds.Add(cmdConfigInstall)
	}

	return cmds
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
