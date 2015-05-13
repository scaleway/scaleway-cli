package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

// ScalewayCache is used not to query the API to resolve full identifiers
type ScalewayCache struct {
	// Images contains names of Scaleway images indexed by identifier
	Images map[string]string `json:"images"`

	// Servers contains names of Scaleway C1 servers indexed by identifier
	Servers map[string]string `json:"servers"`

	// Path is the path to the cache file
	Path string `json:"-"`

	// Modified tells if the cache needs to be overwritten or not
	Modified bool `json:"-"`
}

// NewScalewayCache loads a per-user cache
func NewScalewayCache() (*ScalewayCache, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	cache_path := fmt.Sprintf("%s/.scw-cache.db", u.HomeDir)
	_, err = os.Stat(cache_path)
	if os.IsNotExist(err) {
		return &ScalewayCache{
			Images:  make(map[string]string),
			Servers: make(map[string]string),
			Path:    cache_path,
		}, nil
	} else if err != nil {
		return nil, err
	}
	file, err := ioutil.ReadFile(cache_path)
	if err != nil {
		return nil, err
	}
	var cache ScalewayCache
	cache.Path = cache_path
	err = json.Unmarshal(file, &cache)
	if err != nil {
		return nil, err
	}
	return &cache, nil
}

// Save atomically overwrites the current cache database
func (c *ScalewayCache) Save() error {
	if c.Modified {
		file, err := ioutil.TempFile("", "")
		if err != nil {
			return err
		}
		encoder := json.NewEncoder(file)
		err = encoder.Encode(*c)
		if err != nil {
			return err
		}
		return os.Rename(file.Name(), c.Path)
	}
	return nil
}

// LookupImages attempts to return identifiers matching a pattern
func (c *ScalewayCache) LookUpImages(needle string) []string {
	var res []string
	for identifier, name := range c.Images {
		if strings.HasPrefix(identifier, needle) || strings.HasPrefix(name, needle) {
			res = append(res, identifier)
		}
	}
	return res
}

// LookupServers attempts to return identifiers matching a pattern
func (c *ScalewayCache) LookUpServers(needle string) []string {
	var res []string
	for identifier, name := range c.Servers {
		if strings.HasPrefix(identifier, needle) || strings.HasPrefix(name, needle) {
			res = append(res, identifier)
		}
	}
	return res
}

// InsertServer registers a server in the cache
func (c *ScalewayCache) InsertServer(identifier, name string) {
	current_name, exists := c.Servers[identifier]
	if !exists || current_name != name {
		c.Servers[identifier] = name
		c.Modified = true
	}
}

// InsertImage registers an image in the cache
func (c *ScalewayCache) InsertImage(identifier, name string) {
	current_name, exists := c.Images[identifier]
	if !exists || current_name != name {
		c.Images[identifier] = name
		c.Modified = true
	}
}
