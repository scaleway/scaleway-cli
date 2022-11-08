package lb

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
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

func lbMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	type tmp lb.LB
	loadbalancer := tmp(i.(lb.LB))

	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "IP",
		},
		{
			FieldName: "Instances",
		},
	}

	str, err := human.Marshal(loadbalancer, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

func lbBackendMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	type tmp lb.Backend
	backend := tmp(i.(lb.Backend))

	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "HealthCheck",
		},
		{
			FieldName: "Pool",
		},
		{
			FieldName: "LB",
		},
	}

	str, err := human.Marshal(backend, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}
