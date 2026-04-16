CREATE TABLE IF NOT EXISTS entries (
  id          TEXT PRIMARY KEY,
  title       TEXT NOT NULL,
  type        TEXT NOT NULL,
  language    TEXT NOT NULL,
  tags        TEXT NOT NULL,   -- JSON array
  source_repo TEXT,
  status      TEXT NOT NULL,
  path        TEXT NOT NULL,   -- absolute path to source file
  updated_at  DATETIME NOT NULL
);

-- sqlite-vec virtual table for vectors. Dimension set at creation time.
-- Dimension placeholder is substituted when the code opens the db.
CREATE VIRTUAL TABLE IF NOT EXISTS vec_entries USING vec0(
  id TEXT PRIMARY KEY,
  embedding float[%d]
);
