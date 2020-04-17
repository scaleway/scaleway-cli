package object

import (
	"github.com/scaleway/scaleway-cli/internal/core"
)

var s3cmd_template = `
[default]
access_key = {{ .AccessKey }}
bucket_location = {{ .Region }}
host_base = s3.{{ .Region }}.scw.cloud
host_bucket = %(bucket)s.s3.fr-par.scw.cloud
secret_key = {{ .SecretKey }}
use_https = True
`

var awsCli = ``

var rclone = ``

var mc = ``

func GetCommands() *core.Commands {
	return core.NewCommands(
		objectRoot(),
		getCommand(),
		installCommand(),
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
		Short: `An image is a backup of an instance`,
		Long: `Images are backups of your instances.
You can reuse that image to restore your data or create a series of instances with a predefined configuration.

An image is a complete backup of your server including all volumes.
`,
		Namespace: "object",
	}
}

func getCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "config",
		Verb:      "get",
		Short:     "",
		Long:      "",
		ArgsType:  nil,
		ArgSpecs:  nil,
		Examples:  nil,
		SeeAlsos:  nil,
		Run:       nil,
	}
}

func installCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "config",
		Verb:      "install",
		Short:     "Install a s3 related configuration file",
		Long:      "",
		ArgsType:  nil,
		ArgSpecs:  nil,
		Examples:  nil,
		SeeAlsos:  nil,
		Run:       nil,
	}
}
