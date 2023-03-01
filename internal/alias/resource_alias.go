package alias

import (
	"strings"
)

func JoinResourcePath(path []string) string {
	return strings.Join(path, ".")
}

// SplitResourcePath split a resource from a resource alias to an array of scw words
func SplitResourcePath(resource string) []string {
	return strings.Split(resource, ".")
}

func (c *Config) AddResourceAlias(path []string, alias string) {
	joinedPath := JoinResourcePath(path)

	c.ResourceAliases[joinedPath] = append(c.ResourceAliases[joinedPath], alias)
}

// DeleteResourceAlias delete an alias in resource with given path
// return true if alias has been deleted
func (c *Config) DeleteResourceAlias(path []string, alias string) bool {
	resourcePath := JoinResourcePath(path)
	aliases, exists := c.ResourceAliases[resourcePath]
	if !exists {
		return false
	}
	aliasIndex := -1
	for i, a := range aliases {
		if a == alias {
			aliasIndex = i
			break
		}
	}
	if aliasIndex == -1 {
		return false
	}
	aliases = append(aliases[:aliasIndex], aliases[aliasIndex+1:]...)
	if len(aliases) == 0 {
		delete(c.ResourceAliases, resourcePath)
	} else {
		c.ResourceAliases[resourcePath] = aliases
	}
	return true
}
