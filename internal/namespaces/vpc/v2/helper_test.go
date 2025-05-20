package vpc_test

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
	rdbAPI "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

func createInstance() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"Instance",
		"scw instance server create type=DEV1-S stopped=true image=ubuntu_focal",
	)
}

func deleteInstance() core.AfterFunc {
	return core.ExecAfterCmd("scw instance server delete {{ .Instance.ID }} --wait")
}

func createPN() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"PN",
		"scw vpc private-network create",
	)
}

func deletePN() core.AfterFunc {
	return core.ExecAfterCmd("scw vpc private-network delete {{ .PN.ID }}")
}

func createNIC() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"NIC",
		"scw instance private-nic create server-id={{ .Instance.ID }} private-network-id={{ .PN.ID }}",
	)
}

func createLB() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"LB",
		"scw lb lb create name=cli-test description=cli-test --wait",
	)
}

func attachLB() core.BeforeFunc {
	return core.ExecBeforeCmd(
		"scw lb private-network attach {{ .LB.ID }} private-network-id={{ .PN.ID }}",
	)
}

func detachLB() core.AfterFunc {
	return core.ExecAfterCmd(
		"scw lb private-network detach {{ .LB.ID }} private-network-id={{ .PN.ID }}",
	)
}

func deleteLB() core.AfterFunc {
	return core.ExecAfterCmd("scw lb lb delete {{ .LB.ID }}")
}

func createRdbInstance(metaKey, engineName string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		api := rdbAPI.NewAPI(ctx.Client)
		engine, err := api.FetchLatestEngineVersion(engineName)
		if err != nil {
			return err
		}
		cmd := fmt.Sprintf(
			"scw rdb instance create node-type=DB-DEV-S is-ha-cluster=false name=cli-test engine=%s user-name=foobar password={4xdl*#QOoP+&3XRkGA)] init-endpoints.0.private-network.private-network-id={{ .PN.ID }} init-endpoints.0.private-network.service-ip=192.168.0.1/24 --wait",
			engine.Name,
		)

		return core.ExecStoreBeforeCmd(metaKey, cmd)(ctx)
	}
}

func createMongoDBInstance(metaKey string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cmd := "scw mongodb instance create node-type=MGDB-PLAY2-NANO name=mongo-cli-test user-name=foobar password={4xdl*#QOoP+&3XRkGA)] endpoints.0.private-network.private-network-id={{ .PN.ID }} --wait"

		return core.ExecStoreBeforeCmd(metaKey, cmd)(ctx)
	}
}

func detachRdbInstance() core.AfterFunc {
	return core.ExecAfterCmd(
		"scw rdb endpoint delete {{ (index .RDB.Endpoints 0).ID  }} instance-id={{ .RDB.ID }}",
	)
}

func waitRdbInstance() core.AfterFunc {
	return core.ExecAfterCmd(
		"scw rdb instance wait {{ .RDB.ID }}",
	)
}

func deleteRdbInstance() core.AfterFunc {
	return core.ExecAfterCmd("scw rdb instance delete {{ .RDB.ID }}")
}

func detachMongoDBInstance() core.AfterFunc {
	return core.ExecAfterCmd(
		"scw mongodb endpoint delete {{ (index .mongoDB.Endpoints 0).ID  }}",
	)
}

func deleteMongoDBInstance() core.AfterFunc {
	return core.ExecAfterCmd("scw mongodb instance delete {{ .mongoDB.ID }}")
}

func waitMongoDBInstance() core.AfterFunc {
	return core.ExecAfterCmd(
		"scw mongodb instance wait {{ .mongoDB.ID }}",
	)
}
