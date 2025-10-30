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
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf(`Host %s
  ProxyJump bastion@%s
`,
		b.name(),
		b.address()))

	for _, host := range b.Hosts {
		host.Name = fmt.Sprintf("%s.%s", host.Name, b.Name)
		builder.WriteString(fmt.Sprintf(`Host %s
  User %s
`,
			host.name(),
			host.user()))
	}

	return builder.String()
}

func (b BastionHost) name() string {
	return "*." + b.Name
}

func (b BastionHost) address() string {
	return fmt.Sprintf("%s:%d", b.Address, b.Port)
}
