package instance

import "github.com/scaleway/scaleway-cli/v2/core"

func addWebUrls(cmds *core.Commands) {
	cmds.MustFind("instance", "private-network-interface", "list").WebURL = "https://console.scaleway.com/instance/servers/{{ .Zone }}/{{ .ServerID }}/network"
}
