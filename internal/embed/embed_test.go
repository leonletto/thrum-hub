package embed

import (
	"context"
	"testing"
)

func TestStubEmbedder(t *testing.T) {
	e := NewStub(3)
	vec, err := e.Embed(context.Background(), "hello")
	if err != nil {
		t.Fatal(err)
	}
	if len(vec) != 3 {
		t.Fatalf("expected dim 3, got %d", len(vec))
	}

	vec2, _ := e.Embed(context.Background(), "hello")
	if vec[0] != vec2[0] {
		t.Errorf("stub should be deterministic for same input")
	}

	vec3, _ := e.Embed(context.Background(), "different")
	if vec[0] == vec3[0] && vec[1] == vec3[1] && vec[2] == vec3[2] {
		t.Errorf("stub should produce different vectors for different inputs")
	}

	if e.Dim() != 3 {
		t.Errorf("Dim() should be 3")
	}
}
