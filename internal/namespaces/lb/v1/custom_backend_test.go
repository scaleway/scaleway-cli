package lb_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/lb/v1"
)

func Test_GetBackend(t *testing.T) {
	cmds := lb.GetCommands()
	cmds.Merge(instance.GetCommands())

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createLB(),
			createInstance(),
			createBackend(80),
			addIP2Backend("{{ .Instance.PublicIP.Address }}"),
		),
		Cmd:   "scw lb backend get {{ .Backend.ID }}",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteLB(),
			deleteInstance(),
			deleteLBFlexibleIP(),
		),
	}))
}

func Test_CreateBackend(t *testing.T) {
	cmds := lb.GetCommands()
	cmds.Merge(instance.GetCommands())

	t.Run("With instance ID", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createLB(),
			createRunningInstance(),
		),
		Cmd:   "scw lb backend create lb-id={{ .LB.ID }} name=cli-test instance-server-id.0={{ .Instance.ID }} forward-protocol=tcp forward-port=80 forward-port-algorithm=roundrobin sticky-sessions=none health-check.port=8888 health-check.check-max-retries=5",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteLB(),
			deleteRunningInstance(),
			deleteLBFlexibleIP(),
		),
	}))

	t.Run("With instance ID public IP", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createLB(),
			createRunningInstance(),
		),
		Cmd:   "scw lb backend create lb-id={{ .LB.ID }} name=cli-test instance-server-id.0={{ .Instance.ID }} use-instance-server-public-ip=true forward-protocol=tcp forward-port=80 forward-port-algorithm=roundrobin sticky-sessions=none health-check.port=8888 health-check.check-max-retries=5",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteLB(),
			deleteRunningInstance(),
			deleteLBFlexibleIP(),
		),
	}))

	t.Run("With instance tag", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createLB(),
			createRunningInstanceWithTag(),
		),
		Cmd:   "scw lb backend create lb-id={{ .LB.ID }} name=cli-test instance-server-tag.0={{index .Instance.Tags 0}} forward-protocol=tcp forward-port=80 forward-port-algorithm=roundrobin sticky-sessions=none health-check.port=8888 health-check.check-max-retries=5",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteLB(),
			deleteRunningInstance(),
			deleteLBFlexibleIP(),
		),
	}))

	t.Run("With instance tag public IP", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createLB(),
			createRunningInstanceWithTag(),
		),
		Cmd:   "scw lb backend create lb-id={{ .LB.ID }} name=cli-test instance-server-tag.0={{index .Instance.Tags 0}} use-instance-server-public-ip=true forward-protocol=tcp forward-port=80 forward-port-algorithm=roundrobin sticky-sessions=none health-check.port=8888 health-check.check-max-retries=5",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteLB(),
			deleteRunningInstance(),
			deleteLBFlexibleIP(),
		),
	}))
}

func Test_AddBackendServers(t *testing.T) {
	cmds := lb.GetCommands()
	cmds.Merge(instance.GetCommands())

	t.Run("With instance ID", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createLB(),
			createRunningInstance(),
			createBackend(80),
		),
		Cmd:   "scw lb backend add-servers {{ .Backend.ID }} instance-server-id.0={{ .Instance.ID }}",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteLB(),
			deleteRunningInstance(),
			deleteLBFlexibleIP(),
		),
	}))

	t.Run("With instance ID public IP", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createLB(),
			createRunningInstance(),
			createBackend(80),
		),
		Cmd:   "scw lb backend add-servers {{ .Backend.ID }} instance-server-id.0={{ .Instance.ID }} use-instance-server-public-ip=true",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteLB(),
			deleteRunningInstance(),
			deleteLBFlexibleIP(),
		),
	}))

	t.Run("With instance tag", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createLB(),
			createRunningInstanceWithTag(),
			createBackend(80),
		),
		Cmd:   "scw lb backend add-servers {{ .Backend.ID }} instance-server-tag.0={{index .Instance.Tags 0}}",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteLB(),
			deleteRunningInstance(),
			deleteLBFlexibleIP(),
		),
	}))

	t.Run("With instance tag public IP", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createLB(),
			createRunningInstanceWithTag(),
			createBackend(80),
		),
		Cmd:   "scw lb backend add-servers {{ .Backend.ID }} instance-server-tag.0={{index .Instance.Tags 0}} use-instance-server-public-ip=true",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteLB(),
			deleteRunningInstance(),
			deleteLBFlexibleIP(),
		),
	}))
}
