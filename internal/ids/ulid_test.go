package ids

import (
	"strings"
	"testing"
)

func TestNewEntryID(t *testing.T) {
	a := NewEntryID()
	b := NewEntryID()
	if a == b {
		t.Fatalf("expected distinct IDs, got %q both", a)
	}
	if !strings.HasPrefix(a, "khub_") {
		t.Fatalf("expected khub_ prefix, got %q", a)
	}
	if len(a) != len("khub_")+26 {
		t.Fatalf("expected length %d, got %d (%q)", len("khub_")+26, len(a), a)
	}
}

func TestNewQueryID(t *testing.T) {
	id := NewQueryID()
	if !strings.HasPrefix(id, "kq_") {
		t.Fatalf("expected kq_ prefix, got %q", id)
	}
}
