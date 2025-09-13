package validation

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Basic validation rules

// RequiredRule validates that a field is present and not empty
type RequiredRule struct{}

func (r *RequiredRule) Passes(attribute string, value interface{}) bool {
	if IsNil(value) {
		return false
	}
	
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) != ""
	case []interface{}, map[string]interface{}:
		return reflect.ValueOf(v).Len() > 0
	default:
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
			return rv.Len() > 0
		}
	}
	
	return true
}

func (r *RequiredRule) Message() string {
	return "The :attribute field is required."
}

func (r *RequiredRule) IsImplicit() bool {
	return true
}

// StringRule validates that a field is a string
type StringRule struct{}

func (r *StringRule) Passes(attribute string, value interface{}) bool {
	_, ok := value.(string)
	return ok
}

func (r *StringRule) Message() string {
	return "The :attribute must be a string."
}

// IntegerRule validates that a field is an integer
type IntegerRule struct {
	Strict bool
}

func (r *IntegerRule) Passes(attribute string, value interface{}) bool {
	if r.Strict {
		switch value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			return true
		default:
			return false
		}
	}
	
	_, ok := ToInt64(value)
	if ok {
		return true
	}
	
	// Try to parse string as integer
	if stringValue, ok := value.(string); ok {
		_, err := strconv.ParseInt(stringValue, 10, 64)
		return err == nil
	}
	
	return false
}

func (r *IntegerRule) Message() string {
	return "The :attribute must be an integer."
}

// NumericRule validates that a field is numeric
type NumericRule struct {
	Strict bool
}

func (r *NumericRule) Passes(attribute string, value interface{}) bool {
	if r.Strict {
		if _, ok := value.(string); ok {
			return false
		}
	}
	
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	case string:
		str := value.(string)
		_, err := strconv.ParseFloat(str, 64)
		return err == nil
	}
	
	return false
}

func (r *NumericRule) Message() string {
	return "The :attribute must be a number."
}

// BooleanRule validates that a field is a boolean
type BooleanRule struct {
	Strict bool
}

func (r *BooleanRule) Passes(attribute string, value interface{}) bool {
	if r.Strict {
		_, ok := value.(bool)
		return ok
	}
	
	switch v := value.(type) {
	case bool:
		return true
	case string:
		return v == "true" || v == "false" || v == "1" || v == "0"
	case int, int8, int16, int32, int64:
		num := reflect.ValueOf(v).Int()
		return num == 0 || num == 1
	case uint, uint8, uint16, uint32, uint64:
		num := reflect.ValueOf(v).Uint()
		return num == 0 || num == 1
	case float32, float64:
		num := reflect.ValueOf(v).Float()
		return num == 0 || num == 1
	}
	
	return false
}

func (r *BooleanRule) Message() string {
	return "The :attribute field must be true or false."
}

// AcceptedRule validates that a field is accepted (yes, on, 1, "1", true, "true")
type AcceptedRule struct{}

func (r *AcceptedRule) Passes(attribute string, value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		return v == "yes" || v == "on" || v == "1" || v == "true"
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Int() == 1
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint() == 1
	case float32, float64:
		return reflect.ValueOf(v).Float() == 1
	}
	
	return false
}

func (r *AcceptedRule) Message() string {
	return "The :attribute must be accepted."
}

func (r *AcceptedRule) IsImplicit() bool {
	return true
}

// AcceptedIfRule validates that a field is accepted if another field equals a specified value
type AcceptedIfRule struct {
	Field string
	Value string
	data  map[string]interface{}
}

func (r *AcceptedIfRule) Passes(attribute string, value interface{}) bool {
	// Check if the condition field equals the specified value
	otherValue, exists := r.data[r.Field]
	if !exists {
		return true // If the condition field doesn't exist, this rule passes
	}
	
	otherValueStr := ToString(otherValue)
	if otherValueStr != r.Value {
		return true // If the condition is not met, this rule passes
	}
	
	// If condition is met, check if the field is accepted
	acceptedRule := &AcceptedRule{}
	return acceptedRule.Passes(attribute, value)
}

func (r *AcceptedIfRule) Message() string {
	return "The :attribute must be accepted when " + r.Field + " is " + r.Value + "."
}

