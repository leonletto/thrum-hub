package store

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/leonletto/thrum-hub/internal/schema"
)

// EntryPath returns the absolute file path for a curated entry,
// searching under entries/<language>/<type>/<id>.md. Returns empty
// string if the entry doesn't exist.
func (s *Store) EntryPath(id string) string {
	root := entriesDir(s.root)
	var found string
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, entryFilename(id)) {
			found = path
			return filepath.SkipAll
		}
		return nil
	})
	return found
}

func (s *Store) EntryWrite(e *schema.Entry) error {
	if e.ID == "" || e.Language == "" || !e.Type.Valid() {
		return fmt.Errorf("entry write: id, language, and valid type are required")
	}
	dir := filepath.Join(entriesDir(s.root), e.Language, string(e.Type))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	buf, err := e.Marshal()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, entryFilename(e.ID)), buf, 0644)
}

func (s *Store) EntryRead(id string) (*schema.Entry, error) {
	path := s.EntryPath(id)
	if path == "" {
		return nil, fmt.Errorf("entry not found: %s", id)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return schema.UnmarshalEntry(raw)
}

type EntryFilter struct {
	Language string
	Type     schema.EntryType
	Tag      string
}

func (s *Store) EntryList(filter EntryFilter) ([]*schema.Entry, error) {
	root := entriesDir(s.root)
	var out []*schema.Entry
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		e, err := schema.UnmarshalEntry(raw)
		if err != nil {
			return nil // skip malformed
		}
		if filter.Language != "" && e.Language != filter.Language {
			return nil
		}
		if filter.Type != "" && e.Type != filter.Type {
			return nil
		}
		if filter.Tag != "" && !containsTag(e.Tags, filter.Tag) {
			return nil
		}
		out = append(out, e)
		return nil
	})
	return out, err
}

func containsTag(tags []string, want string) bool {
	for _, t := range tags {
		if t == want {
			return true
		}
	}
	return false
}
