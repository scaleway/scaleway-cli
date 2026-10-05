package docs

var commandTips = map[string][]string{
	"instance.server.create": {
		"Use --boot-type=local to boot from a local image for faster startup.",
		"Use --root-volume=size:50GB to specify the root disk size.",
		"Add --wait to wait for the server to be ready before returning.",
	},
	"instance.server.list": {
		"Use -o json to pipe output to jq for filtering.",
		"Use --organization-id to filter by organization.",
	},
	"k8s.cluster.create": {
		"Use --version to specify the Kubernetes version.",
		"Use --cni to specify the container network interface (cilium or calico).",
	},
}

var commandPitfalls = map[string][]string{
	"instance.server.create": {
		"Forgetting to specify --image or --boot-type will result in a server that cannot boot.",
		"The --root-volume flag must include a size suffix (e.g. 50GB).",
	},
	"instance.server.list": {
		"By default, only servers in the default zone are listed. Use --zone=all to list across all zones.",
	},
}

func GetTips(commandPath string) []string {
	return commandTips[commandPath]
}

func GetPitfalls(commandPath string) []string {
	return commandPitfalls[commandPath]
}

func GetRelatedTutorials(commandPath string) []string {
	return nil
}
