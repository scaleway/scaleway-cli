package main

import (
	"fmt"
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
	inferenceSweeper "github.com/scaleway/scaleway-sdk-go/api/inference/v1/sweepers"
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

	var errors []string

	err = accountSweeper.SweepAll(client)
	if err != nil {
		log.Printf("Error sweeping account: %s", err)
		errors = append(errors, fmt.Sprintf("account: %s", err))
	}

	err = applesiliconSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping applesilicon: %s", err)
		errors = append(errors, fmt.Sprintf("applesilicon: %s", err))
	}

	err = baremetalSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping baremetal: %s", err)
		errors = append(errors, fmt.Sprintf("baremetal: %s", err))
	}

	err = cockpitSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping cockpit: %s", err)
		errors = append(errors, fmt.Sprintf("cockpit: %s", err))
	}

	err = containerSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping container: %s", err)
		errors = append(errors, fmt.Sprintf("container: %s", err))
	}

	err = flexibleipSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping flexibleip: %s", err)
		errors = append(errors, fmt.Sprintf("flexibleip: %s", err))
	}

	err = functionSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping function: %s", err)
		errors = append(errors, fmt.Sprintf("function: %s", err))
	}

	err = iamSweeper.SweepSSHKey(client)
	if err != nil {
		log.Printf("Error sweeping iam: %s", err)
		errors = append(errors, fmt.Sprintf("iam: %s", err))
	}

	err = inferenceSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping inference: %s", err)
		errors = append(errors, fmt.Sprintf("inference: %s", err))
	}

	err = instanceSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping instance: %s", err)
		errors = append(errors, fmt.Sprintf("instance: %s", err))
	}

	// Instance servers need to be swept before volumes and snapshots can be swept
	// because volumes and snapshots are attached to servers.
	err = blockSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping block: %s", err)
		errors = append(errors, fmt.Sprintf("block: %s", err))
	}

	err = iotSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping iot: %s", err)
		errors = append(errors, fmt.Sprintf("iot: %s", err))
	}

	err = jobsSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping jobs: %s", err)
		errors = append(errors, fmt.Sprintf("jobs: %s", err))
	}

	err = k8sSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping k8s: %s", err)
		errors = append(errors, fmt.Sprintf("k8s: %s", err))
	}

	err = lbSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping lb: %s", err)
		errors = append(errors, fmt.Sprintf("lb: %s", err))
	}

	err = mongodbSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping mongodb: %s", err)
		errors = append(errors, fmt.Sprintf("mongodb: %s", err))
	}

	err = mnqSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping mnq: %s", err)
		errors = append(errors, fmt.Sprintf("mnq: %s", err))
	}

	err = rdbSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping rdb: %s", err)
		errors = append(errors, fmt.Sprintf("rdb: %s", err))
	}

	err = redisSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping redis: %s", err)
		errors = append(errors, fmt.Sprintf("redis: %s", err))
	}

	err = registrySweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping registry: %s", err)
		errors = append(errors, fmt.Sprintf("registry: %s", err))
	}

	err = secretSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping secret: %s", err)
		errors = append(errors, fmt.Sprintf("secret: %s", err))
	}

	err = sdbSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping sdb: %s", err)
		errors = append(errors, fmt.Sprintf("sdb: %s", err))
	}

	err = vpcSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping vpc: %s", err)
		errors = append(errors, fmt.Sprintf("vpc: %s", err))
	}

	err = vpcgwSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping vpcgw: %s", err)
		errors = append(errors, fmt.Sprintf("vpcgw: %s", err))
	}

	err = webhostingSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping webhosting: %s", err)
		errors = append(errors, fmt.Sprintf("webhosting: %s", err))
	}

	// IPAM IPs need to be swept in the end because we need to be sure
	// that every resource with an attached ip is destroyed before executing it.
	err = ipamSweeper.SweepAllLocalities(client)
	if err != nil {
		log.Printf("Error sweeping ipam: %s", err)
		errors = append(errors, fmt.Sprintf("ipam: %s", err))
	}

	// If there were any errors, log them all and exit with fatal
	if len(errors) > 0 {
		log.Printf("Sweeper completed with %d errors:", len(errors))
		for _, errMsg := range errors {
			log.Printf("  - %s", errMsg)
		}
		log.Fatalf("Sweeper failed with %d errors", len(errors))
	}

	return 0
}
