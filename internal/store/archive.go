package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/leonletto/thrum-hub/internal/schema"
)

func (s *Store) ArchiveSupersede(id string) error {
	e, err := s.EntryRead(id)
	if err != nil {
		return err
	}
	oldPath := s.EntryPath(id)
	if oldPath == "" {
		return fmt.Errorf("archive: entry not found: %s", id)
	}
	e.Status = schema.StatusSuperseded
	buf, err := e.Marshal()
	if err != nil {
		return err
	}
	dst := filepath.Join(archiveDir(s.root), entryFilename(id))
	if err := os.WriteFile(dst, buf, 0644); err != nil {
		return err
	}
	return os.Remove(oldPath)
}

func (s *Store) ArchiveRead(id string) (*schema.Entry, error) {
	path := filepath.Join(archiveDir(s.root), entryFilename(id))
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return schema.UnmarshalEntry(raw)
}
