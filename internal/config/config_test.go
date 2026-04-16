package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "hub.yaml")
	yaml := `
root: /tmp/hub
embedding:
  backend: ollama
  model: nomic-embed-text
`
	if err := os.WriteFile(cfgPath, []byte(yaml), 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Root != "/tmp/hub" {
		t.Errorf("Root: got %q, want /tmp/hub", cfg.Root)
	}
	if cfg.Embedding.Backend != "ollama" {
		t.Errorf("Backend: got %q", cfg.Embedding.Backend)
	}
	if cfg.Embedding.Model != "nomic-embed-text" {
		t.Errorf("Model: got %q", cfg.Embedding.Model)
	}
	// Defaults applied:
	if cfg.Embedding.OllamaEndpoint == "" {
		t.Errorf("expected default OllamaEndpoint")
	}
}

func TestLoadValidatesBackend(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "hub.yaml")
	yaml := `
root: /tmp/hub
embedding:
  backend: bogus
  model: nomic-embed-text
`
	if err := os.WriteFile(cfgPath, []byte(yaml), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(cfgPath)
	if err == nil {
		t.Fatalf("expected error for bogus backend, got nil")
	}
}
