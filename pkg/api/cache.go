// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/scaleway/scaleway-cli/vendor/code.google.com/p/go-uuid/uuid"
	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/vendor/github.com/renstrom/fuzzysearch/fuzzy"
)

// ScalewayCache is used not to query the API to resolve full identifiers
type ScalewayCache struct {
	// Images contains names of Scaleway images indexed by identifier
	Images map[string]string `json:"images"`

	// Snapshots contains names of Scaleway snapshots indexed by identifier
	Snapshots map[string]string `json:"snapshots"`

	// Volumes contains names of Scaleway volumes indexed by identifier
	Volumes map[string]string `json:"volumes"`

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
	// IdentifierUnknown is used when we don't know explicitely the type key of the object (used for nil comparison)
	IdentifierUnknown = 1 << iota
	// IdentifierServer is the type key of cached server objects
	IdentifierServer
	// IdentifierImage is the type key of cached image objects
	IdentifierImage
	// IdentifierSnapshot is the type key of cached snapshot objects
	IdentifierSnapshot
	// IdentifierBootscript is the type key of cached bootscript objects
	IdentifierBootscript
	// IdentifierVolume is the type key of cached volume objects
	IdentifierVolume
)

// ScalewayResolverResult is a structure containing human-readable information
// about resolver results. This structure is used to display the user choices.
type ScalewayResolverResult struct {
	Identifier string
	Type       int
	Name       string
	Needle     string
	RankMatch  int
}

// ScalewayResolverResults is a list of `ScalewayResolverResult`
type ScalewayResolverResults []ScalewayResolverResult

func (s ScalewayResolverResults) Len() int {
	return len(s)
}

func (s ScalewayResolverResults) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ScalewayResolverResults) Less(i, j int) bool {
	return s[i].RankMatch < s[j].RankMatch
}

// TruncIdentifier returns first 8 characters of an Identifier (UUID)
func (s *ScalewayResolverResult) TruncIdentifier() string {
	return s.Identifier[:8]
}

func identifierTypeName(kind int) string {
	switch kind {
	case IdentifierServer:
		return "Server"
	case IdentifierImage:
		return "Image"
	case IdentifierSnapshot:
		return "Snapshot"
	case IdentifierVolume:
		return "Volume"
	case IdentifierBootscript:
		return "Bootscript"
	}
	return ""
}

// CodeName returns a full resource name with typed prefix
func (s *ScalewayResolverResult) CodeName() string {
	name := strings.ToLower(s.Name)
	name = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(name, "-")
	name = regexp.MustCompile(`--+`).ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")

	return fmt.Sprintf("%s:%s", strings.ToLower(identifierTypeName(s.Type)), name)
}

// NewScalewayCache loads a per-user cache
func NewScalewayCache() (*ScalewayCache, error) {
	homeDir := os.Getenv("HOME") // *nix
	if homeDir == "" {           // Windows
		homeDir = os.Getenv("USERPROFILE")
	}
	if homeDir == "" {
		homeDir = "/tmp"
	}
	cachePath := filepath.Join(homeDir, ".scw-cache.db")
	_, err := os.Stat(cachePath)
	if os.IsNotExist(err) {
		return &ScalewayCache{
			Images:      make(map[string]string),
			Snapshots:   make(map[string]string),
			Volumes:     make(map[string]string),
			Bootscripts: make(map[string]string),
			Servers:     make(map[string]string),
			Path:        cachePath,
		}, nil
	} else if err != nil {
		return nil, err
	}
	file, err := ioutil.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}
	var cache ScalewayCache
	cache.Path = cachePath
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
	if cache.Volumes == nil {
		cache.Volumes = make(map[string]string)
	}
	if cache.Servers == nil {
		cache.Servers = make(map[string]string)
	}
	if cache.Bootscripts == nil {
		cache.Bootscripts = make(map[string]string)
	}
	return &cache, nil
}

// Flush flushes the cache database
func (c *ScalewayCache) Flush() error {
	return os.Remove(c.Path)
}

