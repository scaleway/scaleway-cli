package registry

type program string
type programs []program

const (
	docker         = program("docker")
	podman         = program("podman")
	endpointPrefix = "rg."
	endpointSuffix = ".scw.cloud"
)

var (
	availablePrograms = programs{docker, podman}
)

func (p programs) StringArray() []string {
	var res []string
	for _, prog := range p {
		res = append(res, string(prog))
	}
	return res
}
