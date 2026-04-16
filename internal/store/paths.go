package store

import (
	"fmt"
	"path/filepath"
)

func inboxDir(root string) string   { return filepath.Join(root, "inbox") }
func entriesDir(root string) string { return filepath.Join(root, "entries") }
func archiveDir(root string) string { return filepath.Join(root, "archive", "superseded") }

func inboxFilename(id string) string {
	return fmt.Sprintf("%s-submitted.md", id)
}

func entryFilename(id string) string {
	return fmt.Sprintf("%s.md", id)
}
