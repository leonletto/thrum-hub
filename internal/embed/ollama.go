// Ollama embedding backend. The Ollama HTTP API is a simple POST to
// /api/embeddings. semantec_retrieval_engine has a similar client in its
// top-level embeddings/ package, but importing it would pull in chunking,
// dedup, and storage dependencies we don't need. This is a standalone
// reimplementation of just the embedding call.
package embed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
)

type Ollama struct {
	endpoint string
	model    string
	client   *http.Client
	dim      int32 // resolved lazily on first Embed
}

func NewOllama(endpoint, model string) *Ollama {
	return &Ollama{
		endpoint: endpoint,
		model:    model,
		client:   &http.Client{},
	}
}

type ollamaReq struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}
type ollamaResp struct {
	Embedding []float32 `json:"embedding"`
	Error     string    `json:"error,omitempty"`
}

func (o *Ollama) Embed(ctx context.Context, text string) ([]float32, error) {
	body, _ := json.Marshal(ollamaReq{Model: o.model, Prompt: text})
	url := o.endpoint + "/api/embeddings"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ollama embed: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama embed: http %d: %s", resp.StatusCode, raw)
	}
	var r ollamaResp
	if err := json.Unmarshal(raw, &r); err != nil {
		return nil, fmt.Errorf("ollama embed: decode: %w", err)
	}
	if r.Error != "" {
		return nil, fmt.Errorf("ollama embed: %s", r.Error)
	}
	if len(r.Embedding) == 0 {
		return nil, fmt.Errorf("ollama embed: empty vector")
	}
	atomic.CompareAndSwapInt32(&o.dim, 0, int32(len(r.Embedding)))
	return r.Embedding, nil
}

func (o *Ollama) Dim() int {
	return int(atomic.LoadInt32(&o.dim))
}
