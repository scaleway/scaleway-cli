package lb

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

var (
	backendServerStatsHealthCheckStatusMarshalSpecs = human.EnumMarshalSpecs{
		lb.BackendServerStatsHealthCheckStatusPassed:   &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "passed"},
		lb.BackendServerStatsHealthCheckStatusFailed:   &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "failed"},
		lb.BackendServerStatsHealthCheckStatusUnknown:  &human.EnumMarshalSpec{Attribute: color.Faint, Value: "unknown"},
		lb.BackendServerStatsHealthCheckStatusNeutral:  &human.EnumMarshalSpec{Attribute: color.Faint, Value: "neutral"},
		lb.BackendServerStatsHealthCheckStatusCondpass: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "condition passed"},
	}
	backendServerStatsServerStateMarshalSpecs = human.EnumMarshalSpecs{
		lb.BackendServerStatsServerStateStopped:  &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "stopped"},
		lb.BackendServerStatsServerStateStarting: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "starting"},
		lb.BackendServerStatsServerStateRunning:  &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "running"},
		lb.BackendServerStatsServerStateStopping: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "stopping"},
	}
)