// Save atomically overwrites the current cache database
func (c *ScalewayCache) Save() error {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	logrus.Debugf("Writing cache file to disk")

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

// ComputeRankMatch fills `ScalewayResolverResult.RankMatch` with its `fuzzy` score
func (s *ScalewayResolverResult) ComputeRankMatch(needle string) {
	s.Needle = needle
	s.RankMatch = fuzzy.RankMatch(needle, s.Name)
}

// LookUpImages attempts to return identifiers matching a pattern
func (c *ScalewayCache) LookUpImages(needle string, acceptUUID bool) ScalewayResolverResults {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	var res ScalewayResolverResults
	var exactMatches ScalewayResolverResults

	if acceptUUID && uuid.Parse(needle) != nil {
		entry := ScalewayResolverResult{
			Identifier: needle,
			Name:       needle,
			Type:       IdentifierImage,
		}
		entry.ComputeRankMatch(needle)
		res = append(res, entry)
	}

	needle = regexp.MustCompile(`^user/`).ReplaceAllString(needle, "")
	// FIXME: if 'user/' is in needle, only watch for a user image
	nameRegex := regexp.MustCompile(`(?i)` + regexp.MustCompile(`[_-]`).ReplaceAllString(needle, ".*"))
	for identifier, name := range c.Images {
		if name == needle {
			entry := ScalewayResolverResult{
				Identifier: identifier,
				Name:       name,
				Type:       IdentifierImage,
			}
			entry.ComputeRankMatch(needle)
			exactMatches = append(exactMatches, entry)
		}
		if strings.HasPrefix(identifier, needle) || nameRegex.MatchString(name) {
			entry := ScalewayResolverResult{
				Identifier: identifier,
				Name:       name,
				Type:       IdentifierImage,
			}
			entry.ComputeRankMatch(needle)
			res = append(res, entry)
		}
	}

	if len(exactMatches) == 1 {
		return exactMatches
	}

	return removeDuplicatesResults(res)
}

// LookUpSnapshots attempts to return identifiers matching a pattern
func (c *ScalewayCache) LookUpSnapshots(needle string, acceptUUID bool) ScalewayResolverResults {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	var res ScalewayResolverResults
	var exactMatches ScalewayResolverResults

	if acceptUUID && uuid.Parse(needle) != nil {
		entry := ScalewayResolverResult{
			Identifier: needle,
			Name:       needle,
			Type:       IdentifierSnapshot,
		}
		entry.ComputeRankMatch(needle)
		res = append(res, entry)
	}

	needle = regexp.MustCompile(`^user/`).ReplaceAllString(needle, "")
	nameRegex := regexp.MustCompile(`(?i)` + regexp.MustCompile(`[_-]`).ReplaceAllString(needle, ".*"))
	for identifier, name := range c.Snapshots {
		if name == needle {
			entry := ScalewayResolverResult{
				Identifier: identifier,
				Name:       name,
				Type:       IdentifierSnapshot,
			}
			entry.ComputeRankMatch(needle)
			exactMatches = append(exactMatches, entry)
		}
		if strings.HasPrefix(identifier, needle) || nameRegex.MatchString(name) {
			entry := ScalewayResolverResult{
				Identifier: identifier,
				Name:       name,
				Type:       IdentifierSnapshot,
			}
			entry.ComputeRankMatch(needle)
			res = append(res, entry)
		}
	}

	if len(exactMatches) == 1 {
		return exactMatches
	}

	return removeDuplicatesResults(res)
}

// LookUpVolumes attempts to return identifiers matching a pattern
func (c *ScalewayCache) LookUpVolumes(needle string, acceptUUID bool) ScalewayResolverResults {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	var res ScalewayResolverResults
	var exactMatches ScalewayResolverResults

	if acceptUUID && uuid.Parse(needle) != nil {
		entry := ScalewayResolverResult{
			Identifier: needle,
			Name:       needle,
			Type:       IdentifierVolume,
		}
		entry.ComputeRankMatch(needle)
		res = append(res, entry)
	}

	nameRegex := regexp.MustCompile(`(?i)` + regexp.MustCompile(`[_-]`).ReplaceAllString(needle, ".*"))
	for identifier, name := range c.Volumes {
		if name == needle {
			entry := ScalewayResolverResult{
				Identifier: identifier,
				Name:       name,
				Type:       IdentifierVolume,
			}
			entry.ComputeRankMatch(needle)
			exactMatches = append(exactMatches, entry)
		}
		if strings.HasPrefix(identifier, needle) || nameRegex.MatchString(name) {
			entry := ScalewayResolverResult{
				Identifier: identifier,
				Name:       name,
				Type:       IdentifierVolume,
			}
			entry.ComputeRankMatch(needle)
			res = append(res, entry)
		}
	}

	if len(exactMatches) == 1 {
		return exactMatches
	}

	return removeDuplicatesResults(res)
}

// LookUpBootscripts attempts to return identifiers matching a pattern
func (c *ScalewayCache) LookUpBootscripts(needle string, acceptUUID bool) ScalewayResolverResults {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	var res ScalewayResolverResults
	var exactMatches ScalewayResolverResults

	if acceptUUID && uuid.Parse(needle) != nil {
		entry := ScalewayResolverResult{
			Identifier: needle,
			Name:       needle,
			Type:       IdentifierBootscript,
		}
		entry.ComputeRankMatch(needle)
		res = append(res, entry)
	}

	nameRegex := regexp.MustCompile(`(?i)` + regexp.MustCompile(`[_-]`).ReplaceAllString(needle, ".*"))
	for identifier, name := range c.Bootscripts {
		if name == needle {
			entry := ScalewayResolverResult{
				Identifier: identifier,
				Name:       name,
				Type:       IdentifierBootscript,
			}
			entry.ComputeRankMatch(needle)
			exactMatches = append(exactMatches, entry)
		}
		if strings.HasPrefix(identifier, needle) || nameRegex.MatchString(name) {
			entry := ScalewayResolverResult{
				Identifier: identifier,
				Name:       name,
				Type:       IdentifierBootscript,
			}
			entry.ComputeRankMatch(needle)
			res = append(res, entry)
		}
	}

	if len(exactMatches) == 1 {
		return exactMatches
	}

	return removeDuplicatesResults(res)
}

// LookUpServers attempts to return identifiers matching a pattern
func (c *ScalewayCache) LookUpServers(needle string, acceptUUID bool) ScalewayResolverResults {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	var res ScalewayResolverResults
	var exactMatches ScalewayResolverResults

	if acceptUUID && uuid.Parse(needle) != nil {
		entry := ScalewayResolverResult{
			Identifier: needle,
			Name:       needle,
			Type:       IdentifierServer,
		}
		entry.ComputeRankMatch(needle)
		res = append(res, entry)
	}

	nameRegex := regexp.MustCompile(`(?i)` + regexp.MustCompile(`[_-]`).ReplaceAllString(needle, ".*"))
	for identifier, name := range c.Servers {
		if name == needle {
			entry := ScalewayResolverResult{
				Identifier: identifier,
				Name:       name,
				Type:       IdentifierServer,
			}
			entry.ComputeRankMatch(needle)
			exactMatches = append(exactMatches, entry)
		}
		if strings.HasPrefix(identifier, needle) || nameRegex.MatchString(name) {
			entry := ScalewayResolverResult{
				Identifier: identifier,
				Name:       name,
				Type:       IdentifierServer,
			}
			entry.ComputeRankMatch(needle)
			res = append(res, entry)
		}
	}

	if len(exactMatches) == 1 {
		return exactMatches
	}

	return removeDuplicatesResults(res)
}

// removeDuplicatesResults transforms an array into a unique array
func removeDuplicatesResults(elements ScalewayResolverResults) ScalewayResolverResults {
	encountered := map[string]ScalewayResolverResult{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v].Identifier] = elements[v]
	}

	// Place all keys from the map into a slice.
	results := ScalewayResolverResults{}
	for _, result := range encountered {
		results = append(results, result)
	}
	return results
}

