package index

import (
	"database/sql"
	_ "embed"
	"fmt"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaSQL string

type Index struct {
	db  *sql.DB
	dim int
}

// Open creates or opens the index database at the given path.
// The dim parameter fixes the vector dimension.
func Open(path string, dim int) (*Index, error) {
	sqlite_vec.Auto() // registers vec0 extension on each new connection
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	stmt := fmt.Sprintf(schemaSQL, dim)
	if _, err := db.Exec(stmt); err != nil {
		db.Close()
		return nil, fmt.Errorf("index schema: %w", err)
	}
	return &Index{db: db, dim: dim}, nil
}

func (i *Index) Close() error { return i.db.Close() }

func (i *Index) Count() (int, error) {
	var n int
	if err := i.db.QueryRow("SELECT COUNT(*) FROM entries").Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}

func (i *Index) Dim() int { return i.dim }
