package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/leonletto/thrum-hub/internal/schema"
	"github.com/leonletto/thrum-hub/internal/submit"
	"github.com/spf13/cobra"
)

func cmdSubmit() *cobra.Command {
	var (
		typeStr    string
		language   string
		title      string
		tagsStr    string
		body       string
		bodyFile   string
		contextS   string
		whyMatters string
		supersedes []string
		sourceRepo string
		sourceRefs []string
	)
	c := &cobra.Command{
		Use:   "submit",
		Short: "Submit a new entry to the inbox",
		RunE: func(cmd *cobra.Command, args []string) error {
			contentBody := body
			switch {
			case bodyFile != "":
				raw, err := os.ReadFile(bodyFile)
				if err != nil {
					return fmt.Errorf("read --body-file: %w", err)
				}
				contentBody = string(raw)
			case body == "":
				raw, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("read stdin: %w", err)
				}
				contentBody = strings.TrimSpace(string(raw))
			}
			if contentBody == "" {
				return fmt.Errorf("body required (via --body, --body-file, or stdin)")
			}

			req := submit.Request{
				Type:       schema.EntryType(typeStr),
				Language:   language,
				Title:      title,
				Tags:       splitCSV(tagsStr),
				SourceRepo: sourceRepo,
				Author:     detectAuthor(),
				Body:       contentBody,
				Context:    contextS,
				WhyMatters: whyMatters,
				Supersedes: supersedes,
				SourceRefs: parseSourceRefs(sourceRefs),
			}
			id, err := submit.Do(a.store, req)
			if err != nil {
				return err
			}
			fmt.Println(id)
			return nil
		},
	}
	c.Flags().StringVar(&typeStr, "type", "", "entry type (pattern|decision|gotcha|runbook|tooling)")
	c.Flags().StringVar(&language, "language", "", "primary language")
	c.Flags().StringVar(&title, "title", "", "entry title")
	c.Flags().StringVar(&tagsStr, "tags", "", "comma-separated tags")
	c.Flags().StringVar(&body, "body", "", "body text (or use --body-file or stdin)")
	c.Flags().StringVar(&bodyFile, "body-file", "", "path to body content")
	c.Flags().StringVar(&contextS, "context", "", "`## Context` section")
	c.Flags().StringVar(&whyMatters, "why-matters", "", "`## Why this matters` section")
	c.Flags().StringSliceVar(&supersedes, "supersedes", nil, "IDs this entry supersedes (repeatable)")
	c.Flags().StringVar(&sourceRepo, "source-repo", "", "source repo (defaults to git origin)")
	c.Flags().StringSliceVar(&sourceRefs, "source-ref", nil, "source refs as key=value (repeatable)")
	_ = c.MarkFlagRequired("type")
	_ = c.MarkFlagRequired("language")
	_ = c.MarkFlagRequired("title")
	return c
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func parseSourceRefs(refs []string) []schema.SourceRef {
	var out []schema.SourceRef
	for _, r := range refs {
		eq := strings.IndexByte(r, '=')
		if eq < 0 {
			continue
		}
		k, v := r[:eq], r[eq+1:]
		sr := schema.SourceRef{}
		switch k {
		case "commit":
			sr.Commit = v
		case "pr":
			sr.PR = v
		case "bead":
			sr.Bead = v
		case "thrum_thread":
			sr.ThrumThread = v
		}
		out = append(out, sr)
	}
	return out
}

func detectAuthor() string {
	if n := os.Getenv("THRUM_AGENT_NAME"); n != "" {
		return "agent:" + n
	}
	if n := os.Getenv("USER"); n != "" {
		return "human:" + n
	}
	return "unknown"
}
