package validation

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// Special validation rules

// UuidRule validates that a field is a valid UUID
type UuidRule struct {
	Version interface{} // int for specific version, "max" for maximum version, or nil for any
}

func (r *UuidRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if str == "" {
		return false
	}
	
	parsedUUID, err := uuid.Parse(str)
	if err != nil {
		return false
	}
	
	if r.Version != nil {
		switch v := r.Version.(type) {
		case int:
			if parsedUUID.Version() != uuid.Version(v) {
				return false
			}
		case string:
			if v == "max" {
				// Check if it's a valid version (1-8)
				version := parsedUUID.Version()
				if version < 1 || version > 8 {
					return false
				}
			}
		}
	}
	
	return true
}

func (r *UuidRule) Message() string {
	if r.Version != nil {
		return "The :attribute must be a valid UUID version " + ToString(r.Version) + "."
	}
	return "The :attribute must be a valid UUID."
}

// UlidRule validates that a field is a valid ULID
type UlidRule struct{}

func (r *UlidRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if len(str) != 26 {
		return false
	}
	
	// ULID uses Crockford's Base32 encoding
	// Valid characters: 0-9, A-Z (excluding I, L, O, U)
	validChars := "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	
	for _, char := range str {
		found := false
		for _, validChar := range validChars {
			if char == validChar {
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

func (r *UlidRule) Message() string {
	return "The :attribute must be a valid ULID."
}

// NullableRule indicates that a field can be null
type NullableRule struct{}

func (r *NullableRule) Passes(attribute string, value interface{}) bool {
	return true // This rule always passes - it just indicates the field can be null
}

func (r *NullableRule) Message() string {
	return "" // This rule doesn't have a validation message
}

// SometimesRule indicates conditional validation
type SometimesRule struct{}

func (r *SometimesRule) Passes(attribute string, value interface{}) bool {
	return true // This rule always passes - it's just a marker
}

func (r *SometimesRule) Message() string {
	return "" // This rule doesn't have a validation message
}

// BailRule indicates to stop validation on first failure
type BailRule struct{}

func (r *BailRule) Passes(attribute string, value interface{}) bool {
	return true // This rule always passes - it's just a marker
}

func (r *BailRule) Message() string {
	return "" // This rule doesn't have a validation message
}

// MissingRule validates that a field is not present
type MissingRule struct {
	validator Validator
}

func (r *MissingRule) Passes(attribute string, value interface{}) bool {
	return !r.validator.HasField(attribute)
}

func (r *MissingRule) Message() string {
	return "The :attribute field must not be present."
}

func (r *MissingRule) IsImplicit() bool {
	return true
}

func (r *MissingRule) SetValidator(validator Validator) {
	r.validator = validator
}

// MissingIfRule validates that a field is not present when another field equals specified value
type MissingIfRule struct {
	Field string
	Value string
	data  map[string]interface{}
	validator Validator
}

func (r *MissingIfRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		return true // If the other field doesn't exist, this rule doesn't apply
	}
	
	// Convert both values to strings for comparison
	otherStr := ToString(otherValue)
	if otherStr == r.Value {
		// The condition is met, so this field must not be present
		return !r.validator.HasField(attribute)
	}
	
	return true // Condition not met, so this rule passes
}

func (r *MissingIfRule) Message() string {
	return "The :attribute field must not be present when " + r.Field + " is " + r.Value + "."
}

func (r *MissingIfRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *MissingIfRule) SetValidator(validator Validator) {
	r.validator = validator
}

func (r *MissingIfRule) IsImplicit() bool {
	return true
}

// MissingUnlessRule validates that a field is not present unless another field equals specified value
type MissingUnlessRule struct {
	Field string
	Value string
	data  map[string]interface{}
	validator Validator
}

func (r *MissingUnlessRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		// If the other field doesn't exist, this field must not be present
		return !r.validator.HasField(attribute)
	}
	
	// Convert both values to strings for comparison
	otherStr := ToString(otherValue)
	if otherStr != r.Value {
		// The condition is not met, so this field must not be present
		return !r.validator.HasField(attribute)
	}
	
	return true // Condition met (field equals exception value), so this rule passes
}

func (r *MissingUnlessRule) Message() string {
	return "The :attribute field must not be present unless " + r.Field + " is " + r.Value + "."
}

func (r *MissingUnlessRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *MissingUnlessRule) SetValidator(validator Validator) {
	r.validator = validator
}

func (r *MissingUnlessRule) IsImplicit() bool {
	return true
}

// MissingWithRule validates that a field is not present when any of the other fields are present
type MissingWithRule struct {
	Fields []string
	data   map[string]interface{}
	validator Validator
}

func (r *MissingWithRule) Passes(attribute string, value interface{}) bool {
	// Check if any of the specified fields are present
	anyFieldPresent := false
	for _, field := range r.Fields {
		if _, exists := r.data[field]; exists {
			anyFieldPresent = true
			break
		}
	}
	
	if anyFieldPresent {
		// At least one field is present, so this field must not be present
		return !r.validator.HasField(attribute)
	}
	
	return true // None of the fields are present, so this rule passes
}

func (r *MissingWithRule) Message() string {
	return "The :attribute field must not be present when " + strings.Join(r.Fields, " / ") + " is present."
}

func (r *MissingWithRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *MissingWithRule) SetValidator(validator Validator) {
	r.validator = validator
}

func (r *MissingWithRule) IsImplicit() bool {
	return true
}

// MissingWithAllRule validates that a field is not present when all of the other fields are present
type MissingWithAllRule struct {
	Fields []string
	data   map[string]interface{}
	validator Validator
}

func (r *MissingWithAllRule) Passes(attribute string, value interface{}) bool {
	// Check if all of the specified fields are present
	allFieldsPresent := true
	for _, field := range r.Fields {
		if _, exists := r.data[field]; !exists {
			allFieldsPresent = false
			break
		}
	}
	
	if allFieldsPresent {
		// All fields are present, so this field must not be present
		return !r.validator.HasField(attribute)
	}
	
	return true // Not all fields are present, so this rule passes
}

func (r *MissingWithAllRule) Message() string {
	return "The :attribute field must not be present when " + strings.Join(r.Fields, " / ") + " are present."
}

func (r *MissingWithAllRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *MissingWithAllRule) SetValidator(validator Validator) {
	r.validator = validator
}

func (r *MissingWithAllRule) IsImplicit() bool {
	return true
}

// PresentIfRule validates that a field is present when another field equals specified value
type PresentIfRule struct {
	Field string
	Value string
	data  map[string]interface{}
	validator Validator
}

func (r *PresentIfRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		return true // If the other field doesn't exist, this rule doesn't apply
	}
	
	// Convert both values to strings for comparison
	otherStr := ToString(otherValue)
	if otherStr == r.Value {
		// The condition is met, so this field must be present
		return r.validator.HasField(attribute)
	}
	
	return true // Condition not met, so this rule passes
}

func (r *PresentIfRule) Message() string {
	return "The :attribute field must be present when " + r.Field + " is " + r.Value + "."
}

func (r *PresentIfRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *PresentIfRule) SetValidator(validator Validator) {
	r.validator = validator
}

func (r *PresentIfRule) IsImplicit() bool {
	return true
}

// PresentUnlessRule validates that a field is present unless another field equals specified value
type PresentUnlessRule struct {
	Field string
	Value string
	data  map[string]interface{}
	validator Validator
}

func (r *PresentUnlessRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		// If the other field doesn't exist, this field must be present
		return r.validator.HasField(attribute)
	}
	
	// Convert both values to strings for comparison
	otherStr := ToString(otherValue)
	if otherStr != r.Value {
		// The condition is not met, so this field must be present
		return r.validator.HasField(attribute)
	}
	
	return true // Condition met (field equals exception value), so this rule passes
}

