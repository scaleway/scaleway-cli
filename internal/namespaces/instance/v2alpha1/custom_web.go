package instance

import "github.com/scaleway/scaleway-cli/v2/core"

func addWebUrls(cmds *core.Commands) {
	cmds.MustFind("instance", "placement-group").WebURL = "https://console.scaleway.com/instance/placement-groups"
	cmds.MustFind("instance", "placement-group", "get").WebURL = "https://console.scaleway.com/instance/{{ .Zone }}/placement-groups/{{ .PlacementGroupID }}/overview"

	cmds.MustFind("instance", "private-network-interface", "list").WebURL = "https://console.scaleway.com/instance/{{ .Zone }}/servers/{{ index .ServerIDs 0 }}/network"
}
