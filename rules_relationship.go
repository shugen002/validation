package validation

import (
	"reflect"
	"strings"
)

// Field relationship validation rules

// SameRule validates that a field has the same value as another field
type SameRule struct {
	Field string
	data  map[string]interface{}
}

func (r *SameRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		return false
	}
	return value == otherValue
}

func (r *SameRule) Message() string {
	return "The :attribute and " + r.Field + " must match."
}

func (r *SameRule) SetData(data map[string]interface{}) {
	r.data = data
}

// DifferentRule validates that a field has a different value from other fields
type DifferentRule struct {
	Fields []string
	data   map[string]interface{}
}

func (r *DifferentRule) Passes(attribute string, value interface{}) bool {
	for _, field := range r.Fields {
		if otherValue, exists := r.data[field]; exists {
			if value == otherValue {
				return false
			}
		}
	}
	return true
}

func (r *DifferentRule) Message() string {
	return "The :attribute and " + strings.Join(r.Fields, ", ") + " must be different."
}

func (r *DifferentRule) SetData(data map[string]interface{}) {
	r.data = data
}

// ConfirmedRule validates that a field has a matching confirmation field
type ConfirmedRule struct {
	Field string // optional custom confirmation field name
	data  map[string]interface{}
}

func (r *ConfirmedRule) Passes(attribute string, value interface{}) bool {
	confirmationField := r.Field
	if confirmationField == "" {
		confirmationField = attribute + "_confirmation"
	}
	
	confirmationValue, exists := r.data[confirmationField]
	if !exists {
		return false
	}
	
	return value == confirmationValue
}

func (r *ConfirmedRule) Message() string {
	return "The :attribute confirmation does not match."
}

func (r *ConfirmedRule) SetData(data map[string]interface{}) {
	r.data = data
}

// RequiredIfRule validates that a field is required when another field equals specified value
type RequiredIfRule struct {
	Field string
	Value string
	data  map[string]interface{}
}

func (r *RequiredIfRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		return true // If the other field doesn't exist, this rule doesn't apply
	}
	
	// Convert both values to strings for comparison
	otherStr := ToString(otherValue)
	if otherStr == r.Value {
		// The condition is met, so this field is required
		requiredRule := &RequiredRule{}
		return requiredRule.Passes(attribute, value)
	}
	
	return true // Condition not met, so this rule passes
}

func (r *RequiredIfRule) Message() string {
	return "The :attribute field is required when " + r.Field + " is " + r.Value + "."
}

func (r *RequiredIfRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *RequiredIfRule) IsImplicit() bool {
	return true
}

// RequiredUnlessRule validates that a field is required unless another field equals specified value
type RequiredUnlessRule struct {
	Field string
	Value string
	data  map[string]interface{}
}

func (r *RequiredUnlessRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		// If the other field doesn't exist, this field is required
		requiredRule := &RequiredRule{}
		return requiredRule.Passes(attribute, value)
	}
	
	// Convert both values to strings for comparison
	otherStr := ToString(otherValue)
	if otherStr != r.Value {
		// The condition is not met, so this field is required
		requiredRule := &RequiredRule{}
		return requiredRule.Passes(attribute, value)
	}
	
	return true // Condition met (field equals exception value), so this rule passes
}

func (r *RequiredUnlessRule) Message() string {
	return "The :attribute field is required unless " + r.Field + " is " + r.Value + "."
}

func (r *RequiredUnlessRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *RequiredUnlessRule) IsImplicit() bool {
	return true
}

// RequiredWithRule validates that a field is required when any of the other fields are present
type RequiredWithRule struct {
	Fields []string
	data   map[string]interface{}
}

func (r *RequiredWithRule) Passes(attribute string, value interface{}) bool {
	// Check if any of the specified fields are present
	anyFieldPresent := false
	for _, field := range r.Fields {
		if _, exists := r.data[field]; exists {
			anyFieldPresent = true
			break
		}
	}
	
	if anyFieldPresent {
		// At least one field is present, so this field is required
		requiredRule := &RequiredRule{}
		return requiredRule.Passes(attribute, value)
	}
	
	return true // None of the fields are present, so this rule passes
}

func (r *RequiredWithRule) Message() string {
	return "The :attribute field is required when " + strings.Join(r.Fields, " / ") + " is present."
}

func (r *RequiredWithRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *RequiredWithRule) IsImplicit() bool {
	return true
}

