package object

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
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
	if cmdObjectBucketLifecycle := objectBucketLifecycle(); cmdObjectBucketLifecycle != nil {
		cmds.Add(cmdObjectBucketLifecycle)
	}
	if cmdObjectBucketCors := objectBucketCors(); cmdObjectBucketCors != nil {
		cmds.Add(cmdObjectBucketCors)
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
	if cmdLifecycleCreate := lifecycleCreateCommand(); cmdLifecycleCreate != nil {
		cmds.Add(cmdLifecycleCreate)
	}
	if cmdLifecycleDelete := lifecycleDeleteCommand(); cmdLifecycleDelete != nil {
		cmds.Add(cmdLifecycleDelete)
	}
	if cmdLifecycleGet := lifecycleGetCommand(); cmdLifecycleGet != nil {
		cmds.Add(cmdLifecycleGet)
	}
	if cmdBucketCorsCreate := bucketCorsCreateCommand(); cmdBucketCorsCreate != nil {
		cmds.Add(cmdBucketCorsCreate)
	}
	if cmdBucketCorsDelete := bucketCorsDeleteCommand(); cmdBucketCorsDelete != nil {
		cmds.Add(cmdBucketCorsDelete)
	}
	if cmdBucketCorsGet := bucketCorsGetCommand(); cmdBucketCorsGet != nil {
		cmds.Add(cmdBucketCorsGet)
	}
	if cmdObjectBucketPolicy := objectBucketPolicy(); cmdObjectBucketPolicy != nil {
		cmds.Add(cmdObjectBucketPolicy)
	}
	if cmdBucketPolicyCreate := bucketPolicyCreateCommand(); cmdBucketPolicyCreate != nil {
		cmds.Add(cmdBucketPolicyCreate)
	}
	if cmdBucketPolicyDelete := bucketPolicyDeleteCommand(); cmdBucketPolicyDelete != nil {
		cmds.Add(cmdBucketPolicyDelete)
	}
	if cmdBucketPolicyGet := bucketPolicyGetCommand(); cmdBucketPolicyGet != nil {
		cmds.Add(cmdBucketPolicyGet)
	}

	cmds.Add(bucketCockpitCommand())

	// human.RegisterMarshalerFunc(BucketResponse{}, bucketResponseMarshalerFunc)
	// human.RegisterMarshalerFunc(bucketInfo{}, bucketInfoMarshalerFunc)
	// human.RegisterMarshalerFunc(BucketGetResult{}, bucketGetResultMarshalerFunc)
	// human.RegisterMarshalerFunc(s3.ListBucketsOutput{}.Buckets, bucketMarshalerFunc)
	// TODO
	// func bucketResponseMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	// 	resp := i.(BucketResponse)

	// 	messageStr, err := resp.SuccessResult.MarshalHuman()
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	bucketStr, err := bucketInfoMarshalerFunc(*resp.BucketInfo, opt)
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	return strings.Join([]string{
	// 		messageStr,
	// 		bucketStr,
	// 	}, "\n"), nil
	// }
	//
	// func bucketInfoMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	// 	// To avoid recursion of human.Marshal we create a dummy type
	// 	type tmp bucketInfo
	// 	info := tmp(i.(bucketInfo))

	// 	opt.Sections = []*human.MarshalSection{
	// 		{
	// 			FieldName:   "Tags",
	// 			HideIfEmpty: true,
	// 		},
	// 		{
	// 			FieldName:   "ACL",
	// 			HideIfEmpty: true,
	// 		},
	// 	}
	// 	str, err := human.Marshal(info, opt)
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	return str, nil
	// }

	// func bucketGetResultMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	// 	type tmp BucketGetResult
	// 	result := tmp(i.(BucketGetResult))
	// 	opt.Sections = []*human.MarshalSection{
	// 		{
	// 			FieldName:   "Tags",
	// 			HideIfEmpty: true,
	// 		},
	// 		{
	// 			FieldName:   "ACL",
	// 			HideIfEmpty: true,
	// 		},
	// 	}
	// 	str, err := human.Marshal(result, opt)
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	return str, nil
	// }

	// func bucketMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	// 	type tmp []types.Bucket
	// 	result := tmp(i.([]types.Bucket))
	// 	opt.Fields = []*human.MarshalFieldOpt{
	// 		{
	// 			FieldName: "Name",
	// 			Label:     "Name",
	// 		},
	// 		{
	// 			FieldName: "CreationDate",
	// 			Label:     "Creation Date",
	// 		},
	// 	}
	// 	str, err := human.Marshal(result, opt)
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	return str, nil
	// }

	return cmds
}
