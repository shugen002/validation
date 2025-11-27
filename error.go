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
	return e.Reason
}

type ErrValidationFailed struct {
	Reason string
}

func (e *ErrValidationFailed) Error() string {
	return e.Reason
}

func newErrValidationFailed(reason string) error {
	return &ErrValidationFailed{Reason: reason}
}

func IsErrValidationFailed(err error) bool {
	_, ok := err.(*ErrValidationFailed)
	return ok
}

func IsErrParsingRules(err error) bool {
	_, ok := err.(*ErrParsingRules)
	return ok
}
