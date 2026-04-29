package server

import "sort"

// resourceInfo represents information about an MCP resource
type resourceInfo struct {
	Namespace string `json:"namespace"`
	Resource  string `json:"resource"`
	URI       string `json:"uri"`
	Short     string `json:"short"`
}

// ListResources returns a list of available MCP resources based on the server's configuration.
// Resources are read-only endpoints for list commands that can be accessed via URI.
func (s *MCPServer) ListResources() []resourceInfo {
	resources := make([]resourceInfo, 0, len(s.commands))

	for _, cmd := range s.resources {
		resources = append(resources, resourceInfo{
			Namespace: cmd.Command.Namespace,
			Resource:  cmd.Command.Resource,
			URI:       BuildResourceURI(cmd.Command.Namespace, cmd.Command.Resource),
			Short:     cmd.Command.Short,
		})
	}

	// Sort resources by namespace, resource for consistent output
	sort.Slice(resources, func(i, j int) bool {
		if resources[i].Namespace != resources[j].Namespace {
			return resources[i].Namespace < resources[j].Namespace
		}

		return resources[i].Resource < resources[j].Resource
	})

	return resources
}
