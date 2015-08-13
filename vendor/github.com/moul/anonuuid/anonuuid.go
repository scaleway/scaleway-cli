package anonuuid

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// UUIDRegex is the regex used to find UUIDs in texts
	UUIDRegex = "[a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12}"
)

// AnonUUID is the main structure, it contains the cache map and helpers
type AnonUUID struct {
	cache map[string]string
}

// Sanitize takes a string as input and return sanitized string
func (a *AnonUUID) Sanitize(input string) string {
	r := regexp.MustCompile(UUIDRegex)

	return r.ReplaceAllStringFunc(input, func(m string) string {
		parts := r.FindStringSubmatch(m)
		return a.FakeUUID(parts[0])
	})
}

// FakeUUID takes a realUUID and return its corresponding fakeUUID
func (a *AnonUUID) FakeUUID(realUUID string) string {
	if _, ok := a.cache[realUUID]; !ok {
		nextID := len(a.cache)
		fakeUUID := strings.Repeat(fmt.Sprintf("%x", nextID), 32)[:32]
		fakeUUID = fakeUUID[:8] + "-" + fakeUUID[8:12] + "-" + fakeUUID[12:16] + "-" + fakeUUID[16:20] + "-" + fakeUUID[20:32]
		a.cache[realUUID] = fakeUUID
	}
	return a.cache[realUUID]
}

// New returns a prepared AnonUUID structure
func New() *AnonUUID {
	return &AnonUUID{
		cache: make(map[string]string),
	}
}
