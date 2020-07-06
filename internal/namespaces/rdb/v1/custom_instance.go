package rdb

import (
	"context"
	"reflect"
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
