package object

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
)

func GetCommands() *core.Commands {
	cmds := core.NewCommands()

	human.RegisterMarshalerFunc(BucketResponse{}, bucketResponseMarshalerFunc)
	human.RegisterMarshalerFunc(bucketInfo{}, bucketInfoMarshalerFunc)
	human.RegisterMarshalerFunc(BucketGetResult{}, bucketGetResultMarshalerFunc)
	human.RegisterMarshalerFunc(s3.ListBucketsOutput{}.Buckets, bucketMarshalerFunc)

	if cmdObjectRoot := objectRoot(); cmdObjectRoot != nil {
		cmds.Add(cmdObjectRoot)
	}
	if cmdObjectConfig := objectConfig(); cmdObjectConfig != nil {
		cmds.Add(cmdObjectConfig)
	}
	if cmdObjectBucket := objectBucket(); cmdObjectBucket != nil {
		cmds.Add(cmdObjectBucket)
	}
	if cmdObjectBucketLifecycle := objectBucketLifecycle(); cmdObjectBucketLifecycle != nil {
		cmds.Add(cmdObjectBucketLifecycle)
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
	if cmdLifecycleConfigurationCreate := lifecycleConfigurationCreateCommand(); cmdLifecycleConfigurationCreate != nil {
		cmds.Add(cmdLifecycleConfigurationCreate)
	}
	if cmdLifecycleConfigurationDelete := lifecycleConfigurationDeleteCommand(); cmdLifecycleConfigurationDelete != nil {
		cmds.Add(cmdLifecycleConfigurationDelete)
	}
	if cmdLifecycleConfigurationGet := lifecycleConfigurationGetCommand(); cmdLifecycleConfigurationGet != nil {
		cmds.Add(cmdLifecycleConfigurationGet)
	}

	return cmds
}