func (r *PresentUnlessRule) Message() string {
	return "The :attribute field must be present unless " + r.Field + " is " + r.Value + "."
}

func (r *PresentUnlessRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *PresentUnlessRule) SetValidator(validator Validator) {
	r.validator = validator
}

func (r *PresentUnlessRule) IsImplicit() bool {
	return true
}

// PresentWithRule validates that a field is present when any of the other fields are present
type PresentWithRule struct {
	Fields []string
	data   map[string]interface{}
	validator Validator
}

func (r *PresentWithRule) Passes(attribute string, value interface{}) bool {
	// Check if any of the specified fields are present
	anyFieldPresent := false
	for _, field := range r.Fields {
		if _, exists := r.data[field]; exists {
			anyFieldPresent = true
			break
		}
	}
	
	if anyFieldPresent {
		// At least one field is present, so this field must be present
		return r.validator.HasField(attribute)
	}
	
	return true // None of the fields are present, so this rule passes
}

func (r *PresentWithRule) Message() string {
	return "The :attribute field must be present when " + strings.Join(r.Fields, " / ") + " is present."
}

func (r *PresentWithRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *PresentWithRule) SetValidator(validator Validator) {
	r.validator = validator
}

func (r *PresentWithRule) IsImplicit() bool {
	return true
}

