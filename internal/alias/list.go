package alias

import "strings"

type Alias struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

func (c *Config) List() []Alias {
	list := []Alias(nil)
	for name, alias := range c.Aliases {
		list = append(list, Alias{
			Name:    name,
			Command: strings.Join(alias, "."),
		})
	}
	return list
}

type ResourceAlias struct {
	Path    []string
	Aliases []string
}
