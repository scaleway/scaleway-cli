package manpage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/scaleway/scaleway-cli/v2/core"
)

// Install generates the man pages for every visible command plus the
// top-level `scw` page and installs them in <target>/man1. It creates the
// target directory (and the man1 subdirectory) when missing, synchronizes the
// pages owned by the CLI (overwriting current pages and removing stale ones)
// and leaves any other file in the target untouched.
//
// It returns the number of pages installed.
func Install(ctx context.Context, commands *core.Commands, target string) (int, error) {
	sectionDir := filepath.Join(target, SectionDir)

	// Preflight: the target, when it exists, must be a directory.
	if stats, err := os.Stat(target); err == nil {
		if !stats.IsDir() {
			return 0, fmt.Errorf("target %s is not a directory", target)
		}
	} else if !os.IsNotExist(err) {
		return 0, fmt.Errorf("cannot stat target %s: %w", target, err)
	}

	// Preflight: page names must be collision free before writing anything.
	currentFiles, err := DetectNameCollisions(commands)
	if err != nil {
		return 0, err
	}

	if err := os.MkdirAll(sectionDir, 0o755); err != nil {
		return 0, fmt.Errorf("cannot create %s: %w", sectionDir, err)
	}

	// Render all pages in memory before writing anything.
	rendered, err := renderAll(ctx, commands)
	if err != nil {
		return 0, err
	}

	// Sync: remove stale pages owned by the CLI only.
	if err := syncOwnedFiles(sectionDir, currentFiles); err != nil {
		return 0, err
	}

	// Write the pages.
	for fileName, content := range rendered {
		if err := os.WriteFile(
			filepath.Join(sectionDir, fileName),
			[]byte(content),
			0o644,
		); err != nil {
			return 0, fmt.Errorf("cannot write %s: %w", filepath.Join(sectionDir, fileName), err)
		}
	}

	return len(rendered), nil
}

// renderAll renders every page (top-level + one per visible command) to roff.
func renderAll(ctx context.Context, commands *core.Commands) (map[string]string, error) {
	pages, root, err := BuildPages(ctx, commands)
	if err != nil {
		return nil, err
	}

	rendered := make(map[string]string, len(pages)+1)
	for _, page := range pages {
		content, err := Render(page)
		if err != nil {
			return nil, err
		}
		rendered[page.FileName] = content
	}

	rootContent, err := RenderRoot(root)
	if err != nil {
		return nil, err
	}
	rendered[FileName(PageName())] = rootContent

	return rendered, nil
}

// syncOwnedFiles removes files in sectionDir that are owned by the CLI
// (IsOwnedFileName) but are not part of the current command set. Files that
// are not owned by the CLI are never touched.
func syncOwnedFiles(sectionDir string, currentFiles map[string]bool) error {
	entries, err := os.ReadDir(sectionDir)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", sectionDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !IsOwnedFileName(entry.Name()) {
			continue
		}
		if currentFiles[entry.Name()] {
			continue
		}
		if err := os.Remove(filepath.Join(sectionDir, entry.Name())); err != nil {
			return fmt.Errorf(
				"cannot remove stale page %s: %w",
				filepath.Join(sectionDir, entry.Name()),
				err,
			)
		}
	}

	return nil
}
