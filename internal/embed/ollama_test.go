package embed

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestOllamaEmbedder_Embed(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/embeddings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"embedding": []float32{0.1, 0.2, 0.3, 0.4},
		})
	}))
	defer srv.Close()

	e := NewOllama(srv.URL, "nomic-embed-text")
	vec, err := e.Embed(context.Background(), "hello world")
	if err != nil {
		t.Fatalf("Embed: %v", err)
	}
	if len(vec) != 4 || vec[0] != 0.1 {
		t.Errorf("unexpected vector: %v", vec)
	}
	if gotBody["model"] != "nomic-embed-text" {
		t.Errorf("unexpected model: %v", gotBody["model"])
	}
	if gotBody["prompt"] != "hello world" {
		t.Errorf("unexpected prompt: %v", gotBody["prompt"])
	}
}

func TestOllamaEmbedder_Live(t *testing.T) {
	if os.Getenv("THRUMHUB_LIVE_OLLAMA") == "" {
		t.Skip("set THRUMHUB_LIVE_OLLAMA=1 to run live test against http://localhost:11434")
	}
	e := NewOllama("http://localhost:11434", "nomic-embed-text")
	vec, err := e.Embed(context.Background(), "test context cancellation in goroutines")
	if err != nil {
		t.Fatalf("live Embed: %v", err)
	}
	if len(vec) < 64 {
		t.Errorf("unexpectedly small vector: %d", len(vec))
	}
}
