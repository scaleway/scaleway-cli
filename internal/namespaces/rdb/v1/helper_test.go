package rdb_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
	rdbSDK "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	name     = "cli-test"
	user     = "foobar"
	password = "{4xdl*#QOoP+&3XRkGA)]"
	engine   = "PostgreSQL-15"
)

func fetchLatestEngine(engine string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		api := rdbSDK.NewAPI(ctx.Client)
		dbEngine, err := api.FetchLatestEngineVersion(engine)
		if err != nil {
			return err
		}
		ctx.Meta["latestEngine"] = dbEngine.Name

		return nil
	}
}

func createInstance(engine string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"Instance",
		fmt.Sprintf(baseCommand, name, engine, user, password),
	)
}

func createInstanceWithPrivateNetwork() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"Instance",
		fmt.Sprintf(baseCommand+privateNetworkStaticSpec, name, engine, user, password),
	)
}

func createInstanceWithPrivateNetworkAndLoadBalancer() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"Instance",
		fmt.Sprintf(
			baseCommand+privateNetworkStaticSpec+loadBalancerSpec,
			name,
			engine,
			user,
			password,
		),
	)
}

func createPN() core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		api := vpc.NewAPI(ctx.Client)
		pn, err := api.CreatePrivateNetwork(&vpc.CreatePrivateNetworkRequest{})
		if err != nil {
			return err
		}
		ctx.Meta["PN"] = pn
		if len(pn.Subnets) > 0 {
			ctx.Meta["IPNet"], err = getIPSubnet(pn.Subnets[0].Subnet)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func getIPSubnet(ipNet scw.IPNet) (*string, error) {
	addr := ipNet.IP.To4()
	if addr == nil {
		return nil, errors.New("could get ip 4 bytes")
	}
	addr = addr.Mask(addr.DefaultMask())
	addr[3] = +3

	sz, _ := ipNet.Mask.Size()
	ipNetStr := fmt.Sprintf("%s/%d", addr.String(), sz)

	return &ipNetStr, nil
}

func listEndpointsInMeta() core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		instance := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance
		for _, endpoint := range instance.Endpoints {
			if endpoint.PrivateNetwork != nil {
				ctx.Meta["PrivateEndpoint"] = endpoint
			} else if endpoint.LoadBalancer != nil || endpoint.DirectAccess != nil {
				ctx.Meta["PublicEndpoint"] = endpoint
			}
		}

		return nil
	}
}

func deleteInstance() core.AfterFunc {
	return core.ExecAfterCmd("scw rdb instance delete {{ .Instance.ID }}")
}

func deleteInstanceAndWait() core.AfterFunc {
	return core.ExecAfterCmd("scw rdb instance delete {{ .Instance.ID }} --wait")
}

func createInstanceDirect(_ string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		result := ctx.ExecuteCmd([]string{
			"scw", "rdb", "instance", "create",
			"node-type=DB-DEV-S",
			"is-ha-cluster=false",
			"name=" + name,
			"engine=" + engine,
			"user-name=" + user,
			"password=" + password,
			"--wait",
		})
		ctx.Meta["Instance"] = result

		return nil
	}
}

func createBackupDirect(metaKey string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		instanceResult := ctx.Meta["Instance"].(rdb.CreateInstanceResult)
		instance := instanceResult.Instance

		result := ctx.ExecuteCmd([]string{
			"scw", "rdb", "backup", "create",
			"name=cli-test-backup",
			"expires-at=2032-01-02T15:04:05-07:00",
			"instance-id=" + instance.ID,
			"database-name=rdb",
			"--wait",
		})
		ctx.Meta[metaKey] = result

		return nil
	}
}

func deleteBackupDirect(metaKey string) core.AfterFunc {
	return func(ctx *core.AfterFuncCtx) error {
		backup := ctx.Meta[metaKey].(*rdbSDK.DatabaseBackup)
		ctx.ExecuteCmd([]string{
			"scw", "rdb", "backup", "delete",
			backup.ID,
		})

		return nil
	}
}

func deleteInstanceDirect() core.AfterFunc {
	return func(ctx *core.AfterFuncCtx) error {
		instance := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance
		ctx.ExecuteCmd([]string{
			"scw", "rdb", "instance", "delete",
			instance.ID,
		})

		return nil
	}
}

func waitForInstanceReady(
	executeCmd func([]string) any,
	instanceID string,
	maxWait time.Duration,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), maxWait)
	defer cancel()

	backoff := time.Second
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf(
				"timeout waiting for instance %s to be ready for operations",
				instanceID,
			)
		default:
			result := executeCmd([]string{"scw", "rdb", "instance", "get", instanceID})

			// Try direct type assertion first
			if instance, ok := result.(*rdbSDK.Instance); ok {
				if instance.Status == rdbSDK.InstanceStatusReady {
					time.Sleep(5 * time.Second)

					return nil
				}
			} else {
				v := result.(struct {
					*rdbSDK.Instance
					ACLs []*rdbSDK.ACLRule `json:"acls"`
				})
				if v.Instance != nil && v.Instance.Status == rdbSDK.InstanceStatusReady {
					time.Sleep(5 * time.Second)

					return nil
				}
			}

			time.Sleep(backoff)
			if backoff < 10*time.Second {
				backoff *= 2
			}
		}
	}
}
