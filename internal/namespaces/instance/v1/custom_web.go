package instance

import "github.com/scaleway/scaleway-cli/v2/core"

func addWebUrls(cmds *core.Commands) {
	const imageURL = "https://console.scaleway.com/instance/images"

	cmds.MustFind("instance").WebURL = "https://console.scaleway.com/instance/servers"

	cmds.MustFind("instance", "server").WebURL = "https://console.scaleway.com/instance/servers"
	cmds.MustFind("instance", "server", "get").WebURL = "https://console.scaleway.com/instance/servers/{{ .Zone }}/{{ .ServerID }}/overview"
	cmds.MustFind("instance", "server", "create").WebURL = "https://console.scaleway.com/instance/servers/create?zone={{ .Zone }}&offerName={{ .Type }}"

	cmds.MustFind("instance", "image").WebURL = imageURL

	cmds.MustFind("instance", "ip").WebURL = "https://console.scaleway.com/instance/ips"

	cmds.MustFind("instance", "placement-group").WebURL = "https://console.scaleway.com/instance/placement-groups"
	cmds.MustFind("instance", "placement-group", "get").WebURL = "https://console.scaleway.com/instance/placement-groups/{{ .Zone }}/{{ .PlacementGroupID }}/overview"

	cmds.MustFind("instance", "private-nic", "list").WebURL = "https://console.scaleway.com/instance/servers/{{ .Zone }}/{{ .ServerID }}/private-networks"

	cmds.MustFind("instance", "security-group").WebURL = "https://console.scaleway.com/instance/security-groups"
	cmds.MustFind("instance", "snapshot").WebURL = "https://console.scaleway.com/instance/snapshots"
	cmds.MustFind("instance", "volume").WebURL = "https://console.scaleway.com/instance/volumes"

	cmds.MustFind("instance", "snapshot").WebURL = imageURL
	cmds.MustFind("instance", "snapshot").WebURL = imageURL
}
