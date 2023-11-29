package sshconfig

import "fmt"

type SimpleHost struct {
	Name    string
	Address string
	User    string
}

func (h SimpleHost) Config() string {
	return fmt.Sprintf(`Host %s
  Hostname %s
  User %s
`,
		h.name(),
		h.address(),
		h.user())
}

func (h SimpleHost) name() string {
	return h.Name
}

func (h SimpleHost) address() string {
	return h.Address
}

func (h SimpleHost) user() string {
	if h.User == "" {
		return "root"
	}

	return h.User
}
