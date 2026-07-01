package registry

import "github.com/scaleway/scaleway-sdk-go/scw"

type (
	program  string
	programs []program
)

const (
	docker         = program("docker")
	podman         = program("podman")
	endpointPrefix = "rg."
	cloudDomain    = ".scw.cloud"
	euDomain       = ".scw.eu"
)

var regionRegistryDomains = map[scw.Region]string{
	scw.RegionItMil: euDomain,
}

func getRegistryEndpoint(region scw.Region) string {
	domain := cloudDomain
	if regionDomain, ok := regionRegistryDomains[region]; ok {
		domain = regionDomain
	}

	return endpointPrefix + region.String() + domain
}

var availablePrograms = programs{docker, podman}

func (p programs) StringArray() []string {
	res := make([]string, 0, len(p))
	for _, prog := range p {
		res = append(res, string(prog))
	}

	return res
}
