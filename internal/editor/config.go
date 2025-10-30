package editor

import (
	"os"
	"runtime"
)

// List of env variables where to find the editor to use
// Order in slice is override order, the latest will override the first ones
var editorEnvVariables = []string{"EDITOR", "VISUAL"}

func GetSystemDefaultEditor() string {
	switch runtime.GOOS {
	case "windows":
		return "notepad"
	default:
		return "vi"
	}
}

func GetDefaultEditor() string {
	editor := ""
	for _, envVar := range editorEnvVariables {
		tmp := os.Getenv(envVar)
		if tmp != "" {
			editor = tmp
		}
	}

	if editor == "" {
		return GetSystemDefaultEditor()
	}

	return editor
}
