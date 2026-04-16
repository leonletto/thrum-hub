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

// API is an OpenAI-compatible embedding client. Works against any
// provider whose embeddings endpoint matches the OpenAI schema
// (OpenAI, Together, Fireworks, local OpenAI-compatible proxies).
type API struct {
	endpoint string
	model    string
	apiKey   string
	client   *http.Client
	dim      int32
}

func NewAPI(endpoint, model, apiKey string) *API {
	return &API{
		endpoint: endpoint,
		model:    model,
		apiKey:   apiKey,
		client:   &http.Client{},
	}
}

type apiReq struct {
	Model string `json:"model"`
	Input string `json:"input"`
}
type apiResp struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (a *API) Embed(ctx context.Context, text string) ([]float32, error) {
	body, _ := json.Marshal(apiReq{Model: a.model, Input: text})
	url := a.endpoint + "/embeddings"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("api embed: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api embed: http %d: %s", resp.StatusCode, raw)
	}
	var r apiResp
	if err := json.Unmarshal(raw, &r); err != nil {
		return nil, fmt.Errorf("api embed: decode: %w", err)
	}
	if r.Error.Message != "" {
		return nil, fmt.Errorf("api embed: %s", r.Error.Message)
	}
	if len(r.Data) == 0 || len(r.Data[0].Embedding) == 0 {
		return nil, fmt.Errorf("api embed: empty response")
	}
	atomic.CompareAndSwapInt32(&a.dim, 0, int32(len(r.Data[0].Embedding)))
	return r.Data[0].Embedding, nil
}

func (a *API) Dim() int {
	return int(atomic.LoadInt32(&a.dim))
}
