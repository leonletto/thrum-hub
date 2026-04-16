package schema

import (
	"strings"
	"testing"
	"time"
)

func TestEntryRoundTrip(t *testing.T) {
	original := &Entry{
		ID:         "khub_01HZTESTULID000000000000",
		Type:       TypeDecision,
		Title:      "Sonic over encoding/json for high-QPS handlers",
		Tags:       []string{"go", "json", "performance"},
		Language:   "go",
		SourceRepo: "falcondev/falcon-backend",
		Author:     "agent:implementer_api",
		Created:    time.Date(2026, 4, 15, 19, 30, 0, 0, time.UTC),
		Status:     StatusActive,
		Context:    "We were benchmarking a hot handler and saw json.Marshal dominating CPU.",
		Content:    "Switched to sonic. 3x throughput improvement.",
		WhyMatters: "Any high-QPS handler in this codebase should default to sonic.",
	}

	buf, err := original.Marshal()
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if !strings.HasPrefix(string(buf), "---\n") {
		t.Fatalf("expected frontmatter delimiter prefix, got %q", string(buf)[:10])
	}

	parsed, err := UnmarshalEntry(buf)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if parsed.ID != original.ID {
		t.Errorf("ID mismatch: got %q, want %q", parsed.ID, original.ID)
	}
	if parsed.Type != original.Type {
		t.Errorf("Type mismatch: got %q, want %q", parsed.Type, original.Type)
	}
	if parsed.Title != original.Title {
		t.Errorf("Title mismatch")
	}
	if len(parsed.Tags) != 3 || parsed.Tags[0] != "go" {
		t.Errorf("Tags mismatch: %v", parsed.Tags)
	}
	if parsed.Content != original.Content {
		t.Errorf("Content mismatch")
	}
}

func TestEntryValidate(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*Entry)
		wantErr string
	}{
		{"valid", func(e *Entry) {}, ""},
		{"missing id", func(e *Entry) { e.ID = "" }, "id: required"},
		{"bad type", func(e *Entry) { e.Type = "bogus" }, "type: invalid"},
		{"missing language", func(e *Entry) { e.Language = "" }, "language: required"},
		{"missing content", func(e *Entry) { e.Content = "" }, "content: required"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := &Entry{
				ID: "khub_x", Type: TypeDecision, Title: "t", Language: "go",
				Created: time.Now(), Status: StatusActive, Content: "body",
			}
			tc.mutate(e)
			err := e.Validate()
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("want no error, got %v", err)
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("want error containing %q, got %v", tc.wantErr, err)
				}
			}
		})
	}
}
