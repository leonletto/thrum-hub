package embed

import (
	"context"
	"hash/fnv"
	"math"
)

// Embedder produces dense vectors for text. Implementations must be
// safe for concurrent use and must return vectors of a fixed dimension.
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
	Dim() int
}

// Stub is a deterministic hash-based pseudo-embedder for tests.
// It does NOT produce semantically meaningful vectors.
type Stub struct {
	dim int
}

func NewStub(dim int) *Stub { return &Stub{dim: dim} }

func (s *Stub) Dim() int { return s.dim }

func (s *Stub) Embed(_ context.Context, text string) ([]float32, error) {
	out := make([]float32, s.dim)
	for i := 0; i < s.dim; i++ {
		h := fnv.New64a()
		h.Write([]byte(text))
		h.Write([]byte{byte(i)})
		v := float32(h.Sum64()&0xFFFF) / 0xFFFF
		out[i] = v*2 - 1 // [-1, 1]
	}
	// Normalize.
	var norm float64
	for _, v := range out {
		norm += float64(v) * float64(v)
	}
	norm = math.Sqrt(norm)
	if norm > 0 {
		for i := range out {
			out[i] = float32(float64(out[i]) / norm)
		}
	}
	return out, nil
}
