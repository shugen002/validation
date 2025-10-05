package validation

import (
	"fmt"
	"strings"
)

// Booleans:
// Accepted
// Accepted If
// Boolean
// Declined
// Declined If

// accepted
// The field under validation must be "yes", "on", 1, "1", true, or "true". This is useful for validating "Terms of Service" acceptance or similar fields.
func constructAcceptedRule(_cfg map[string]interface{}, _args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		val := strings.ToLower(ctx.FieldValue)
		if val == "yes" || val == "on" || val == "1" || val == "true" {
			return true, nil
		}
		return false, fmt.Errorf("the %s field must be accepted", ctx.FieldName)
	}, nil
}

// accepted_if:anotherfield,value,...
// The field under validation must be "yes", "on", 1, "1", true, or "true" if another field under validation is equal to a specified value. This is useful for validating "Terms of Service" acceptance or similar fields.
func constructAcceptedIfRule(cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("accepted_if rule requires at least 2 arguments")
	}
	otherField := args[0]
	expectedValues := args[1:]

	return func(ctx *ValidationContext) (bool, error) {
		otherValue, ok := ctx.Raw[otherField]
		if !ok {
			return false, fmt.Errorf("the %s field is not present", otherField)
		}
		match := false
		for _, v := range expectedValues {
			if otherValue == v {
				match = true
				break
			}
		}
		if match {
			val := strings.ToLower(ctx.FieldValue)
			if val == "yes" || val == "on" || val == "1" || val == "true" {
				return true, nil
			}
			return false, fmt.Errorf("the %s field must be accepted when %s is %s", ctx.FieldName, otherField, strings.Join(expectedValues, ", "))
		}
		return true, nil
	}, nil
}

// boolean
// The field under validation must be able to be cast as a boolean. Accepted input are true, false, 1, 0, "1", and "0".
// You may use the strict parameter to only consider the field valid if its value is true or false:
// 'foo' => 'boolean:strict'
func constructBooleanRule(cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	strict := false
	if len(args) > 0 && args[0] == "strict" {
		strict = true
	}

	return func(ctx *ValidationContext) (bool, error) {
		val := strings.ToLower(ctx.FieldValue)
		if strict {
			if val == "true" || val == "false" {
				return true, nil
			}
			return false, fmt.Errorf("the %s field must be a boolean (true or false)", ctx.FieldName)
		} else {
			if val == "true" || val == "false" || val == "1" || val == "0" {
				return true, nil
			}
			return false, fmt.Errorf("the %s field must be a boolean", ctx.FieldName)
		}
	}, nil
}

// declined
// The field under validation must be "no", "off", 0, "0", false, or "false".

func constructDeclinedRule(_cfg map[string]interface{}, _args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		val := strings.ToLower(ctx.FieldValue)
		if val == "no" || val == "off" || val == "0" || val == "false" {
			return true, nil
		}
		return false, fmt.Errorf("the %s field must be declined", ctx.FieldName)
	}, nil
}

// declined_if:anotherfield,value,...
// The field under validation must be "no", "off", 0, "0", false, or "false" if another field under validation is equal to a specified value.
func constructDeclinedIfRule(cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("declined_if rule requires at least 2 arguments")
	}
	otherField := args[0]
	expectedValues := args[1:]

	return func(ctx *ValidationContext) (bool, error) {
		otherValue, ok := ctx.Raw[otherField]
		if !ok {
			return false, fmt.Errorf("the %s field is not present", otherField)
		}
		match := false
		for _, v := range expectedValues {
			if otherValue == v {
				match = true
				break
			}
		}
		if match {
			val := strings.ToLower(ctx.FieldValue)
			if val == "no" || val == "off" || val == "0" || val == "false" {
				return true, nil
			}
			return false, fmt.Errorf("the %s field must be declined when %s is %s", ctx.FieldName, otherField, strings.Join(expectedValues, ", "))
		}
		return true, nil
	}, nil
}

var embeddedBooleanRules = map[string]RuleConstructor{
	"accepted":    constructAcceptedRule,
	"accepted_if": constructAcceptedIfRule,
	"boolean":     constructBooleanRule,
	"declined":    constructDeclinedRule,
	"declined_if": constructDeclinedIfRule,
}
