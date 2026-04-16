package index

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	"github.com/leonletto/thrum-hub/internal/embed"
	"github.com/leonletto/thrum-hub/internal/schema"
)

// Upsert embeds and stores an entry. Replaces any existing row with the same ID.
func (i *Index) Upsert(ctx context.Context, e *schema.Entry, path string, emb embed.Embedder) error {
	text := entryEmbedText(e)
	vec, err := emb.Embed(ctx, text)
	if err != nil {
		return fmt.Errorf("embed: %w", err)
	}
	if len(vec) != i.dim {
		return fmt.Errorf("embed: dim mismatch: got %d, want %d", len(vec), i.dim)
	}

	tagsJSON, _ := json.Marshal(e.Tags)

	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Replace metadata row.
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO entries (id, title, type, language, tags, source_repo, status, path, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
		  title=excluded.title, type=excluded.type, language=excluded.language,
		  tags=excluded.tags, source_repo=excluded.source_repo, status=excluded.status,
		  path=excluded.path, updated_at=excluded.updated_at
	`, e.ID, e.Title, string(e.Type), e.Language, string(tagsJSON),
		e.SourceRepo, string(e.Status), path, time.Now().UTC()); err != nil {
		return err
	}

	// Replace vector row using sqlite-vec's serialization.
	blob, err := sqlite_vec.SerializeFloat32(vec)
	if err != nil {
		return fmt.Errorf("serialize vector: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM vec_entries WHERE id = ?`, e.ID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO vec_entries(id, embedding) VALUES (?, ?)`, e.ID, blob); err != nil {
		return err
	}
	return tx.Commit()
}

// entryEmbedText is the canonical text to embed for an entry.
// Keep this stable across versions of the index; re-indexing is
// necessary if it changes.
func entryEmbedText(e *schema.Entry) string {
	return e.Title + "\n\n" + e.Context + "\n\n" + e.Content + "\n\n" + e.WhyMatters
}
