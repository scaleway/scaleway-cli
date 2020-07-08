package rdb

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	instanceActionTimeout = 20 * time.Minute
)

type serverWaitRequest struct {
	InstanceID string
	Region     scw.Region
}

func instanceMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	// To avoid recursion of human.Marshal we create a dummy type
	type tmp rdb.Instance
	instance := tmp(i.(rdb.Instance))

	// Sections
	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "Endpoint",
		},
		{
			FieldName: "Volume",
		},
		{
			FieldName: "BackupSchedule",
		},
		{
			FieldName: "Settings",
		},
	}

	str, err := human.Marshal(instance, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

func backupScheduleMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	// To avoid recursion of human.Marshal we create a dummy type
	type LocalBackupSchedule rdb.BackupSchedule
	type tmp struct {
		LocalBackupSchedule
		Frequency *scw.Duration `json:"frequency"`
		Retention *scw.Duration `json:"retention"`
	}

	backupSchedule := tmp{
		LocalBackupSchedule: LocalBackupSchedule(i.(rdb.BackupSchedule)),
		Frequency: &scw.Duration{
			Seconds: int64(i.(rdb.BackupSchedule).Frequency) * 24 * 3600,
		},
		Retention: &scw.Duration{
			Seconds: int64(i.(rdb.BackupSchedule).Retention) * 24 * 3600,
		},
	}

	str, err := human.Marshal(backupSchedule, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

func instanceCloneBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := rdb.NewAPI(core.ExtractClient(ctx))
		return api.WaitForInstance(&rdb.WaitForInstanceRequest{
			InstanceID:    respI.(*rdb.Instance).ID,
			Region:        respI.(*rdb.Instance).Region,
			Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

func instanceCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("node-type").Default = core.DefaultValueSetter("DB-DEV-S")
	c.ArgSpecs.GetByName("node-type").EnumValues = nodeTypes

	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := rdb.NewAPI(core.ExtractClient(ctx))
		return api.WaitForInstance(&rdb.WaitForInstanceRequest{
			InstanceID:    respI.(*rdb.Instance).ID,
			Region:        respI.(*rdb.Instance).Region,
			Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	// Waiting for API to accept uppercase node-type
	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		args := argsI.(*rdb.CreateInstanceRequest)
		args.NodeType = strings.ToLower(args.NodeType)
		return runner(ctx, args)
	}

	return c
}

func instanceUpgradeBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("node-type").EnumValues = nodeTypes

	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := rdb.NewAPI(core.ExtractClient(ctx))
		return api.WaitForInstance(&rdb.WaitForInstanceRequest{
			InstanceID:    respI.(*rdb.Instance).ID,
			Region:        respI.(*rdb.Instance).Region,
			Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

func instanceWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for an instance to reach a stable state`,
		Long:      `Wait for an instance to reach a stable state. This is similar to using --wait flag.`,
		Namespace: "rdb",
		Resource:  "instance",
		Verb:      "wait",
		ArgsType:  reflect.TypeOf(serverWaitRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			api := rdb.NewAPI(core.ExtractClient(ctx))
			return api.WaitForInstance(&rdb.WaitForInstanceRequest{
				Region:        argsI.(*serverWaitRequest).Region,
				InstanceID:    argsI.(*serverWaitRequest).InstanceID,
				Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
				RetryInterval: core.DefaultRetryInterval,
			})
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "instance-id",
				Short:      `ID of the instance you want to wait for.`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for an instance to reach a stable state",
				ArgsJSON: `{"instance_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}