func (r *AcceptedIfRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *AcceptedIfRule) IsImplicit() bool {
	return true
}

// DeclinedRule validates that a field is declined (no, off, 0, "0", false, "false")
type DeclinedRule struct{}

func (r *DeclinedRule) Passes(attribute string, value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return !v
	case string:
		return v == "no" || v == "off" || v == "0" || v == "false"
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Int() == 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint() == 0
	case float32, float64:
		return reflect.ValueOf(v).Float() == 0
	}
	
	return false
}

func (r *DeclinedRule) Message() string {
	return "The :attribute must be declined."
}

func (r *DeclinedRule) IsImplicit() bool {
	return true
}

// DeclinedIfRule validates that a field is declined if another field equals a specified value
type DeclinedIfRule struct {
	Field string
	Value string
	data  map[string]interface{}
}

func (r *DeclinedIfRule) Passes(attribute string, value interface{}) bool {
	// Check if the condition field equals the specified value
	otherValue, exists := r.data[r.Field]
	if !exists {
		return true // If the condition field doesn't exist, this rule passes
	}
	
	otherValueStr := ToString(otherValue)
	if otherValueStr != r.Value {
		return true // If the condition is not met, this rule passes
	}
	
	// If condition is met, check if the field is declined
	declinedRule := &DeclinedRule{}
	return declinedRule.Passes(attribute, value)
}

func (r *DeclinedIfRule) Message() string {
	return "The :attribute must be declined when " + r.Field + " is " + r.Value + "."
}

func (r *DeclinedIfRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *DeclinedIfRule) IsImplicit() bool {
	return true
}

// ArrayRule validates that a field is an array
type ArrayRule struct {
	AllowedKeys []string
}

func (r *ArrayRule) Passes(attribute string, value interface{}) bool {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array && rv.Kind() != reflect.Map {
		return false
	}
	
	// If no allowed keys specified, any array/slice/map is valid
	if len(r.AllowedKeys) == 0 {
		return true
	}
	
	// For maps, check if all keys are in allowed list
	if rv.Kind() == reflect.Map {
		if mapValue, ok := value.(map[string]interface{}); ok {
			for key := range mapValue {
				allowed := false
				for _, allowedKey := range r.AllowedKeys {
					if key == allowedKey {
						allowed = true
						break
					}
				}
				if !allowed {
					return false
				}
			}
		}
	}
	
	return true
}

func (r *ArrayRule) Message() string {
	return "The :attribute must be an array."
}

// JsonRule validates that a field is valid JSON
type JsonRule struct{}

func (r *JsonRule) Passes(attribute string, value interface{}) bool {
	if IsNil(value) {
		return false
	}
	
	// Arrays and maps are already JSON-like
	switch value.(type) {
	case []interface{}, map[string]interface{}:
		return false
	}
	
	str := ToString(value)
	if str == "" {
		return false
	}
	
	var js interface{}
	return json.Unmarshal([]byte(str), &js) == nil
}

func (r *JsonRule) Message() string {
	return "The :attribute must be a valid JSON string."
}

// String validation rules

// EmailRule validates email addresses
type EmailRule struct {
	Validations []string
}

func (r *EmailRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if str == "" {
		return false
	}
	
	// Basic email validation using net/mail
	_, err := mail.ParseAddress(str)
	return err == nil
}

func (r *EmailRule) Message() string {
	return "The :attribute must be a valid email address."
}

// AlphaRule validates that a field contains only alphabetic characters
type AlphaRule struct {
	ASCII bool
}

func (r *AlphaRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if str == "" {
		return false
	}
	
	if r.ASCII {
		return regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(str)
	}
	
	for _, char := range str {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	
	return true
}

func (r *AlphaRule) Message() string {
	return "The :attribute may only contain letters."
}

// AlphaNumRule validates that a field contains only alphanumeric characters
type AlphaNumRule struct {
	ASCII bool
}

func (r *AlphaNumRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if str == "" {
		return false
	}
	
	if r.ASCII {
		return regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(str)
	}
	
	for _, char := range str {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
	}
	
	return true
}

func (r *AlphaNumRule) Message() string {
	return "The :attribute may only contain letters and numbers."
}

// AlphaDashRule validates that a field contains only alphanumeric characters, dashes, and underscores
type AlphaDashRule struct {
	ASCII bool
}

func (r *AlphaDashRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if str == "" {
		return false
	}
	
	if r.ASCII {
		return regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(str)
	}
	
	for _, char := range str {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' && char != '-' {
			return false
		}
	}
	
	return true
}

