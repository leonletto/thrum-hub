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

func TestIndexSearch(t *testing.T) {
	dir := t.TempDir()
	idx, _ := Open(filepath.Join(dir, "index.db"), 3)
	defer idx.Close()

	emb := embed.NewStub(3)
	entries := []*schema.Entry{
		{ID: "khub_01AAAAAAAAAAAAAAAAAAAAAAAA", Title: "goroutine context cancellation", Type: schema.TypeGotcha, Language: "go", Status: schema.StatusActive, Content: "shutdown pattern"},
		{ID: "khub_01BBBBBBBBBBBBBBBBBBBBBBBB", Title: "sonic json", Type: schema.TypeDecision, Language: "go", Status: schema.StatusActive, Content: "encoding perf"},
		{ID: "khub_01CCCCCCCCCCCCCCCCCCCCCCCC", Title: "sql migration", Type: schema.TypePattern, Language: "sql", Status: schema.StatusActive, Content: "not go"},
	}
	for _, e := range entries {
		if err := idx.Upsert(context.Background(), e, "/tmp/"+e.ID+".md", emb); err != nil {
			t.Fatal(err)
		}
	}

	// Query using one of the indexed titles — should surface some ranking.
	results, err := idx.Search(context.Background(), "goroutine context cancellation", 3, emb, SearchFilter{})
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("expected non-empty results")
	}
	if results[0].ID != "khub_01AAAAAAAAAAAAAAAAAAAAAAAA" {
		t.Errorf("expected top result to match query title, got %q", results[0].ID)
	}

	// Language filter excludes SQL.
	results, _ = idx.Search(context.Background(), "anything", 10, emb, SearchFilter{Language: "go"})
	for _, r := range results {
		if r.Language != "go" {
			t.Errorf("filter violation: got language %q", r.Language)
		}
	}
}
