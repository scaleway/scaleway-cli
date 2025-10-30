package sshconfig

import (
	"fmt"
	"strings"
)

type BastionHost struct {
	Name    string
	Address string
	Port    uint32

	Hosts []SimpleHost
}

func (b BastionHost) Config() string {
	bastionConfig := make([]string, 1, len(b.Hosts)+1)
	bastionConfig[0] = fmt.Sprintf(`Host %s
  ProxyJump bastion@%s
`,
		b.name(),
		b.address())

	for _, host := range b.Hosts {
		host.Name = fmt.Sprintf("%s.%s", host.Name, b.Name)
		bastionConfig = append(bastionConfig, fmt.Sprintf(`Host %s
  User %s
`,
			host.name(),
			host.user()))
	}

	return strings.Join(bastionConfig, "")
}

func (b BastionHost) name() string {
	return "*." + b.Name
}

func (b BastionHost) address() string {
	return fmt.Sprintf("%s:%d", b.Address, b.Port)
}
