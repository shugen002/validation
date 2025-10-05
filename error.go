package validation

import "fmt"

type ErrUnknownRule struct {
	Rule string
}

func (e *ErrUnknownRule) Error() string {
	return fmt.Sprintf("unknown validation rule: %s", e.Rule)
}

type ErrParsingRules struct {
	Reason string
}

func (e *ErrParsingRules) Error() string {
	return fmt.Sprintf("error parsing validation rules: %s", e.Reason)
}
