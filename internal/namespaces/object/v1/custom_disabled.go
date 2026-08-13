//go:build !(darwin || linux || windows)

package object

import "github.com/scaleway/scaleway-cli/v2/core"

func objectRoot() *core.Command {
	return nil
}

func objectConfig() *core.Command {
	return nil
}

func objectBucket() *core.Command {
	return nil
}

func bucketCreateCommand() *core.Command {
	return nil
}

func bucketDeleteCommand() *core.Command {
	return nil
}

func bucketListCommand() *core.Command {
	return nil
}

func bucketGetCommand() *core.Command {
	return nil
}

func bucketUpdateCommand() *core.Command {
	return nil
}

func configGetCommand() *core.Command {
	return nil
}

func configInstallCommand() *core.Command {
	return nil
}

func objectBucketPolicy() *core.Command {
	return nil
}

func bucketPolicyCreateCommand() *core.Command {
	return nil
}

func bucketPolicyDeleteCommand() *core.Command {
	return nil
}

func bucketPolicyGetCommand() *core.Command {
	return nil
}
