package schema

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type SourceRef struct {
	Commit      string `yaml:"commit,omitempty"`
	PR          string `yaml:"pr,omitempty"`
	Bead        string `yaml:"bead,omitempty"`
	ThrumThread string `yaml:"thrum_thread,omitempty"`
}

type Entry struct {
	ID         string      `yaml:"id"`
	Type       EntryType   `yaml:"type"`
	Title      string      `yaml:"title"`
	Tags       []string    `yaml:"tags,omitempty"`
	Language   string      `yaml:"language"`
	SourceRepo string      `yaml:"source_repo,omitempty"`
	SourceRefs []SourceRef `yaml:"source_refs,omitempty"`
	Author     string      `yaml:"author,omitempty"`
	Created    time.Time   `yaml:"created"`
	Supersedes []string    `yaml:"supersedes,omitempty"`
	Status     Status      `yaml:"status"`

	// Body sections (not in frontmatter)
	Context    string `yaml:"-"`
	Content    string `yaml:"-"`
	WhyMatters string `yaml:"-"`
}

const delim = "---\n"

// Marshal serializes the entry to its on-disk representation:
// YAML frontmatter delimited by `---`, followed by three markdown
// body sections ("## Context", "## Content", "## Why this matters").
func (e *Entry) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(delim)

	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(e); err != nil {
		return nil, fmt.Errorf("yaml encode: %w", err)
	}
	_ = enc.Close()

	buf.WriteString(delim)
	buf.WriteString("\n")

	if e.Context != "" {
		buf.WriteString("## Context\n")
		buf.WriteString(strings.TrimSpace(e.Context))
		buf.WriteString("\n\n")
	}
	if e.Content != "" {
		buf.WriteString("## Content\n")
		buf.WriteString(strings.TrimSpace(e.Content))
		buf.WriteString("\n\n")
	}
	if e.WhyMatters != "" {
		buf.WriteString("## Why this matters\n")
		buf.WriteString(strings.TrimSpace(e.WhyMatters))
		buf.WriteString("\n")
	}
	return buf.Bytes(), nil
}

// UnmarshalEntry parses an on-disk entry file.
func UnmarshalEntry(raw []byte) (*Entry, error) {
	s := string(raw)
	if !strings.HasPrefix(s, delim) {
		return nil, fmt.Errorf("missing frontmatter delimiter")
	}
	rest := s[len(delim):]
	end := strings.Index(rest, delim)
	if end < 0 {
		return nil, fmt.Errorf("missing closing frontmatter delimiter")
	}
	frontmatter := rest[:end]
	body := rest[end+len(delim):]

	var e Entry
	if err := yaml.Unmarshal([]byte(frontmatter), &e); err != nil {
		return nil, fmt.Errorf("yaml decode: %w", err)
	}
	e.Context, e.Content, e.WhyMatters = parseBody(body)
	return &e, nil
}

func parseBody(body string) (context, content, why string) {
	sections := map[string]*string{
		"## Context":          &context,
		"## Content":          &content,
		"## Why this matters": &why,
	}
	lines := strings.Split(body, "\n")
	var current *string
	var buf strings.Builder
	flush := func() {
		if current != nil {
			*current = strings.TrimSpace(buf.String())
			buf.Reset()
		}
	}
	for _, line := range lines {
		if p, ok := sections[strings.TrimSpace(line)]; ok {
			flush()
			current = p
			continue
		}
		if current != nil {
			buf.WriteString(line)
			buf.WriteString("\n")
		}
	}
	flush()
	return
}