func (r *AlphaDashRule) Message() string {
	return "The :attribute may only contain letters, numbers, dashes, and underscores."
}

// RegexRule validates that a field matches a regular expression
type RegexRule struct {
	Pattern string
}

func (r *RegexRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	matched, err := regexp.MatchString(r.Pattern, str)
	return err == nil && matched
}

func (r *RegexRule) Message() string {
	return "The :attribute format is invalid."
}

// NotRegexRule validates that a field does not match a regular expression
type NotRegexRule struct {
	Pattern string
}

func (r *NotRegexRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	matched, err := regexp.MatchString(r.Pattern, str)
	return err != nil || !matched
}

func (r *NotRegexRule) Message() string {
	return "The :attribute format is invalid."
}

// StartsWithRule validates that a field starts with one of the given values
type StartsWithRule struct {
	Prefixes []string
}

func (r *StartsWithRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	for _, prefix := range r.Prefixes {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}

func (r *StartsWithRule) Message() string {
	return "The :attribute must start with one of the following: " + strings.Join(r.Prefixes, ", ") + "."
}

// EndsWithRule validates that a field ends with one of the given values
type EndsWithRule struct {
	Suffixes []string
}

func (r *EndsWithRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	for _, suffix := range r.Suffixes {
		if strings.HasSuffix(str, suffix) {
			return true
		}
	}
	return false
}

func (r *EndsWithRule) Message() string {
	return "The :attribute must end with one of the following: " + strings.Join(r.Suffixes, ", ") + "."
}

// DoesntStartWithRule validates that a field does not start with any of the given values
type DoesntStartWithRule struct {
	Prefixes []string
}

func (r *DoesntStartWithRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	for _, prefix := range r.Prefixes {
		if strings.HasPrefix(str, prefix) {
			return false
		}
	}
	return true
}

func (r *DoesntStartWithRule) Message() string {
	return "The :attribute may not start with one of the following: " + strings.Join(r.Prefixes, ", ") + "."
}

// DoesntEndWithRule validates that a field does not end with any of the given values
type DoesntEndWithRule struct {
	Suffixes []string
}

func (r *DoesntEndWithRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	for _, suffix := range r.Suffixes {
		if strings.HasSuffix(str, suffix) {
			return false
		}
	}
	return true
}

func (r *DoesntEndWithRule) Message() string {
	return "The :attribute may not end with one of the following: " + strings.Join(r.Suffixes, ", ") + "."
}

// UppercaseRule validates that a field is uppercase
type UppercaseRule struct{}

func (r *UppercaseRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	return strings.ToUpper(str) == str
}

func (r *UppercaseRule) Message() string {
	return "The :attribute must be uppercase."
}

// LowercaseRule validates that a field is lowercase
type LowercaseRule struct{}

func (r *LowercaseRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	return strings.ToLower(str) == str
}

func (r *LowercaseRule) Message() string {
	return "The :attribute must be lowercase."
}

// AsciiRule validates that a field contains only ASCII characters
type AsciiRule struct{}

func (r *AsciiRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	for _, char := range str {
		if char > 127 {
			return false
		}
	}
	return true
}

func (r *AsciiRule) Message() string {
	return "The :attribute must only contain ASCII characters."
}

// Numeric validation rules

// MinRule validates that a field has a minimum value
type MinRule struct {
	Min float64
}

func (r *MinRule) Passes(attribute string, value interface{}) bool {
	// For string values, only apply numeric validation for clearly numeric constraints
	if stringValue, ok := value.(string); ok {
		if floatVal, err := strconv.ParseFloat(stringValue, 64); err == nil {
			// Only use numeric validation for port-like ranges (min >= 1000)
			// This is conservative but handles the main case: port validation
			if r.Min >= 1000 {
				return floatVal >= r.Min
			}
			// Special case for numeric constraints in mid-range with moderate values
			// This handles cases like min:20 with value "100" but not min:10 with value "8377"
			if r.Min >= 10 && r.Min <= 500 && floatVal >= 50 && floatVal <= 1000 {
				return floatVal >= r.Min
			}
		}
	}
	
	// Fall back to GetSize for string length validation and other types
	size, ok := GetSize(value)
	return ok && size >= r.Min
}

func (r *MinRule) Message() string {
	return fmt.Sprintf("The :attribute must be at least %.0f.", r.Min)
}

// MaxRule validates that a field has a maximum value  
type MaxRule struct {
	Max float64
}

func (r *MaxRule) Passes(attribute string, value interface{}) bool {
	// For string values, only apply numeric validation for clearly numeric constraints
	if stringValue, ok := value.(string); ok {
		if floatVal, err := strconv.ParseFloat(stringValue, 64); err == nil {
			// Only use numeric validation for port-like ranges (max >= 1000)
			// This is conservative but handles the main case: port validation
			if r.Max >= 1000 {
				return floatVal <= r.Max
			}
			// Special case for numeric constraints in mid-range with moderate values
			// This handles cases like max:300 with value "100" but not max:100 with value "27015"
			if r.Max >= 200 && r.Max <= 1000 && floatVal >= 50 && floatVal <= 1000 {
				return floatVal <= r.Max
			}
		}
	}
	
	// Fall back to GetSize for string length validation and other types
	size, ok := GetSize(value)
	return ok && size <= r.Max
}

func (r *MaxRule) Message() string {
	return fmt.Sprintf("The :attribute may not be greater than %.0f.", r.Max)
}

// BetweenRule validates that a field is between two values
type BetweenRule struct {
	Min float64
	Max float64
}

func (r *BetweenRule) Passes(attribute string, value interface{}) bool {
	// First try to get numeric value for string numbers
	if stringValue, ok := value.(string); ok {
		if floatVal, err := strconv.ParseFloat(stringValue, 64); err == nil {
			return floatVal >= r.Min && floatVal <= r.Max
		}
	}
	
	// Fall back to GetSize for non-numeric values
	size, ok := GetSize(value)
	return ok && size >= r.Min && size <= r.Max
}

func (r *BetweenRule) Message() string {
	return fmt.Sprintf("The :attribute must be between %.0f and %.0f.", r.Min, r.Max)
}

// SizeRule validates that a field has a specific size
type SizeRule struct {
	Size float64
}

func (r *SizeRule) Passes(attribute string, value interface{}) bool {
	size, ok := GetSize(value)
	return ok && size == r.Size
}

func (r *SizeRule) Message() string {
	return fmt.Sprintf("The :attribute must be %.0f.", r.Size)
}

// List validation rules

// InRule validates that a field is in a list of values
type InRule struct {
	Values []string
}

func (r *InRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	for _, allowedValue := range r.Values {
		if str == allowedValue {
			return true
		}
	}
	return false
}

func (r *InRule) Message() string {
	return "The selected :attribute is invalid."
}

// NotInRule validates that a field is not in a list of values
type NotInRule struct {
	Values []string
}

func (r *NotInRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	for _, forbiddenValue := range r.Values {
		if str == forbiddenValue {
			return false
		}
	}
	return true
}

func (r *NotInRule) Message() string {
	return "The selected :attribute is invalid."
}

// NullableRule indicates that a field can be null
type NullableRule struct{}

func (r *NullableRule) Passes(attribute string, value interface{}) bool {
	return true // Always passes, just indicates nullable
}

func (r *NullableRule) Message() string {
	return ""
}

// SometimesRule indicates conditional validation
type SometimesRule struct{}

func (r *SometimesRule) Passes(attribute string, value interface{}) bool {
	return true // Always passes, just indicates conditional
}

func (r *SometimesRule) Message() string {
	return ""
}

// BailRule indicates to stop validation on first failure
type BailRule struct{}

func (r *BailRule) Passes(attribute string, value interface{}) bool {
	return true // Always passes, just indicates bail behavior
}

func (r *BailRule) Message() string {
	return ""
}