package interactive

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// WriteFileOptions contains options for writing user files
type WriteFileOptions struct {
	// Confirmed enables interactive confirmation before writing files
	Confirmed bool
	// FileMode specifies the file permissions to use when writing the file
	// If nil, defaults to 0o600
	FileMode *os.FileMode
}

// WriteFile writes data to a file in the user's home directory.
// It ensures the parent directory exists before writing the file.
// The function takes a homeDir parameter to specify the user's home directory,
// a filePath that is relative to the home directory, and the data and file permissions to write.
func WriteFile(
	ctx context.Context,
	path string,
	data []byte,
	opts *WriteFileOptions,
) error {
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
		Printf("File %s is already with an identical config\n", path)

		return nil
	}

	// If the command is not confirmed (for instance by using the --yes flag), show diff and ask for confirmation
	if opts != nil && !opts.Confirmed && ctx != nil {
		dmp := diffmatchpatch.New()

		// Encode texts into runes representing unique lines
		existingContentTokens, dataTokens, lineArray := dmp.DiffLinesToChars(string(existingContent), string(data))

		// Run diff on the tokenized version
		diffs := dmp.DiffMain(existingContentTokens, dataTokens, false)

		// Optional cleanup for nicer output
		dmp.DiffCleanupSemantic(diffs)

		// Convert back to text form
		diffs = dmp.DiffCharsToLines(diffs, lineArray)

		// Format diff
		// Format and show diff using interactive package
		Printf("Differences in %s:\n", path)
		Print(dmp.DiffPrettyText(diffs))
		Print("\n")

		// Ask for confirmation
		confirmed, err := PromptBoolWithConfig(&PromptBoolConfig{
			Ctx: ctx,
			Prompt: fmt.Sprintf(
				"Do you want to overwrite the existing configuration file (%s)?",
				path,
			),
			DefaultValue: false,
		})
		if err != nil {
			return err
		}
		if !confirmed {
			return fmt.Errorf("user declined to write file %s", path)
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
