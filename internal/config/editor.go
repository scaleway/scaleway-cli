package config

import "os"

var (
	// List of env variables where to find the editor to use
	// Order in slice is override order, the latest will override the first ones
	editorEnvVariables = []string{"EDITOR", "VISUAL"}
)

func GetDefaultEditor() string {
	editor := ""
	for _, envVar := range editorEnvVariables {
		tmp := os.Getenv(envVar)
		if tmp != "" {
			editor = tmp
		}
	}

	if editor == "" {
		return "vi"
	}

	return editor
}
