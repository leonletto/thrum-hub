package telemetry

import (
	"path/filepath"
	"testing"
	"time"
)

func TestRecordQueryAndFeedback(t *testing.T) {
	dir := t.TempDir()
	tl, err := Open(filepath.Join(dir, "queries.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer tl.Close()

	q := QueryRecord{
		ID:         "kq_01TEST",
		Timestamp:  time.Now().UTC(),
		QueryText:  "goroutine shutdown",
		Agent:      "agent:test",
		SourceRepo: "some/repo",
		NumResults: 2,
		TopScore:   0.87,
		Results: []Result{
			{EntryID: "khub_01A", Rank: 0, Score: 0.87},
			{EntryID: "khub_01B", Rank: 1, Score: 0.55},
		},
	}
	if err := tl.RecordQuery(q); err != nil {
		t.Fatal(err)
	}
	if err := tl.RecordFeedback(Feedback{
		QueryID: "kq_01TEST", EntryID: "khub_01A", Signal: "up", Note: "helpful",
		Timestamp: time.Now().UTC(),
	}); err != nil {
		t.Fatal(err)
	}

	// Verify by raw select.
	var count int
	if err := tl.db.QueryRow("SELECT COUNT(*) FROM queries").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("queries count: got %d, want 1", count)
	}
	if err := tl.db.QueryRow("SELECT COUNT(*) FROM query_results").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("query_results count: got %d, want 2", count)
	}
	if err := tl.db.QueryRow("SELECT COUNT(*) FROM feedback").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("feedback count: got %d, want 1", count)
	}
}

func TestRecentQueries(t *testing.T) {
	dir := t.TempDir()
	tl, _ := Open(filepath.Join(dir, "queries.db"))
	defer tl.Close()

	for i := 0; i < 3; i++ {
		_ = tl.RecordQuery(QueryRecord{
			ID: "kq_" + string(rune('A'+i)), Timestamp: time.Now().UTC(), QueryText: "q",
			NumResults: 1,
		})
	}
	recent, err := tl.RecentQueries("", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(recent) != 3 {
		t.Fatalf("got %d, want 3", len(recent))
	}
}
