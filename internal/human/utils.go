package human

import "unicode"

// Capitalize returns the given string with a first character in uppercase.
func Capitalize(s string) string {
	for i, c := range s {
		return string(unicode.ToUpper(c)) + s[i+1:]
	}
	return ""
}
