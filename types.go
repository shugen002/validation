package validation

import (
	"fmt"
	"reflect"
)

// Validator interface defines the validation contract
type Validator interface {
	Validate() error
	ValidateWithBag(errorBag string) error
	Passes() bool
	Fails() bool
	Errors() *ErrorBag
	Valid() map[string]interface{}
	Invalid() map[string]interface{}
	AddRule(field string, rules ...Rule) Validator
	Sometimes(field string, rules []Rule, callback func(data map[string]interface{}) bool) Validator
	StopOnFirstFailure() Validator
}

// Rule interface defines a validation rule
type Rule interface {
	Passes(attribute string, value interface{}) bool
	Message() string
}

// DataAwareRule interface for rules that need access to all validation data
type DataAwareRule interface {
	Rule
	SetData(data map[string]interface{})
}

// ValidatorAwareRule interface for rules that need access to the validator
type ValidatorAwareRule interface {
	Rule
	SetValidator(validator Validator)
}

// ImplicitRule interface marks rules that imply the field is required
type ImplicitRule interface {
	Rule
	IsImplicit() bool
}

// ErrorBag holds validation error messages
type ErrorBag struct {
	messages map[string][]string
}

// NewErrorBag creates a new error bag
func NewErrorBag() *ErrorBag {
	return &ErrorBag{
		messages: make(map[string][]string),
	}
}

// Add adds an error message for a field
func (e *ErrorBag) Add(field, message string) {
	if e.messages == nil {
		e.messages = make(map[string][]string)
	}
	e.messages[field] = append(e.messages[field], message)
}

// Has checks if a field has any errors
func (e *ErrorBag) Has(field string) bool {
	_, exists := e.messages[field]
	return exists
}

// Get returns all error messages for a field
func (e *ErrorBag) Get(field string) []string {
	return e.messages[field]
}

// First returns the first error message for a field
func (e *ErrorBag) First(field string) string {
	if messages, exists := e.messages[field]; exists && len(messages) > 0 {
		return messages[0]
	}
	return ""
}

// All returns all error messages
func (e *ErrorBag) All() map[string][]string {
	return e.messages
}

// IsEmpty checks if the error bag is empty
func (e *ErrorBag) IsEmpty() bool {
	return len(e.messages) == 0
}

// IsNotEmpty checks if the error bag has any errors
func (e *ErrorBag) IsNotEmpty() bool {
	return !e.IsEmpty()
}

// Count returns the total number of error messages
func (e *ErrorBag) Count() int {
	count := 0
	for _, messages := range e.messages {
		count += len(messages)
	}
	return count
}

// ValidationException represents a validation failure
type ValidationException struct {
	Message string
	Errors  *ErrorBag
}

func (e *ValidationException) Error() string {
	return e.Message
}

// NewValidationException creates a new validation exception
func NewValidationException(message string, errors *ErrorBag) *ValidationException {
	return &ValidationException{
		Message: message,
		Errors:  errors,
	}
}

// Helper functions for type checking and conversion

// IsNil checks if a value is nil
func IsNil(value interface{}) bool {
	if value == nil {
		return true
	}
	
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

// ToString converts a value to string
func ToString(value interface{}) string {
	if value == nil {
		return ""
	}
	
	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ToFloat64 converts a value to float64
func ToFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}

// ToInt64 converts a value to int64
func ToInt64(value interface{}) (int64, bool) {
	switch v := value.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case int32:
		return int64(v), true
	case uint:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		return int64(v), true
	case float64:
		if v == float64(int64(v)) {
			return int64(v), true
		}
	case float32:
		if v == float32(int64(v)) {
			return int64(v), true
		}
	}
	return 0, false
}

// GetSize returns the size of a value (length for strings/slices, value for numbers)
func GetSize(value interface{}) (float64, bool) {
	if value == nil {
		return 0, false
	}
	
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return float64(len([]rune(v.String()))), true
	case reflect.Slice, reflect.Array, reflect.Map:
		return float64(v.Len()), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint()), true
	case reflect.Float32, reflect.Float64:
		return v.Float(), true
	}
	
	return 0, false
}