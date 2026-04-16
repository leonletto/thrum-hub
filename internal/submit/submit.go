package submit

import (
	"fmt"
	"time"

	"github.com/leonletto/thrum-hub/internal/ids"
	"github.com/leonletto/thrum-hub/internal/schema"
	"github.com/leonletto/thrum-hub/internal/store"
)

type Request struct {
	Type       schema.EntryType
	Language   string
	Title      string
	Tags       []string
	SourceRepo string
	SourceRefs []schema.SourceRef
	Author     string
	Supersedes []string

	// Body sections
	Context    string
	Body       string // becomes Content
	WhyMatters string
}

// Do validates the request, stamps identity + timestamp, and writes to the inbox.
// Returns the generated entry ID.
func Do(s *store.Store, r Request) (string, error) {
	e := &schema.Entry{
		ID:         ids.NewEntryID(),
		Type:       r.Type,
		Title:      r.Title,
		Tags:       r.Tags,
		Language:   r.Language,
		SourceRepo: r.SourceRepo,
		SourceRefs: r.SourceRefs,
		Author:     r.Author,
		Created:    time.Now().UTC(),
		Supersedes: r.Supersedes,
		Status:     schema.StatusActive,
		Context:    r.Context,
		Content:    r.Body,
		WhyMatters: r.WhyMatters,
	}
	if err := e.Validate(); err != nil {
		return "", fmt.Errorf("submit: %w", err)
	}
	if err := s.InboxWrite(e); err != nil {
		return "", fmt.Errorf("submit: write inbox: %w", err)
	}
	return e.ID, nil
}
