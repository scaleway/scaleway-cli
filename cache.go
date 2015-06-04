package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"regexp"
	"strings"
	"sync"
)

// ScalewayCache is used not to query the API to resolve full identifiers
type ScalewayCache struct {
	// Images contains names of Scaleway images indexed by identifier
	Images map[string]string `json:"images"`

	// Snapshots contains names of Scaleway snapshots indexed by identifier
	Snapshots map[string]string `json:"snapshots"`

	// Bootscripts contains names of Scaleway bootscripts indexed by identifier
	Bootscripts map[string]string `json:"bootscripts"`

	// Servers contains names of Scaleway C1 servers indexed by identifier
	Servers map[string]string `json:"servers"`

	// Path is the path to the cache file
	Path string `json:"-"`

	// Modified tells if the cache needs to be overwritten or not
	Modified bool `json:"-"`

	// Lock allows ScalewayCache to be used concurrently
	Lock sync.Mutex `json:"-"`
}

const (
	IDENTIFIER_SERVER = iota
	IDENTIFIER_IMAGE
	IDENTIFIER_SNAPSHOT
	IDENTIFIER_BOOTSCRIPT
)

// ScalewayIdentifier is a unique identifier on Scaleway
type ScalewayIdentifier struct {
	// Identifier is a unique identifier on
	Identifier string

	// Type of the identifier
	Type int
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
			Images:      make(map[string]string),
			Snapshots:   make(map[string]string),
			Bootscripts: make(map[string]string),
			Servers:     make(map[string]string),
			Path:        cache_path,
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
	if cache.Images == nil {
		cache.Images = make(map[string]string)
	}
	if cache.Snapshots == nil {
		cache.Snapshots = make(map[string]string)
	}
	if cache.Servers == nil {
		cache.Servers = make(map[string]string)
	}
	if cache.Bootscripts == nil {
		cache.Bootscripts = make(map[string]string)
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
	needle = regexp.MustCompile(`^user/`).ReplaceAllString(needle, "")
	// FIXME: if 'user/' is in needle, only watch for a user image
	nameRegex := regexp.MustCompile(`(?i)` + regexp.MustCompile(`[_-]`).ReplaceAllString(needle, ".*"))
	for identifier, name := range c.Images {
		if strings.HasPrefix(identifier, needle) || nameRegex.MatchString(name) {
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
	needle = regexp.MustCompile(`^user/`).ReplaceAllString(needle, "")
	nameRegex := regexp.MustCompile(`(?i)` + regexp.MustCompile(`[_-]`).ReplaceAllString(needle, ".*"))
	for identifier, name := range c.Snapshots {
		if strings.HasPrefix(identifier, needle) || nameRegex.MatchString(name) {
			res = append(res, identifier)
		}
	}
	return res
}

// LookupBootscripts attempts to return identifiers matching a pattern
func (c *ScalewayCache) LookUpBootscripts(needle string) []string {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	var res []string
	nameRegex := regexp.MustCompile(`(?i)` + regexp.MustCompile(`[_-]`).ReplaceAllString(needle, ".*"))
	for identifier, name := range c.Bootscripts {
		if strings.HasPrefix(identifier, needle) || nameRegex.MatchString(name) {
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
	nameRegex := regexp.MustCompile(`(?i)` + regexp.MustCompile(`[_-]`).ReplaceAllString(needle, ".*"))
	for identifier, name := range c.Servers {
		if strings.HasPrefix(identifier, needle) || nameRegex.MatchString(name) {
			res = append(res, identifier)
		}
	}
	return res
}

// LookupIdentifier attempts to return identifiers matching a pattern
func (c *ScalewayCache) LookUpIdentifiers(needle string) []ScalewayIdentifier {
	result := []ScalewayIdentifier{}

	for _, identifier := range c.LookUpServers(needle) {
		result = append(result, ScalewayIdentifier{
			Identifier: identifier,
			Type:       IDENTIFIER_SERVER,
		})
	}

	for _, identifier := range c.LookUpImages(needle) {
		result = append(result, ScalewayIdentifier{
			Identifier: identifier,
			Type:       IDENTIFIER_IMAGE,
		})
	}

	for _, identifier := range c.LookUpSnapshots(needle) {
		result = append(result, ScalewayIdentifier{
			Identifier: identifier,
			Type:       IDENTIFIER_SNAPSHOT,
		})
	}

	for _, identifier := range c.LookUpBootscripts(needle) {
		result = append(result, ScalewayIdentifier{
			Identifier: identifier,
			Type:       IDENTIFIER_BOOTSCRIPT,
		})
	}

	return result
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

// RemoveServer removes a server from the cache
func (c *ScalewayCache) RemoveServer(identifier string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	delete(c.Servers, identifier)
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

// RemoveImage removes a server from the cache
func (c *ScalewayCache) RemoveImage(identifier string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	delete(c.Images, identifier)
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

// RemoveSnapshot removes a server from the cache
func (c *ScalewayCache) RemoveSnapshot(identifier string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	delete(c.Snapshots, identifier)
}

// InsertBootscript registers an bootscript in the cache
func (c *ScalewayCache) InsertBootscript(identifier, name string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	current_name, exists := c.Bootscripts[identifier]
	if !exists || current_name != name {
		c.Bootscripts[identifier] = name
		c.Modified = true
	}
}

// RemoveBootscript removes a bootscript from the cache
func (c *ScalewayCache) RemoveBootscript(identifier string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	delete(c.Bootscripts, identifier)
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

// GetNbBootscripts returns the number of bootscripts in the cache
func (c *ScalewayCache) GetNbBootscripts() int {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	return len(c.Bootscripts)
}
