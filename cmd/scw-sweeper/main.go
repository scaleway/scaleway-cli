package main

import (
	"log"
	"os"

	accountSweeper "github.com/scaleway/scaleway-sdk-go/api/account/v3/sweepers"
	applesiliconSweeper "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1/sweepers"
	baremetalSweeper "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1/sweepers"
	blockSweeper "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1/sweepers"
	cockpitSweeper "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1/sweepers"
	containerSweeper "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1/sweepers"
	flexibleipSweeper "github.com/scaleway/scaleway-sdk-go/api/flexibleip/v1alpha1/sweepers"
	functionSweeper "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1/sweepers"
	iamSweeper "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1/sweepers"
	inferenceSweeper "github.com/scaleway/scaleway-sdk-go/api/inference/v1beta1/sweepers"
	instanceSweeper "github.com/scaleway/scaleway-sdk-go/api/instance/v1/sweepers"
	iotSweeper "github.com/scaleway/scaleway-sdk-go/api/iot/v1/sweepers"
	ipamSweeper "github.com/scaleway/scaleway-sdk-go/api/ipam/v1/sweepers"
	jobsSweeper "github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha1/sweepers"
	k8sSweeper "github.com/scaleway/scaleway-sdk-go/api/k8s/v1/sweepers"
	lbSweeper "github.com/scaleway/scaleway-sdk-go/api/lb/v1/sweepers"
	mnqSweeper "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1/sweepers"
	mongodbSweeper "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1alpha1/sweepers"
	rdbSweeper "github.com/scaleway/scaleway-sdk-go/api/rdb/v1/sweepers"
	redisSweeper "github.com/scaleway/scaleway-sdk-go/api/redis/v1/sweepers"
	registrySweeper "github.com/scaleway/scaleway-sdk-go/api/registry/v1/sweepers"
	secretSweeper "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1/sweepers"
	sdbSweeper "github.com/scaleway/scaleway-sdk-go/api/serverless_sqldb/v1alpha1/sweepers"
	vpcSweeper "github.com/scaleway/scaleway-sdk-go/api/vpc/v2/sweepers"
	vpcgwSweeper "github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1/sweepers"
	webhostingSweeper "github.com/scaleway/scaleway-sdk-go/api/webhosting/v1/sweepers"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func main() {
	exitCode := mainNoExit()
	os.Exit(exitCode)
}

func getConfigProfile() *scw.Profile {
	config, err := scw.LoadConfig()
	if err != nil {
		return &scw.Profile{}
	}
	profile, err := config.GetActiveProfile()
	if err != nil {
		return &scw.Profile{}
	}

	return profile
}

func mainNoExit() int {
	configProfile := getConfigProfile()
	envProfile := scw.LoadEnvProfile()
	profile := scw.MergeProfiles(configProfile, envProfile)

	client, err := scw.NewClient(
		scw.WithProfile(profile),
		scw.WithUserAgent("scw-sweeper"),
		scw.WithEnv(),
	)
	if err != nil {
		log.Fatalf("Cannot create Scaleway client: %s", err)
	}

	err = accountSweeper.SweepAll(client)
	if err != nil {
		log.Fatalf("Error sweeping account: %s", err)

		return -1
	}

	err = applesiliconSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping applesilicon: %s", err)

		return -1
	}

	err = baremetalSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping baremetal: %s", err)

		return -1
	}

	err = cockpitSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping cockpit: %s", err)

		return -1
	}

	err = containerSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping container: %s", err)

		return -1
	}

	err = flexibleipSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping flexibleip: %s", err)

		return -1
	}

	err = functionSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping function: %s", err)

		return -1
	}

	err = iamSweeper.SweepSSHKey(client)
	if err != nil {
		log.Fatalf("Error sweeping iam: %s", err)

		return -1
	}

	err = inferenceSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping inference: %s", err)

		return -1
	}

	err = instanceSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping instance: %s", err)

		return -1
	}

	// Instance servers need to be swept before volumes and snapshots can be swept
	// because volumes and snapshots are attached to servers.
	err = blockSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping block: %s", err)

		return -1
	}

	err = iotSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping iot: %s", err)

		return -1
	}

	err = jobsSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping jobs: %s", err)

		return -1
	}

	err = k8sSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping k8s: %s", err)

		return -1
	}

	err = lbSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping lb: %s", err)

		return -1
	}

	err = mongodbSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping mongodb: %s", err)

		return -1
	}

	err = mnqSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping mnq: %s", err)

		return -1
	}

	err = rdbSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping rdb: %s", err)

		return -1
	}

	err = redisSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping redis: %s", err)

		return -1
	}

	err = registrySweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping registry: %s", err)

		return -1
	}

	err = secretSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping secret: %s", err)

		return -1
	}

	err = sdbSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping sdb: %s", err)

		return -1
	}

	err = vpcSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping vpc: %s", err)

		return -1
	}

	err = vpcgwSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping vpcgw: %s", err)

		return -1
	}

	err = webhostingSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping webhosting: %s", err)

		return -1
	}

	// IPAM IPs need to be swept in the end because we need to be sure
	// that every resource with an attached ip is destroyed before executing it.
	err = ipamSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Fatalf("Error sweeping ipam: %s", err)

		return -1
	}

	return 0
}
