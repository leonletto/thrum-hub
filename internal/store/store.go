package store

import (
	"fmt"
	"os"
)

type Store struct {
	root string
}

func New(root string) (*Store, error) {
	for _, d := range []string{inboxDir(root), entriesDir(root), archiveDir(root)} {
		if err := os.MkdirAll(d, 0755); err != nil {
			return nil, fmt.Errorf("mkdir %s: %w", d, err)
		}
	}
	return &Store{root: root}, nil
}

func (s *Store) Root() string { return s.root }
