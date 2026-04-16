package embed

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIEmbedder_OpenAIStyle(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("missing Authorization header")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"embedding": []float32{0.5, 0.6, 0.7}},
			},
		})
	}))
	defer srv.Close()

	e := NewAPI(srv.URL, "text-embedding-3-small", "test-key")
	vec, err := e.Embed(context.Background(), "ping")
	if err != nil {
		t.Fatalf("Embed: %v", err)
	}
	if len(vec) != 3 || vec[0] != 0.5 {
		t.Errorf("unexpected vector: %v", vec)
	}
}
