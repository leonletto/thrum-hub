package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/leonletto/thrum-hub/internal/schema"
)

func newTestEntry(id string) *schema.Entry {
	return &schema.Entry{
		ID:       id,
		Type:     schema.TypeGotcha,
		Title:    "test gotcha",
		Language: "go",
		Status:   schema.StatusActive,
		Created:  time.Now().UTC().Truncate(time.Second),
		Content:  "body goes here",
	}
}

func TestInboxRoundTrip(t *testing.T) {
	root := t.TempDir()
	s, err := New(root)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	e := newTestEntry("khub_01HZINBOXTESTID0000000000")
	if err := s.InboxWrite(e); err != nil {
		t.Fatalf("InboxWrite: %v", err)
	}
	ids, err := s.InboxList()
	if err != nil {
		t.Fatalf("InboxList: %v", err)
	}
	if len(ids) != 1 || ids[0] != e.ID {
		t.Fatalf("InboxList got %v, want [%q]", ids, e.ID)
	}
	got, err := s.InboxRead(e.ID)
	if err != nil {
		t.Fatalf("InboxRead: %v", err)
	}
	if got.Title != e.Title {
		t.Errorf("Title mismatch: got %q", got.Title)
	}
	if err := s.InboxDelete(e.ID); err != nil {
		t.Fatalf("InboxDelete: %v", err)
	}
	ids, _ = s.InboxList()
	if len(ids) != 0 {
		t.Fatalf("expected empty inbox after delete, got %v", ids)
	}
}

func TestEntryWriteAndRead(t *testing.T) {
	root := t.TempDir()
	s, _ := New(root)
	e := newTestEntry("khub_01HZENTRYTESTID0000000000")
	if err := s.EntryWrite(e); err != nil {
		t.Fatalf("EntryWrite: %v", err)
	}
	path := s.EntryPath(e.ID)
	if !filepath.IsAbs(path) {
		t.Errorf("EntryPath not absolute: %q", path)
	}
	got, err := s.EntryRead(e.ID)
	if err != nil {
		t.Fatalf("EntryRead: %v", err)
	}
	if got.ID != e.ID {
		t.Errorf("ID mismatch")
	}
}

func TestArchiveSupersede(t *testing.T) {
	root := t.TempDir()
	s, _ := New(root)
	e := newTestEntry("khub_01HZARCHIVETESTID000000000")
	_ = s.EntryWrite(e)
	if err := s.ArchiveSupersede(e.ID); err != nil {
		t.Fatalf("ArchiveSupersede: %v", err)
	}
	if _, err := s.EntryRead(e.ID); err == nil {
		t.Errorf("expected EntryRead to fail after archive")
	}
	archived, err := s.ArchiveRead(e.ID)
	if err != nil {
		t.Fatalf("ArchiveRead: %v", err)
	}
	if archived.Status != schema.StatusSuperseded {
		t.Errorf("expected StatusSuperseded, got %q", archived.Status)
	}
}
