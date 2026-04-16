package submit

import (
	"testing"
	"time"

	"github.com/leonletto/thrum-hub/internal/schema"
	"github.com/leonletto/thrum-hub/internal/store"
)

func TestSubmit(t *testing.T) {
	root := t.TempDir()
	s, _ := store.New(root)

	req := Request{
		Type:       schema.TypeGotcha,
		Language:   "go",
		Title:      "ctx cancel gotcha",
		Tags:       []string{"context", "goroutine"},
		SourceRepo: "falcondev/falcon-backend",
		Author:     "agent:tester",
		Body:       "body",
		Context:    "ctx",
		WhyMatters: "why",
	}
	id, err := Do(s, req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if id == "" {
		t.Fatal("empty id")
	}
	e, err := s.InboxRead(id)
	if err != nil {
		t.Fatalf("InboxRead: %v", err)
	}
	if e.Title != req.Title {
		t.Errorf("title: got %q", e.Title)
	}
	if e.Created.IsZero() {
		t.Errorf("Created not set")
	}
	if time.Since(e.Created) > time.Minute {
		t.Errorf("Created too old: %v", e.Created)
	}
	if e.Status != schema.StatusActive {
		t.Errorf("Status: got %q", e.Status)
	}
}

func TestSubmitRejectsBadType(t *testing.T) {
	root := t.TempDir()
	s, _ := store.New(root)
	_, err := Do(s, Request{
		Type: "bogus", Language: "go", Title: "t", Body: "b",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}
