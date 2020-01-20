package instance

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

//
// Marshalers
//

// serverStateMarshalerFunc marshals a instance.ServerState.
var (
	serverStateAttributes = human.Attributes{
		instance.ServerStateRunning:        color.FgGreen,
		instance.ServerStateStopped:        color.Faint,
		instance.ServerStateStoppedInPlace: color.Faint,
		instance.ServerStateStarting:       color.FgBlue,
		instance.ServerStateStopping:       color.FgBlue,
		instance.ServerStateLocked:         color.FgRed,
	}
)