// parseNeedle parses a user needle and try to extract a forced object type
// i.e:
//   - server:blah-blah -> kind=server, needle=blah-blah
//   - blah-blah -> kind="", needle=blah-blah
//   - not-existing-type:blah-blah
func parseNeedle(input string) (identifierType int, needle string) {
	parts := strings.Split(input, ":")
	if len(parts) == 2 {
		switch parts[0] {
		case "server":
			return IdentifierServer, parts[1]
		case "image":
			return IdentifierImage, parts[1]
		case "snapshot":
			return IdentifierSnapshot, parts[1]
		case "bootscript":
			return IdentifierBootscript, parts[1]
		case "volume":
			return IdentifierVolume, parts[1]
		}
	}
	return IdentifierUnknown, input
}

// LookUpIdentifiers attempts to return identifiers matching a pattern
func (c *ScalewayCache) LookUpIdentifiers(needle string) ScalewayResolverResults {
	results := ScalewayResolverResults{}

	identifierType, needle := parseNeedle(needle)

	if identifierType&(IdentifierUnknown|IdentifierServer) > 0 {
		for _, result := range c.LookUpServers(needle, false) {
			entry := ScalewayResolverResult{
				Identifier: result.Identifier,
				Name:       result.Name,
				Type:       IdentifierServer,
			}
			entry.ComputeRankMatch(needle)
			results = append(results, entry)
		}
	}

	if identifierType&(IdentifierUnknown|IdentifierImage) > 0 {
		for _, result := range c.LookUpImages(needle, false) {
			entry := ScalewayResolverResult{
				Identifier: result.Identifier,
				Name:       result.Name,
				Type:       IdentifierImage,
			}
			entry.ComputeRankMatch(needle)
			results = append(results, entry)
		}
	}

	if identifierType&(IdentifierUnknown|IdentifierSnapshot) > 0 {
		for _, result := range c.LookUpSnapshots(needle, false) {
			entry := ScalewayResolverResult{
				Identifier: result.Identifier,
				Name:       result.Name,
				Type:       IdentifierSnapshot,
			}
			entry.ComputeRankMatch(needle)
			results = append(results, entry)
		}
	}

	if identifierType&(IdentifierUnknown|IdentifierVolume) > 0 {
		for _, result := range c.LookUpVolumes(needle, false) {
			entry := ScalewayResolverResult{
				Identifier: result.Identifier,
				Name:       result.Name,
				Type:       IdentifierVolume,
			}
			entry.ComputeRankMatch(needle)
			results = append(results, entry)
		}
	}

	if identifierType&(IdentifierUnknown|IdentifierBootscript) > 0 {
		for _, result := range c.LookUpBootscripts(needle, false) {
			entry := ScalewayResolverResult{
				Identifier: result.Identifier,
				Name:       result.Name,
				Type:       IdentifierBootscript,
			}
			entry.ComputeRankMatch(needle)
			results = append(results, entry)
		}
	}

	return results
}

