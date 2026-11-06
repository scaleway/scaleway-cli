package instance

import "github.com/scaleway/scaleway-cli/v2/core"

func addWebUrls(cmds *core.Commands) {
	cmds.MustFind("instance", "placement-group").WebURL = "https://console.scaleway.com/instance/placement-groups"
	cmds.MustFind("instance", "placement-group", "get").WebURL = "https://console.scaleway.com/instance/placement-groups/{{ .Zone }}/{{ .PlacementGroupID }}/overview"
}