// PresentWithAllRule validates that a field is present when all of the other fields are present
type PresentWithAllRule struct {
	Fields []string
	data   map[string]interface{}
	validator Validator
}

func (r *PresentWithAllRule) Passes(attribute string, value interface{}) bool {
	// Check if all of the specified fields are present
	allFieldsPresent := true
	for _, field := range r.Fields {
		if _, exists := r.data[field]; !exists {
			allFieldsPresent = false
			break
		}
	}
	
	if allFieldsPresent {
		// All fields are present, so this field must be present
		return r.validator.HasField(attribute)
	}
	
	return true // Not all fields are present, so this rule passes
}

func (r *PresentWithAllRule) Message() string {
	return "The :attribute field must be present when " + strings.Join(r.Fields, " / ") + " are present."
}

func (r *PresentWithAllRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *PresentWithAllRule) SetValidator(validator Validator) {
	r.validator = validator
}

func (r *PresentWithAllRule) IsImplicit() bool {
	return true
}

// ProhibitedIfAcceptedRule validates that a field is prohibited when another field is accepted
type ProhibitedIfAcceptedRule struct {
	Field string
	data  map[string]interface{}
}

func (r *ProhibitedIfAcceptedRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		return true // If the other field doesn't exist, this rule doesn't apply
	}
	
	// Check if the other field is accepted
	acceptedRule := &AcceptedRule{}
	if acceptedRule.Passes(r.Field, otherValue) {
		// The other field is accepted, so this field must be prohibited
		prohibitedRule := &ProhibitedRule{}
		return prohibitedRule.Passes(attribute, value)
	}
	
	return true // Other field is not accepted, so this rule passes
}

func (r *ProhibitedIfAcceptedRule) Message() string {
	return "The :attribute field is prohibited when " + r.Field + " is accepted."
}

func (r *ProhibitedIfAcceptedRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *ProhibitedIfAcceptedRule) IsImplicit() bool {
	return true
}

// ProhibitedIfDeclinedRule validates that a field is prohibited when another field is declined
type ProhibitedIfDeclinedRule struct {
	Field string
	data  map[string]interface{}
}

func (r *ProhibitedIfDeclinedRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		return true // If the other field doesn't exist, this rule doesn't apply
	}
	
	// Check if the other field is declined
	declinedRule := &DeclinedRule{}
	if declinedRule.Passes(r.Field, otherValue) {
		// The other field is declined, so this field must be prohibited
		prohibitedRule := &ProhibitedRule{}
		return prohibitedRule.Passes(attribute, value)
	}
	
	return true // Other field is not declined, so this rule passes
}

func (r *ProhibitedIfDeclinedRule) Message() string {
	return "The :attribute field is prohibited when " + r.Field + " is declined."
}

func (r *ProhibitedIfDeclinedRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *ProhibitedIfDeclinedRule) IsImplicit() bool {
	return true
}

// ProhibitedUnlessRule validates that a field is prohibited unless another field equals specified value
type ProhibitedUnlessRule struct {
	Field string
	Value string
	data  map[string]interface{}
}

