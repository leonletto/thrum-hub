package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/leonletto/thrum-hub/internal/schema"
)

func (s *Store) InboxWrite(e *schema.Entry) error {
	if e.ID == "" {
		return fmt.Errorf("inbox write: entry has empty ID")
	}
	buf, err := e.Marshal()
	if err != nil {
		return err
	}
	path := filepath.Join(inboxDir(s.root), inboxFilename(e.ID))
	return os.WriteFile(path, buf, 0644)
}

func (s *Store) InboxList() ([]string, error) {
	files, err := os.ReadDir(inboxDir(s.root))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var ids []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if !strings.HasSuffix(name, "-submitted.md") {
			continue
		}
		ids = append(ids, strings.TrimSuffix(name, "-submitted.md"))
	}
	return ids, nil
}

func (s *Store) InboxRead(id string) (*schema.Entry, error) {
	path := filepath.Join(inboxDir(s.root), inboxFilename(id))
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return schema.UnmarshalEntry(raw)
}

func (s *Store) InboxDelete(id string) error {
	path := filepath.Join(inboxDir(s.root), inboxFilename(id))
	return os.Remove(path)
}
