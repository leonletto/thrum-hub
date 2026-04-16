package index

import (
	"context"
	"fmt"
	"strings"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	"github.com/leonletto/thrum-hub/internal/embed"
)

type SearchFilter struct {
	Language string
	Type     string
	Status   string // defaults to "active" if empty
}

type Hit struct {
	ID         string
	Title      string
	Type       string
	Language   string
	SourceRepo string
	Path       string
	Score      float64
}

func (i *Index) Search(ctx context.Context, query string, k int, emb embed.Embedder, f SearchFilter) ([]Hit, error) {
	qvec, err := emb.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("embed query: %w", err)
	}
	if len(qvec) != i.dim {
		return nil, fmt.Errorf("query dim mismatch: got %d, want %d", len(qvec), i.dim)
	}
	status := f.Status
	if status == "" {
		status = "active"
	}

	blob, err := sqlite_vec.SerializeFloat32(qvec)
	if err != nil {
		return nil, fmt.Errorf("serialize query: %w", err)
	}

	// sqlite-vec KNN: WHERE embedding MATCH ? AND k = ?
	// Join to the entries metadata table for filtering and enrichment.
	sb := strings.Builder{}
	sb.WriteString(`
SELECT e.id, e.title, e.type, e.language, e.source_repo, e.path, v.distance
FROM vec_entries v
JOIN entries e ON e.id = v.id
WHERE v.embedding MATCH ? AND k = ? AND e.status = ?
`)
	args := []any{blob, k, status}
	if f.Language != "" {
		sb.WriteString(" AND e.language = ?")
		args = append(args, f.Language)
	}
	if f.Type != "" {
		sb.WriteString(" AND e.type = ?")
		args = append(args, f.Type)
	}
	sb.WriteString(" ORDER BY v.distance ASC")

	rows, err := i.db.QueryContext(ctx, sb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Hit
	for rows.Next() {
		var h Hit
		if err := rows.Scan(&h.ID, &h.Title, &h.Type, &h.Language, &h.SourceRepo, &h.Path, &h.Score); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, rows.Err()
}
