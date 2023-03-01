package core

// aliasDisabled returns if alias are disabled for a command
func aliasDisabled(cmd string) bool {
	switch cmd {
	case "alias":
		return true
	case "autocomplete":
		return true
	}
	return false
}
