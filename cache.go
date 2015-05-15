package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"sync"
)

// ScalewayCache is used not to query the API to resolve full identifiers
type ScalewayCache struct {
	// Images contains names of Scaleway images indexed by identifier
	Images map[string]string `json:"images"`

	// Snapshots contains names of Scaleway snapshots indexed by identifier
	Snapshots map[string]string `json:"snapshots"`

	// Servers contains names of Scaleway C1 servers indexed by identifier
	Servers map[string]string `json:"servers"`

	// Path is the path to the cache file
	Path string `json:"-"`

	// Modified tells if the cache needs to be overwritten or not
	Modified bool `json:"-"`

	// Lock allows ScalewayCache to be used concurrently
	Lock sync.Mutex `json:"-"`
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
			Images:    make(map[string]string),
			Snapshots: make(map[string]string),
			Servers:   make(map[string]string),
			Path:      cache_path,
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
	c.Lock.Lock()
	defer c.Lock.Unlock()

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
	c.Lock.Lock()
	defer c.Lock.Unlock()

	var res []string
	for identifier, name := range c.Images {
		if strings.HasPrefix(identifier, needle) || strings.HasPrefix(name, needle) {
			res = append(res, identifier)
		}
	}
	return res
}

// LookupSnapshots attempts to return identifiers matching a pattern
func (c *ScalewayCache) LookUpSnapshots(needle string) []string {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	var res []string
	for identifier, name := range c.Snapshots {
		if strings.HasPrefix(identifier, needle) || strings.HasPrefix(name, needle) {
			res = append(res, identifier)
		}
	}
	return res
}

// LookupServers attempts to return identifiers matching a pattern
func (c *ScalewayCache) LookUpServers(needle string) []string {
	c.Lock.Lock()
	defer c.Lock.Unlock()

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
	c.Lock.Lock()
	defer c.Lock.Unlock()

	current_name, exists := c.Servers[identifier]
	if !exists || current_name != name {
		c.Servers[identifier] = name
		c.Modified = true
	}
}

// InsertImage registers an image in the cache
func (c *ScalewayCache) InsertImage(identifier, name string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	current_name, exists := c.Images[identifier]
	if !exists || current_name != name {
		c.Images[identifier] = name
		c.Modified = true
	}
}

// InsertSnapshot registers an snapshot in the cache
func (c *ScalewayCache) InsertSnapshot(identifier, name string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	current_name, exists := c.Snapshots[identifier]
	if !exists || current_name != name {
		c.Snapshots[identifier] = name
		c.Modified = true
	}
}

// GetNbServers returns the number of servers in the cache
func (c *ScalewayCache) GetNbServers() int {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	return len(c.Servers)
}

// GetNbImages returns the number of images in the cache
func (c *ScalewayCache) GetNbImages() int {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	return len(c.Images)
}

// GetNbSnapshots returns the number of snapshots in the cache
func (c *ScalewayCache) GetNbSnapshots() int {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	return len(c.Snapshots)
}
