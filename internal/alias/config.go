package alias

type Config struct {
	// Aliases are raw aliases that allow to expand a command
	// "scw instance sl", sl may be an alias and would expand command
	// "scw instance server list"
	// key = sl
	// value = server, list
	Aliases map[string][]string `yaml:"aliases"`

	// map of alias using their first word as key
	// value can contain multiple aliases with the same first word
	// key = instance
	// value = isl, isc
	aliasesByFirstWord map[string][]Alias
}

func EmptyConfig() *Config {
	return &Config{
		Aliases: map[string][]string{},
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

// ResolveAliases resolve aliases in given command
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

// ResolveAliasesByFirstWord return list of aliases that start with given first word
// firstWord: instance
// may return
// isl => instance server list
// isc => instance server create
func (c *Config) ResolveAliasesByFirstWord(firstWord string) ([]Alias, bool) {
	if c.aliasesByFirstWord == nil {
		c.fillAliasByFirstWord()
	}
	alias, ok := c.aliasesByFirstWord[firstWord]

	return alias, ok
}

func (c *Config) fillAliasByFirstWord() {
	c.aliasesByFirstWord = make(map[string][]Alias, len(c.Aliases))
	for alias, cmd := range c.Aliases {
		if len(cmd) == 0 {
			continue
		}
		path := cmd[0]
		c.aliasesByFirstWord[path] = append(c.aliasesByFirstWord[path], Alias{
			Name:    alias,
			Command: cmd,
		})
	}
}
