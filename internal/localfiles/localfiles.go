package localfiles

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// WriteUserFileOptions contains options for writing user files
type WriteUserFileOptions struct {
	// Confirm enables interactive confirmation before writing files
	Confirm bool
	// FileMode specifies the file permissions to use when writing the file
	// If nil, defaults to 0o600
	FileMode *os.FileMode
}

// WriteUserFile writes data to a file in the user's home directory.
// It ensures the parent directory exists before writing the file.
// The function takes a homeDir parameter to specify the user's home directory,
// a filePath that is relative to the home directory, and the data and file permissions to write.
func WriteUserFile(ctx context.Context, path string, data []byte, opts *WriteUserFileOptions) error {
	// Use the provided path directly
	fullPath := path

	// Ensure the parent directory exists
	parentDir := filepath.Dir(fullPath)
	err := os.MkdirAll(parentDir, 0o755)
	if err != nil {
		return err
	}

	// Check if file exists and get current content
	existingContent, err := os.ReadFile(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// If file doesn't exist, existing content is empty
	if os.IsNotExist(err) {
		existingContent = []byte{}
	}

	// If no changes, return early
	if string(existingContent) == string(data) {
		interactive.Println("File %s is already with an identical config", path)

		return nil
	}

	// If confirmation is requested, show diff and ask for confirmation
	if opts != nil && opts.Confirm && ctx != nil {

		// Create diff
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(existingContent), string(data), false)

		// Format diff
		// Format and show diff using interactive package
		interactive.Printf("Differences in %s:\n", path)
		interactive.Print(dmp.DiffPrettyText(diffs))
		interactive.Print("\n")

		// Ask for confirmation
		confirmed, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
			Ctx:          ctx,
			Prompt:       fmt.Sprintf("Write changes to %s?", path),
			DefaultValue: true,
		})
		if err != nil {
			return err
		}
		if !confirmed {
			return fmt.Errorf("user declined to write file")
		}
	}

	// Write the file with the specified permissions
	// Use the FileMode from options if provided, otherwise default to 0o600
	fileMode := os.FileMode(0o600)
	if opts != nil && opts.FileMode != nil {
		fileMode = *opts.FileMode
	}
	err = os.WriteFile(fullPath, data, fileMode)
	if err != nil {
		return err
	}

	return nil
}

// Ask whether to remove previous configuration file if it exists
//	if _, err := os.Stat(configPath); err == nil {
//	doIt, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
//	Ctx:          ctx,
//		Prompt:       "Do you want to overwrite the existing configuration file (" + configPath + ")?",
//		DefaultValue: false,
//	})
//	if err != nil {
//		return nil, err
//	}
//	if !doIt {
//		return nil, errors.New("installation aborted by user")
//	}
//}
