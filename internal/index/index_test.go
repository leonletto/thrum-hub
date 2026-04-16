package index

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/leonletto/thrum-hub/internal/embed"
	"github.com/leonletto/thrum-hub/internal/schema"
)

func TestIndexUpsertAndCount(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "index.db")
	idx, err := Open(dbPath, 3)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer idx.Close()

	e := &schema.Entry{
		ID:       "khub_01HZIDX0TEST00000000000000",
		Title:    "test",
		Type:     schema.TypeDecision,
		Language: "go",
		Status:   schema.StatusActive,
		Created:  time.Now(),
		Content:  "payload",
	}
	emb := embed.NewStub(3)
	if err := idx.Upsert(context.Background(), e, "/tmp/fake.md", emb); err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	n, err := idx.Count()
	if err != nil {
		t.Fatalf("Count: %v", err)
	}
	if n != 1 {
		t.Errorf("Count: got %d, want 1", n)
	}
	// Re-upsert should be idempotent.
	if err := idx.Upsert(context.Background(), e, "/tmp/fake.md", emb); err != nil {
		t.Fatalf("Upsert 2: %v", err)
	}
	n, _ = idx.Count()
	if n != 1 {
		t.Errorf("Count after upsert 2: got %d, want 1", n)
	}
}
