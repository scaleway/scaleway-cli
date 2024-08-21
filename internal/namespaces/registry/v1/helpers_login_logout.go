package registry

type (
	program  string
	programs []program
)

const (
	docker         = program("docker")
	podman         = program("podman")
	endpointPrefix = "rg."
	endpointSuffix = ".scw.cloud"
)

var availablePrograms = programs{docker, podman}

func (p programs) StringArray() []string {
	res := make([]string, 0, len(p))
	for _, prog := range p {
		res = append(res, string(prog))
	}
	return res
}
