package ids

import (
	"crypto/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	entropyMu sync.Mutex
	entropy   = ulid.Monotonic(rand.Reader, 0)
)

func newULID() string {
	entropyMu.Lock()
	defer entropyMu.Unlock()
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}

// NewEntryID returns a "khub_<ULID>" identifier for a new entry.
func NewEntryID() string {
	return "khub_" + newULID()
}

// NewQueryID returns a "kq_<ULID>" identifier for a new query.
func NewQueryID() string {
	return "kq_" + newULID()
}
