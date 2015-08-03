package anonuuid

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var (
	// UUIDRegex is the regex used to find UUIDs in texts
	UUIDRegex = "[a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12}"
)

// AnonUUID is the main structure, it contains the cache map and helpers
type AnonUUID struct {
	cache map[string]string

	// Hexspeak flag will generate hexspeak style fake UUIDs
	Hexspeak bool

	// Random flag will generate random fake UUIDs
	Random bool
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

		var fakeUUID string
		if a.Hexspeak {
			fakeUUID = GenerateHexspeakUUID(len(a.cache))
		} else if a.Random {
			fakeUUID = GenerateRandomUUID(10)
		} else {
			fakeUUID = GenerateLenUUID(len(a.cache))
		}

		// FIXME: check for duplicates and retry

		a.cache[realUUID] = fakeUUID
	}
	return a.cache[realUUID]
}

// New returns a prepared AnonUUID structure
func New() *AnonUUID {
	return &AnonUUID{
		cache:    make(map[string]string),
		Hexspeak: false,
		Random:   false,
	}
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// FormatUUID takes a string in input and return an UUID formatted string by repeating the string and placing dashes if necessary
func FormatUUID(part string) string {
	if len(part) < 32 {
		part = strings.Repeat(part, 32)
	}
	if len(part) > 32 {
		part = part[:32]
	}
	return part[:8] + "-" + part[8:12] + "-" + part[12:16] + "-" + part[16:20] + "-" + part[20:32]
}

// GenerateRandomUUID returns an UUID based on random strings
func GenerateRandomUUID(length int) string {
	var letters = []rune("abcdef0123456789")

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return FormatUUID(string(b))
}

// GenerateHexspeakUUID returns an UUID formatted string containing hexspeak words
func GenerateHexspeakUUID(i int) string {
	hexspeaks := []string{
		"0ff1ce",
		"31337",
		"4b1d",
		"badc0de",
		"badcafe",
		"badf00d",
		"deadbabe",
		"deadbeef",
		"deadc0de",
		"deadfeed",
		"fee1bad",
	}
	return FormatUUID(hexspeaks[i%len(hexspeaks)])
}

// GenerateLenUUID returns an UUID formatted string based on an index number
func GenerateLenUUID(i int) string {
	return FormatUUID(fmt.Sprintf("%x", i))
}
