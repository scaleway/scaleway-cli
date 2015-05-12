package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

// Cache is used not to query the API to resolve full identifiers
type Cache struct {
	// Images contains names of Scaleway images indexed by identifier
	Images map[string]string `json:"images"`

	// Servers contains names of Scaleway C1 servers indexed by identifier
	Servers map[string]string `json:"servers"`

	// Path is the path to the cache file
	Path string
}

// NewCache loads a per-user cache
func NewCache() (*Cache, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	cache_path := fmt.Sprintf("%s/.scw-cache.db", u.HomeDir)
	_, err = os.Stat(cache_path)
	if os.IsNotExist(err) {
		return &Cache{
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
	var cache Cache
	cache.Path = cache_path
	err = json.Unmarshal(file, &cache)
	if err != nil {
		return nil, err
	}
	return &cache, nil
}

// Save atomically overwrites the current cache database
func (c *Cache) Save() error {
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
