package tui

import (
	"context"
	"os/exec"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Connector lets the TUI connect to the selected row of a resource by
// running an interactive client (ssh, psql, mongosh, k9s, ...) on the
// terminal. The TUI suspends itself (tea.ExecProcess) while the client
// runs, then resumes and re-fetches the list.
//
// Connectors are registered per TUI resource name (e.g. "instance server")
// by the namespace packages, reusing the CLI command implementations.
// Adding a new connector is one RegisterConnector call plus a Connect
// function returning the *exec.Cmd to run.
type Connector struct {
	// Hint is the label shown in the header hints ("[ c ] ssh").
	Hint string
	// Connect builds the interactive command for row. It may call the
	// API: the app runs it in a tea.Cmd, outside the render loop.
	// cleanup (may be nil) runs after the command exits.
	Connect func(ctx context.Context, client *scw.Client, loc *Locality, row Row) (*exec.Cmd, func(), error)
}

// connectorRegistry holds the connectors, keyed by TUI resource name.
var connectorRegistry = map[string]Connector{}

// RegisterConnector makes the selected row of resource connectable.
func RegisterConnector(resource string, c Connector) {
	connectorRegistry[resource] = c
}

// ConnectorFor returns the connector of resource, if any.
func ConnectorFor(resource string) (Connector, bool) {
	c, ok := connectorRegistry[resource]

	return c, ok
}

// connectRequestMsg asks the app to connect to a row (the app owns the
// client and the session locality).
type connectRequestMsg struct {
	resource string
	row      Row
}

// connectReadyMsg carries the built connect command (or its error).
type connectReadyMsg struct {
	cmd     *exec.Cmd
	cleanup func()
	err     error
}

// connectDoneMsg reports the exit of the connect command.
type connectDoneMsg struct {
	cleanup func()
	err     error
}
