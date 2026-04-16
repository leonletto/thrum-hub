package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/leonletto/thrum-hub/internal/config"
	"github.com/leonletto/thrum-hub/internal/embed"
	"github.com/leonletto/thrum-hub/internal/index"
	"github.com/leonletto/thrum-hub/internal/store"
	"github.com/leonletto/thrum-hub/internal/telemetry"
	"github.com/spf13/cobra"
)

var version = "0.1.0-dev"

// Global state wired once per invocation via persistentPreRun.
type app struct {
	hubRoot string
	cfg     *config.Config
	store   *store.Store
	index   *index.Index
	telem   *telemetry.Telemetry
	embed   embed.Embedder
}

var a = &app{}

func main() {
	root := &cobra.Command{
		Use:   "thrumhub",
		Short: "Thrum Hub — shared coding-knowledge corpus for agents",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// `version` and `init` skip wiring.
			if cmd.Name() == "version" || cmd.Name() == "init" {
				return nil
			}
			return wireApp()
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return tearDownApp()
		},
	}
	root.PersistentFlags().StringVar(&a.hubRoot, "hub", "", "path to hub root (defaults to $THRUMHUB_ROOT or CWD)")

	root.AddCommand(cmdVersion())
	root.AddCommand(cmdInit())
	root.AddCommand(cmdStatus())
	root.AddCommand(cmdSubmit())
	root.AddCommand(cmdSearch())
	root.AddCommand(cmdShow())
	root.AddCommand(cmdRelated())
	root.AddCommand(cmdLs())
	root.AddCommand(cmdFeedback())
	root.AddCommand(cmdHistory())
	root.AddCommand(cmdDebug())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func cmdVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("thrumhub", version)
		},
	}
}

func wireApp() error {
	if a.hubRoot == "" {
		a.hubRoot = os.Getenv("THRUMHUB_ROOT")
	}
	if a.hubRoot == "" {
		cwd, _ := os.Getwd()
		a.hubRoot = cwd
	}
	cfgPath := filepath.Join(a.hubRoot, "config", "hub.yaml")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	a.cfg = cfg
	a.store, err = store.New(cfg.Root)
	if err != nil {
		return err
	}
	a.embed = buildEmbedder(cfg.Embedding)
	// Probe embedder to learn dimension.
	vec, err := a.embed.Embed(context.Background(), "init probe")
	if err != nil {
		return fmt.Errorf("embed probe: %w", err)
	}
	dim := len(vec)
	a.index, err = index.Open(filepath.Join(cfg.Root, "index.db"), dim)
	if err != nil {
		return err
	}
	a.telem, err = telemetry.Open(filepath.Join(cfg.Root, "telemetry", "queries.db"))
	if err != nil {
		return err
	}
	return nil
}

func buildEmbedder(cfg config.EmbeddingConfig) embed.Embedder {
	switch cfg.Backend {
	case "ollama":
		return embed.NewOllama(cfg.OllamaEndpoint, cfg.Model)
	case "api":
		key := os.Getenv(cfg.APIKeyEnv)
		return embed.NewAPI(cfg.APIBase, cfg.APIModel, key)
	}
	return embed.NewStub(3) // should never reach this; config validation rejected it
}

func tearDownApp() error {
	if a.index != nil {
		a.index.Close()
	}
	if a.telem != nil {
		a.telem.Close()
	}
	return nil
}
