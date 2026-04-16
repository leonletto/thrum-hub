package schema

type EntryType string

const (
	TypePattern  EntryType = "pattern"
	TypeDecision EntryType = "decision"
	TypeGotcha   EntryType = "gotcha"
	TypeRunbook  EntryType = "runbook"
	TypeTooling  EntryType = "tooling"
)

func (t EntryType) Valid() bool {
	switch t {
	case TypePattern, TypeDecision, TypeGotcha, TypeRunbook, TypeTooling:
		return true
	}
	return false
}

type Status string

const (
	StatusActive     Status = "active"
	StatusSuperseded Status = "superseded"
	StatusDeprecated Status = "deprecated"
)

func (s Status) Valid() bool {
	switch s {
	case StatusActive, StatusSuperseded, StatusDeprecated:
		return true
	}
	return false
}
