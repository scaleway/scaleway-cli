package lb

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func createLB() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"LB",
		"scw lb lb create name=cli-test description=cli-test --wait",
	)
}

func deleteLB() core.AfterFunc {
	return core.ExecAfterCmd("scw lb lb delete {{ .LB.ID }}")
}

func createInstance() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"Instance",
		"scw instance server create stopped=true image=ubuntu_focal",
	)
}

func deleteInstance() core.AfterFunc {
	return core.ExecAfterCmd("scw instance server delete {{ .Instance.ID }}")
}

func createBackend(forwardPort int32) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"Backend",
		fmt.Sprintf("scw lb backend create lb-id={{ .LB.ID }} name=cli-test forward-protocol=tcp forward-port=%d forward-port-algorithm=roundrobin sticky-sessions=none health-check.port=8888 health-check.check-max-retries=5", forwardPort),
	)
}

func addIP2Backend(ip string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"AddIP2Backend",
		fmt.Sprintf("scw lb backend add-servers {{ .Backend.ID }} server-ip.0=%s", ip),
	)
}

func createFrontend(inboundPort int32) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"Frontend",
		fmt.Sprintf("scw lb frontend create lb-id={{ .LB.ID }} backend-id={{ .Backend.ID }} name=cli-test inbound-port=%d", inboundPort),
	)
}
