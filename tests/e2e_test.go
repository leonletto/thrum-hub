package tests

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leonletto/thrum-hub/internal/embed"
	"github.com/leonletto/thrum-hub/internal/index"
	"github.com/leonletto/thrum-hub/internal/query"
	"github.com/leonletto/thrum-hub/internal/schema"
	"github.com/leonletto/thrum-hub/internal/store"
	"github.com/leonletto/thrum-hub/internal/submit"
	"github.com/leonletto/thrum-hub/internal/telemetry"
)

func TestSubmitProcessSearchFeedback(t *testing.T) {
	root := t.TempDir()
	s, err := store.New(root)
	if err != nil {
		t.Fatal(err)
	}
	idx, err := index.Open(filepath.Join(root, "index.db"), 3)
	if err != nil {
		t.Fatal(err)
	}
	defer idx.Close()
	tl, err := telemetry.Open(filepath.Join(root, "telemetry", "queries.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer tl.Close()
	emb := embed.NewStub(3)

	// 1. Submit.
	id, err := submit.Do(s, submit.Request{
		Type:     schema.TypeGotcha,
		Language: "go",
		Title:    "context cancellation in long-running goroutines",
		Tags:     []string{"go", "context", "goroutine"},
		Author:   "agent:e2e-test",
		Body:     "Always wire a done channel into the goroutine so it exits when ctx.Done() fires.",
		Context:  "We hit a shutdown hang on a worker pool.",
	})
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	inbox, _ := s.InboxList()
	if len(inbox) != 1 || inbox[0] != id {
		t.Fatalf("inbox state wrong: %v", inbox)
	}

	// 2. Process inbox (Plan 1 debug path).
	for _, pid := range inbox {
		e, err := s.InboxRead(pid)
		if err != nil {
			t.Fatal(err)
		}
		if err := s.EntryWrite(e); err != nil {
			t.Fatal(err)
		}
		path := s.EntryPath(pid)
		if err := idx.Upsert(context.Background(), e, path, emb); err != nil {
			t.Fatal(err)
		}
		_ = s.InboxDelete(pid)
	}
	n, _ := idx.Count()
	if n != 1 {
		t.Fatalf("index count: got %d, want 1", n)
	}

	// 3. Search.
	res, err := query.Do(context.Background(), query.Deps{
		Store: s, Index: idx, Telemetry: tl, Embedder: emb,
	}, query.Request{
		Query: "context cancellation in long-running goroutines",
		K:     5, Agent: "agent:e2e-test", SourceRepo: "test/repo",
	})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(res.Hits) == 0 {
		t.Fatal("no hits")
	}
	if res.Hits[0].ID != id {
		t.Errorf("top hit: got %q, want %q", res.Hits[0].ID, id)
	}

	// 4. Feedback.
	if err := tl.RecordFeedback(telemetry.Feedback{
		QueryID: res.QueryID, EntryID: id, Signal: "up", Note: "helpful",
	}); err != nil {
		t.Fatalf("feedback: %v", err)
	}

	// 5. Verify show (EntryRead).
	got, err := s.EntryRead(id)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got.Content, "done channel") {
		t.Errorf("content not preserved")
	}
}
