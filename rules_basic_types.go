package validation

import (
	"reflect"
	"strings"
)

// Basic type validation rules

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

// IntegerRule validates that a field is an integer or a string representing an integer
type IntegerRule struct {
	Strict bool // if true, only accept actual integer types
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
	
	return IsInteger(value)
}

func (r *IntegerRule) Message() string {
	return "The :attribute must be an integer."
}

// NumericRule validates that a field is numeric (int, float, or string containing numeric value)
type NumericRule struct {
	Strict bool // if true, only accept actual numeric types, not string representations
}

func (r *NumericRule) Passes(attribute string, value interface{}) bool {
	if r.Strict {
		switch value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
			float32, float64:
			return true
		default:
			return false
		}
	}
	
	return IsNumeric(value)
}

func (r *NumericRule) Message() string {
	return "The :attribute must be a number."
}

// BooleanRule validates that a field is a boolean or can be converted to boolean
type BooleanRule struct {
	Strict bool // if true, only accept actual boolean types
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
		lower := strings.ToLower(v)
		return lower == "true" || lower == "false" || lower == "1" || lower == "0" ||
			lower == "yes" || lower == "no" || lower == "on" || lower == "off"
	case int, int8, int16, int32, int64:
		val := reflect.ValueOf(v).Int()
		return val == 0 || val == 1
	case uint, uint8, uint16, uint32, uint64:
		val := reflect.ValueOf(v).Uint()
		return val == 0 || val == 1
	case float32, float64:
		val := reflect.ValueOf(v).Float()
		return val == 0.0 || val == 1.0
	}
	
	return false
}

func (r *BooleanRule) Message() string {
	return "The :attribute field must be true or false."
}

// AcceptedRule validates that a field is accepted (useful for terms of service, agreements, etc.)
type AcceptedRule struct{}

func (r *AcceptedRule) Passes(attribute string, value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		lower := strings.ToLower(v)
		return lower == "yes" || lower == "on" || lower == "1" || lower == "true"
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Int() == 1
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint() == 1
	case float32, float64:
		return reflect.ValueOf(v).Float() == 1.0
	}
	return false
}

func (r *AcceptedRule) Message() string {
	return "The :attribute must be accepted."
}

func (r *AcceptedRule) IsImplicit() bool {
	return true
}

// AcceptedIfRule validates that a field is accepted if another field equals the specified value
type AcceptedIfRule struct {
	Field string
	Value string
	data  map[string]interface{}
}

func (r *AcceptedIfRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		return true // If the other field doesn't exist, this rule doesn't apply
	}
	
	// Convert both values to strings for comparison
	otherStr := ToString(otherValue)
	if otherStr == r.Value {
		// The condition is met, so this field must be accepted
		acceptedRule := &AcceptedRule{}
		return acceptedRule.Passes(attribute, value)
	}
	
	return true // Condition not met, so this rule passes
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

// DeclinedRule validates that a field is declined (opposite of accepted)
type DeclinedRule struct{}

func (r *DeclinedRule) Passes(attribute string, value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return !v
	case string:
		lower := strings.ToLower(v)
		return lower == "no" || lower == "off" || lower == "0" || lower == "false"
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Int() == 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint() == 0
	case float32, float64:
		return reflect.ValueOf(v).Float() == 0.0
	}
	return false
}

func (r *DeclinedRule) Message() string {
	return "The :attribute must be declined."
}

func (r *DeclinedRule) IsImplicit() bool {
	return true
}

// DeclinedIfRule validates that a field is declined if another field equals the specified value
type DeclinedIfRule struct {
	Field string
	Value string
	data  map[string]interface{}
}

func (r *DeclinedIfRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		return true // If the other field doesn't exist, this rule doesn't apply
	}
	
	// Convert both values to strings for comparison
	otherStr := ToString(otherValue)
	if otherStr == r.Value {
		// The condition is met, so this field must be declined
		declinedRule := &DeclinedRule{}
		return declinedRule.Passes(attribute, value)
	}
	
	return true // Condition not met, so this rule passes
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

// ArrayRule validates that a field is an array/slice with optional key restrictions
type ArrayRule struct {
	AllowedKeys []string
}

func (r *ArrayRule) Passes(attribute string, value interface{}) bool {
	rv := reflect.ValueOf(value)
	
	// Check if it's a slice or array
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array && rv.Kind() != reflect.Map {
		return false
	}
	
	// If we have allowed keys restriction and it's a map, check keys
	if len(r.AllowedKeys) > 0 && rv.Kind() == reflect.Map {
		for _, key := range rv.MapKeys() {
			keyStr := ToString(key.Interface())
			allowed := false
			for _, allowedKey := range r.AllowedKeys {
				if keyStr == allowedKey {
					allowed = true
					break
				}
			}
			if !allowed {
				return false
			}
		}
	}
	
	return true
}

func (r *ArrayRule) Message() string {
	if len(r.AllowedKeys) > 0 {
		return "The :attribute must be an array with allowed keys: " + strings.Join(r.AllowedKeys, ", ") + "."
	}
	return "The :attribute must be an array."
}

// JsonRule validates that a field contains valid JSON
type JsonRule struct{}

func (r *JsonRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	return IsJSON(str)
}

func (r *JsonRule) Message() string {
	return "The :attribute must be a valid JSON string."
}

// FilledRule validates that a field has a value when present (can be absent)
type FilledRule struct {
	validator Validator
}

func (r *FilledRule) Passes(attribute string, value interface{}) bool {
	// If the field is not in the data, the rule passes
	if !r.validator.HasField(attribute) {
		return true
	}
	
	// If the field is present, it must not be empty
	requiredRule := &RequiredRule{}
	return requiredRule.Passes(attribute, value)
}

func (r *FilledRule) Message() string {
	return "The :attribute field must have a value when present."
}

func (r *FilledRule) IsImplicit() bool {
	return true
}

func (r *FilledRule) SetValidator(validator Validator) {
	r.validator = validator
}

// PresentRule validates that a field is present in input (can be empty)
type PresentRule struct {
	validator Validator
}

func (r *PresentRule) Passes(attribute string, value interface{}) bool {
	return r.validator.HasField(attribute)
}

func (r *PresentRule) Message() string {
	return "The :attribute field must be present."
}

func (r *PresentRule) IsImplicit() bool {
	return true
}

func (r *PresentRule) SetValidator(validator Validator) {
	r.validator = validator
}

// ProhibitedRule validates that a field is not present or is empty
type ProhibitedRule struct{}

func (r *ProhibitedRule) Passes(attribute string, value interface{}) bool {
	if IsNil(value) {
		return true
	}
	
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case []interface{}, map[string]interface{}:
		return reflect.ValueOf(v).Len() == 0
	default:
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
			return rv.Len() == 0
		}
	}
	
	return false
}

func (r *ProhibitedRule) Message() string {
	return "The :attribute field is prohibited."
}

func (r *ProhibitedRule) IsImplicit() bool {
	return true
}