// RequiredWithoutRule validates that a field is required when any of the other fields are not present
type RequiredWithoutRule struct {
	Fields []string
	data   map[string]interface{}
}

func (r *RequiredWithoutRule) Passes(attribute string, value interface{}) bool {
	// Check if any of the specified fields are missing
	anyFieldMissing := false
	for _, field := range r.Fields {
		if _, exists := r.data[field]; !exists {
			anyFieldMissing = true
			break
		}
	}
	
	if anyFieldMissing {
		// At least one field is missing, so this field is required
		requiredRule := &RequiredRule{}
		return requiredRule.Passes(attribute, value)
	}
	
	return true // All fields are present, so this rule passes
}

func (r *RequiredWithoutRule) Message() string {
	return "The :attribute field is required when " + strings.Join(r.Fields, " / ") + " is not present."
}

func (r *RequiredWithoutRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *RequiredWithoutRule) IsImplicit() bool {
	return true
}

// RequiredWithAllRule validates that a field is required when all of the other fields are present
type RequiredWithAllRule struct {
	Fields []string
	data   map[string]interface{}
}

func (r *RequiredWithAllRule) Passes(attribute string, value interface{}) bool {
	// Check if all of the specified fields are present
	allFieldsPresent := true
	for _, field := range r.Fields {
		if _, exists := r.data[field]; !exists {
			allFieldsPresent = false
			break
		}
	}
	
	if allFieldsPresent {
		// All fields are present, so this field is required
		requiredRule := &RequiredRule{}
		return requiredRule.Passes(attribute, value)
	}
	
	return true // Not all fields are present, so this rule passes
}

func (r *RequiredWithAllRule) Message() string {
	return "The :attribute field is required when " + strings.Join(r.Fields, " / ") + " are present."
}

func (r *RequiredWithAllRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *RequiredWithAllRule) IsImplicit() bool {
	return true
}

// RequiredIfAcceptedRule validates that a field is required when another field is accepted
type RequiredIfAcceptedRule struct {
	Field string
	data  map[string]interface{}
}

func (r *RequiredIfAcceptedRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		return true // If the other field doesn't exist, this rule doesn't apply
	}
	
	// Check if the other field is accepted
	acceptedRule := &AcceptedRule{}
	if acceptedRule.Passes(r.Field, otherValue) {
		// The other field is accepted, so this field is required
		requiredRule := &RequiredRule{}
		return requiredRule.Passes(attribute, value)
	}
	
	return true // Other field is not accepted, so this rule passes
}

func (r *RequiredIfAcceptedRule) Message() string {
	return "The :attribute field is required when " + r.Field + " is accepted."
}

func (r *RequiredIfAcceptedRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *RequiredIfAcceptedRule) IsImplicit() bool {
	return true
}

// RequiredIfDeclinedRule validates that a field is required when another field is declined
type RequiredIfDeclinedRule struct {
	Field string
	data  map[string]interface{}
}

func (r *RequiredIfDeclinedRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		return true // If the other field doesn't exist, this rule doesn't apply
	}
	
	// Check if the other field is declined
	declinedRule := &DeclinedRule{}
	if declinedRule.Passes(r.Field, otherValue) {
		// The other field is declined, so this field is required
		requiredRule := &RequiredRule{}
		return requiredRule.Passes(attribute, value)
	}
	
	return true // Other field is not declined, so this rule passes
}

func (r *RequiredIfDeclinedRule) Message() string {
	return "The :attribute field is required when " + r.Field + " is declined."
}

func (r *RequiredIfDeclinedRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *RequiredIfDeclinedRule) IsImplicit() bool {
	return true
}

// RequiredArrayKeysRule validates that an array must contain the specified keys
type RequiredArrayKeysRule struct {
	Keys []string
	data map[string]interface{}
}

func (r *RequiredArrayKeysRule) Passes(attribute string, value interface{}) bool {
	rv := reflect.ValueOf(value)
	
	// Must be a map
	if rv.Kind() != reflect.Map {
		return false
	}
	
	// Check that all required keys are present
	for _, requiredKey := range r.Keys {
		found := false
		for _, key := range rv.MapKeys() {
			if ToString(key.Interface()) == requiredKey {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	return true
}

func (r *RequiredArrayKeysRule) Message() string {
	return "The :attribute array must contain entries for: " + strings.Join(r.Keys, ", ") + "."
}

func (r *RequiredArrayKeysRule) SetData(data map[string]interface{}) {
	r.data = data
}