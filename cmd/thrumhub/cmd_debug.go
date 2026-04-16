package main

import (
	"context"
	"fmt"

	"github.com/leonletto/thrum-hub/internal/embed"
	"github.com/leonletto/thrum-hub/internal/index"
	"github.com/leonletto/thrum-hub/internal/store"
)

// DebugProcessInbox moves all pending inbox entries to entries/ and indexes
// them with no judgment applied. This is a Plan 1 test affordance; Plan 2
// will replace it with real admin-agent curation.
func DebugProcessInbox(ctx context.Context, s *store.Store, idx *index.Index, emb embed.Embedder) (processed int, err error) {
	ids, err := s.InboxList()
	if err != nil {
		return 0, err
	}
	for _, id := range ids {
		e, err := s.InboxRead(id)
		if err != nil {
			return processed, fmt.Errorf("read %s: %w", id, err)
		}
		if err := s.EntryWrite(e); err != nil {
			return processed, fmt.Errorf("promote %s: %w", id, err)
		}
		path := s.EntryPath(id)
		if err := idx.Upsert(ctx, e, path, emb); err != nil {
			return processed, fmt.Errorf("index %s: %w", id, err)
		}
		if err := s.InboxDelete(id); err != nil {
			return processed, fmt.Errorf("delete inbox %s: %w", id, err)
		}
		processed++
	}
	return processed, nil
}