// InsertServer registers a server in the cache
func (c *ScalewayCache) InsertServer(identifier, name string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	currentName, exists := c.Servers[identifier]
	if !exists || currentName != name {
		c.Servers[identifier] = name
		c.Modified = true
	}
}

// RemoveServer removes a server from the cache
func (c *ScalewayCache) RemoveServer(identifier string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	delete(c.Servers, identifier)
	c.Modified = true
}

// ClearServers removes all servers from the cache
func (c *ScalewayCache) ClearServers() {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.Servers = make(map[string]string)
	c.Modified = true
}

// InsertImage registers an image in the cache
func (c *ScalewayCache) InsertImage(identifier, name string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	currentName, exists := c.Images[identifier]
	if !exists || currentName != name {
		c.Images[identifier] = name
		c.Modified = true
	}
}

// RemoveImage removes a server from the cache
func (c *ScalewayCache) RemoveImage(identifier string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	delete(c.Images, identifier)
	c.Modified = true
}

// ClearImages removes all images from the cache
func (c *ScalewayCache) ClearImages() {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.Images = make(map[string]string)
	c.Modified = true
}

// InsertSnapshot registers an snapshot in the cache
func (c *ScalewayCache) InsertSnapshot(identifier, name string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	currentName, exists := c.Snapshots[identifier]
	if !exists || currentName != name {
		c.Snapshots[identifier] = name
		c.Modified = true
	}
}

// RemoveSnapshot removes a server from the cache
func (c *ScalewayCache) RemoveSnapshot(identifier string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	delete(c.Snapshots, identifier)
	c.Modified = true
}

// ClearSnapshots removes all snapshots from the cache
func (c *ScalewayCache) ClearSnapshots() {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.Snapshots = make(map[string]string)
	c.Modified = true
}

// InsertVolume registers an volume in the cache
func (c *ScalewayCache) InsertVolume(identifier, name string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	currentName, exists := c.Volumes[identifier]
	if !exists || currentName != name {
		c.Volumes[identifier] = name
		c.Modified = true
	}
}

// RemoveVolume removes a server from the cache
func (c *ScalewayCache) RemoveVolume(identifier string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	delete(c.Volumes, identifier)
	c.Modified = true
}

// ClearVolumes removes all volumes from the cache
func (c *ScalewayCache) ClearVolumes() {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.Volumes = make(map[string]string)
	c.Modified = true
}

// InsertBootscript registers an bootscript in the cache
func (c *ScalewayCache) InsertBootscript(identifier, name string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	currentName, exists := c.Bootscripts[identifier]
	if !exists || currentName != name {
		c.Bootscripts[identifier] = name
		c.Modified = true
	}
}

// RemoveBootscript removes a bootscript from the cache
func (c *ScalewayCache) RemoveBootscript(identifier string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	delete(c.Bootscripts, identifier)
	c.Modified = true
}

// ClearBootscripts removes all bootscripts from the cache
func (c *ScalewayCache) ClearBootscripts() {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.Bootscripts = make(map[string]string)
	c.Modified = true
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

// GetNbVolumes returns the number of volumes in the cache
func (c *ScalewayCache) GetNbVolumes() int {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	return len(c.Volumes)
}

// GetNbBootscripts returns the number of bootscripts in the cache
func (c *ScalewayCache) GetNbBootscripts() int {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	return len(c.Bootscripts)
}
