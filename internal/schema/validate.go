package schema

import "fmt"

type ValidationError struct {
	Field  string
	Reason string
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("schema: %s: %s", v.Field, v.Reason)
}

func (e *Entry) Validate() error {
	if e.ID == "" {
		return &ValidationError{Field: "id", Reason: "required"}
	}
	if !e.Type.Valid() {
		return &ValidationError{Field: "type", Reason: fmt.Sprintf("invalid or empty: %q", e.Type)}
	}
	if e.Title == "" {
		return &ValidationError{Field: "title", Reason: "required"}
	}
	if e.Language == "" {
		return &ValidationError{Field: "language", Reason: "required"}
	}
	if !e.Status.Valid() {
		return &ValidationError{Field: "status", Reason: fmt.Sprintf("invalid or empty: %q", e.Status)}
	}
	if e.Content == "" {
		return &ValidationError{Field: "content", Reason: "required (body section `## Content`)"}
	}
	return nil
}
