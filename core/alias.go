package core

// aliasDisabled returns if alias are disabled for a command
func aliasDisabled(cmd string) bool {
	switch cmd {
	case "alias", "autocomplete", "init":
		return true
	}

	return false
}