func (r *ProhibitedUnlessRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		// If the other field doesn't exist, this field must be prohibited
		prohibitedRule := &ProhibitedRule{}
		return prohibitedRule.Passes(attribute, value)
	}
	
	// Convert both values to strings for comparison
	otherStr := ToString(otherValue)
	if otherStr != r.Value {
		// The condition is not met, so this field must be prohibited
		prohibitedRule := &ProhibitedRule{}
		return prohibitedRule.Passes(attribute, value)
	}
	
	return true // Condition met (field equals exception value), so this rule passes
}

func (r *ProhibitedUnlessRule) Message() string {
	return "The :attribute field is prohibited unless " + r.Field + " is " + r.Value + "."
}

func (r *ProhibitedUnlessRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *ProhibitedUnlessRule) IsImplicit() bool {
	return true
}

// ProhibitsRule validates that if this field is present, the specified fields must not be present
type ProhibitsRule struct {
	Fields []string
	data   map[string]interface{}
}

func (r *ProhibitsRule) Passes(attribute string, value interface{}) bool {
	// Check if any of the prohibited fields are present
	for _, field := range r.Fields {
		if _, exists := r.data[field]; exists {
			return false
		}
	}
	
	return true
}

func (r *ProhibitsRule) Message() string {
	return "The :attribute field prohibits " + strings.Join(r.Fields, " / ") + " from being present."
}

func (r *ProhibitsRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *ProhibitsRule) IsImplicit() bool {
	return true
}

// GreaterThanRule validates that a field is greater than another field
type GreaterThanRule struct {
	Field string
	data  map[string]interface{}
}

func (r *GreaterThanRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		return false
	}
	
	thisVal, err1 := strconv.ParseFloat(ToString(value), 64)
	otherVal, err2 := strconv.ParseFloat(ToString(otherValue), 64)
	
	if err1 != nil || err2 != nil {
		return false
	}
	
	return thisVal > otherVal
}

func (r *GreaterThanRule) Message() string {
	return "The :attribute must be greater than " + r.Field + "."
}

func (r *GreaterThanRule) SetData(data map[string]interface{}) {
	r.data = data
}

// GreaterThanOrEqualRule validates that a field is greater than or equal to another field
type GreaterThanOrEqualRule struct {
	Field string
	data  map[string]interface{}
}

func (r *GreaterThanOrEqualRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		return false
	}
	
	thisVal, err1 := strconv.ParseFloat(ToString(value), 64)
	otherVal, err2 := strconv.ParseFloat(ToString(otherValue), 64)
	
	if err1 != nil || err2 != nil {
		return false
	}
	
	return thisVal >= otherVal
}

func (r *GreaterThanOrEqualRule) Message() string {
	return "The :attribute must be greater than or equal to " + r.Field + "."
}

func (r *GreaterThanOrEqualRule) SetData(data map[string]interface{}) {
	r.data = data
}

// LessThanRule validates that a field is less than another field
type LessThanRule struct {
	Field string
	data  map[string]interface{}
}

func (r *LessThanRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		return false
	}
	
	thisVal, err1 := strconv.ParseFloat(ToString(value), 64)
	otherVal, err2 := strconv.ParseFloat(ToString(otherValue), 64)
	
	if err1 != nil || err2 != nil {
		return false
	}
	
	return thisVal < otherVal
}

func (r *LessThanRule) Message() string {
	return "The :attribute must be less than " + r.Field + "."
}

func (r *LessThanRule) SetData(data map[string]interface{}) {
	r.data = data
}

// LessThanOrEqualRule validates that a field is less than or equal to another field
type LessThanOrEqualRule struct {
	Field string
	data  map[string]interface{}
}

func (r *LessThanOrEqualRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		return false
	}
	
	thisVal, err1 := strconv.ParseFloat(ToString(value), 64)
	otherVal, err2 := strconv.ParseFloat(ToString(otherValue), 64)
	
	if err1 != nil || err2 != nil {
		return false
	}
	
	return thisVal <= otherVal
}

func (r *LessThanOrEqualRule) Message() string {
	return "The :attribute must be less than or equal to " + r.Field + "."
}

func (r *LessThanOrEqualRule) SetData(data map[string]interface{}) {
	r.data = data
}