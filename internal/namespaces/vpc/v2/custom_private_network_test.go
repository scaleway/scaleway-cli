package vpc_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/lb/v1"
	mongodb "github.com/scaleway/scaleway-cli/v2/internal/namespaces/mongodb/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/redis/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpc/v2"
)

func Test_GetPrivateNetwork(t *testing.T) {
	cmds := vpc.GetCommands()
	cmds.Merge(instance.GetCommands())
	cmds.Merge(lb.GetCommands())
	cmds.Merge(rdb.GetCommands())
	cmds.Merge(redis.GetCommands())
	cmds.Merge(mongodb.GetCommands())

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createPN(),
			createInstance(),
			createNIC(),
		),
		Cmd:   "scw vpc private-network get {{ .PN.ID }}",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteInstance(),
			deletePN(),
		),
	}))

	t.Run("Multiple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createPN(),
			createInstance(),
			createNIC(),
			createLB(),
			attachLB(),
			createRdbInstance("RDB", "PostgreSQL"),
			createMongoDBInstance("mongoDB"),
		),
		Cmd:   "scw vpc private-network get {{ .PN.ID }}",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			detachLB(),
			deleteLB(),
			deleteInstance(),
			detachRdbInstance(),
			waitRdbInstance(),
			deleteRdbInstance(),
			detachMongoDBInstance(),
			waitMongoDBInstance(),
			deleteMongoDBInstance(),
			deletePN(),
		),
	}))
}
