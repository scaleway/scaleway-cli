package alias

type Config struct {
	// Aliases are raw aliases that allow to expand a command
	// "scw instance sl", sl may be an alias and would expand command
	// "scw instance server list"
	Aliases map[string][]string `yaml:"aliases"`
	// ResourceAliases are aliases specific to a resource
	// it allows to add an alternative name to a namespace, resource or verb
	ResourceAliases map[string][]string `yaml:"resources"`
}

func EmptyConfig() *Config {
	return &Config{
		Aliases:         map[string][]string{},
		ResourceAliases: map[string][]string{},
	}
}

// GetAlias return raw alias for a given string
func (c *Config) GetAlias(name string) []string {
	alias, aliasExists := c.Aliases[name]
	if aliasExists {
		return alias
	}
	return nil
}

// ResolveAliases resolve raw aliases in given command
// "scw isl" may return "scw instance server list"
func (c *Config) ResolveAliases(command []string) []string {
	expandedCommand := make([]string, 0, len(command))
	for _, arg := range command {
		if alias := c.GetAlias(arg); alias != nil {
			expandedCommand = append(expandedCommand, alias...)
		} else {
			expandedCommand = append(expandedCommand, arg)
		}
	}
	return expandedCommand
}

// AddAlias add alias to config
// return true if alias has been replaced
func (c *Config) AddAlias(name string, command []string) bool {
	_, exists := c.Aliases[name]
	c.Aliases[name] = command
	return exists
}

// DeleteAlias deletes an alias
// return true if alias was deleted
func (c *Config) DeleteAlias(name string) bool {
	_, exists := c.Aliases[name]
	delete(c.Aliases, name)
	return exists
}
