package query

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/leonletto/thrum-hub/internal/embed"
	"github.com/leonletto/thrum-hub/internal/index"
	"github.com/leonletto/thrum-hub/internal/schema"
	"github.com/leonletto/thrum-hub/internal/store"
	"github.com/leonletto/thrum-hub/internal/telemetry"
)

func TestDoRecordsTelemetry(t *testing.T) {
	root := t.TempDir()
	s, _ := store.New(root)
	idx, _ := index.Open(filepath.Join(root, "index.db"), 3)
	defer idx.Close()
	tl, _ := telemetry.Open(filepath.Join(root, "queries.db"))
	defer tl.Close()
	emb := embed.NewStub(3)

	e := &schema.Entry{
		ID: "khub_01QTEST00000000000000000000",
		Title: "search target", Type: schema.TypeDecision, Language: "go",
		Status: schema.StatusActive, Created: time.Now(), Content: "body",
	}
	_ = s.EntryWrite(e)
	_ = idx.Upsert(context.Background(), e, s.EntryPath(e.ID), emb)

	result, err := Do(context.Background(), Deps{
		Store: s, Index: idx, Telemetry: tl, Embedder: emb,
	}, Request{
		Query: "search target", K: 5, Agent: "agent:tester", SourceRepo: "test/repo",
	})
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if result.QueryID == "" {
		t.Fatal("empty QueryID")
	}
	if len(result.Hits) == 0 {
		t.Fatal("no hits")
	}
	// Verify telemetry got recorded.
	recent, _ := tl.RecentQueries("", 10)
	if len(recent) != 1 {
		t.Errorf("telemetry rows: got %d, want 1", len(recent))
	}
}
