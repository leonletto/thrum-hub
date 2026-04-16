package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func cmdInit() *cobra.Command {
	return &cobra.Command{
		Use:   "init <path>",
		Short: "Scaffold a hub at the given path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}
			for _, d := range []string{
				root,
				filepath.Join(root, "config"),
				filepath.Join(root, "inbox"),
				filepath.Join(root, "entries"),
				filepath.Join(root, "archive", "superseded"),
				filepath.Join(root, "telemetry"),
			} {
				if err := os.MkdirAll(d, 0755); err != nil {
					return err
				}
			}
			cfgPath := filepath.Join(root, "config", "hub.yaml")
			if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
				content := fmt.Sprintf(`root: %s

embedding:
  backend: ollama
  model: nomic-embed-text
  ollama_endpoint: http://localhost:11434
`, root)
				if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
					return err
				}
			}
			fmt.Printf("initialized hub at %s\n", root)
			fmt.Printf("config: %s\n", cfgPath)
			return nil
		},
	}
}
