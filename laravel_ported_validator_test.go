package validation

import (
	"testing"
)

// TestValidatorBasicFunctionality tests ported from Laravel's ValidationValidatorTest.php
func TestValidatorBasicFunctionality(t *testing.T) {
	factory := NewFactory()

	t.Run("basic validation passes", func(t *testing.T) {
		data := map[string]interface{}{
			"name":  "John",
			"email": "john@example.com",
			"age":   25,
		}
		rules := map[string]interface{}{
			"name":  "required",
			"email": "required|email",
			"age":   "integer|min:18",
		}

		validator := factory.Make(data, rules)
		if !validator.Passes() {
			t.Errorf("Expected validation to pass, but it failed. Errors: %v", validator.Errors().All())
		}
	})

	t.Run("basic validation fails", func(t *testing.T) {
		data := map[string]interface{}{
			"name":  "",
			"email": "invalid-email",
			"age":   15,
		}
		rules := map[string]interface{}{
			"name":  "required",
			"email": "required|email",
			"age":   "integer|min:18",
		}

		validator := factory.Make(data, rules)
		if validator.Passes() {
			t.Error("Expected validation to fail, but it passed")
		}
		
		// Check that we have errors for all three fields
		errors := validator.Errors()
		if !errors.Has("name") {
			t.Error("Expected error for name field")
		}
		if !errors.Has("email") {
			t.Error("Expected error for email field")
		}
		if !errors.Has("age") {
			t.Error("Expected error for age field")
		}
	})

	t.Run("nested validation", func(t *testing.T) {
		data := map[string]interface{}{
			"user": map[string]interface{}{
				"name":    "John",
				"email":   "john@example.com",
				"profile": map[string]interface{}{
					"age": 25,
				},
			},
		}
		rules := map[string]interface{}{
			"user.name":         "required",
			"user.email":        "required|email",
			"user.profile.age":  "integer|min:18",
		}

		validator := factory.Make(data, rules)
		if !validator.Passes() {
			t.Errorf("Expected nested validation to pass, but it failed. Errors: %v", validator.Errors().All())
		}
	})

	t.Run("array validation", func(t *testing.T) {
		data := map[string]interface{}{
			"users": []map[string]interface{}{
				{
					"name":  "John",
					"email": "john@example.com",
				},
				{
					"name":  "Jane",
					"email": "jane@example.com",
				},
			},
		}
		rules := map[string]interface{}{
			"users.*.name":  "required",
			"users.*.email": "required|email",
		}

		validator := factory.Make(data, rules)
		if !validator.Passes() {
			t.Errorf("Expected array validation to pass, but it failed. Errors: %v", validator.Errors().All())
		}
	})
}

// TestRequiredRule tests ported from Laravel validation tests
func TestRequiredRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "required passes with non-empty string",
			data: map[string]interface{}{
				"field": "value",
			},
			rules: map[string]interface{}{
				"field": "required",
			},
			valid: true,
		},
		{
			name: "required fails with empty string",
			data: map[string]interface{}{
				"field": "",
			},
			rules: map[string]interface{}{
				"field": "required",
			},
			valid: false,
		},
		{
			name: "required fails with nil",
			data: map[string]interface{}{
				"field": nil,
			},
			rules: map[string]interface{}{
				"field": "required",
			},
			valid: false,
		},
		{
			name: "required fails with missing field",
			data: map[string]interface{}{},
			rules: map[string]interface{}{
				"field": "required",
			},
			valid: false,
		},
		{
			name: "required passes with zero integer",
			data: map[string]interface{}{
				"field": 0,
			},
			rules: map[string]interface{}{
				"field": "required",
			},
			valid: true,
		},
		{
			name: "required passes with false boolean",
			data: map[string]interface{}{
				"field": false,
			},
			rules: map[string]interface{}{
				"field": "required",
			},
			valid: true,
		},
		{
			name: "required passes with empty array",
			data: map[string]interface{}{
				"field": []string{},
			},
			rules: map[string]interface{}{
				"field": "required",
			},
			valid: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validator := factory.Make(test.data, test.rules)

			if test.valid {
				if !validator.Passes() {
					t.Errorf("Expected validation to pass, but it failed. Errors: %v", validator.Errors().All())
				}
			} else {
				if validator.Passes() {
					t.Error("Expected validation to fail, but it passed")
				}
			}
		})
	}
}

// TestTypeValidationRules tests ported from Laravel validation tests
func TestTypeValidationRules(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		// String rule tests
		{
			name: "string rule passes with string",
			data: map[string]interface{}{
				"field": "hello",
			},
			rules: map[string]interface{}{
				"field": "string",
			},
			valid: true,
		},
		{
			name: "string rule fails with integer",
			data: map[string]interface{}{
				"field": 123,
			},
			rules: map[string]interface{}{
				"field": "string",
			},
			valid: false,
		},
		// Integer rule tests
		{
			name: "integer rule passes with integer",
			data: map[string]interface{}{
				"field": 123,
			},
			rules: map[string]interface{}{
				"field": "integer",
			},
			valid: true,
		},
		{
			name: "integer rule passes with numeric string",
			data: map[string]interface{}{
				"field": "123",
			},
			rules: map[string]interface{}{
				"field": "integer",
			},
			valid: true,
		},
		{
			name: "integer rule fails with non-numeric string",
			data: map[string]interface{}{
				"field": "hello",
			},
			rules: map[string]interface{}{
				"field": "integer",
			},
			valid: false,
		},
		// Numeric rule tests
		{
			name: "numeric rule passes with integer",
			data: map[string]interface{}{
				"field": 123,
			},
			rules: map[string]interface{}{
				"field": "numeric",
			},
			valid: true,
		},
		{
			name: "numeric rule passes with float",
			data: map[string]interface{}{
				"field": 123.45,
			},
			rules: map[string]interface{}{
				"field": "numeric",
			},
			valid: true,
		},
		{
			name: "numeric rule passes with numeric string",
			data: map[string]interface{}{
				"field": "123.45",
			},
			rules: map[string]interface{}{
				"field": "numeric",
			},
			valid: true,
		},
		{
			name: "numeric rule fails with non-numeric string",
			data: map[string]interface{}{
				"field": "hello",
			},
			rules: map[string]interface{}{
				"field": "numeric",
			},
			valid: false,
		},
		// Boolean rule tests
		{
			name: "boolean rule passes with true",
			data: map[string]interface{}{
				"field": true,
			},
			rules: map[string]interface{}{
				"field": "boolean",
			},
			valid: true,
		},
		{
			name: "boolean rule passes with false",
			data: map[string]interface{}{
				"field": false,
			},
			rules: map[string]interface{}{
				"field": "boolean",
			},
			valid: true,
		},
		{
			name: "boolean rule passes with 1",
			data: map[string]interface{}{
				"field": 1,
			},
			rules: map[string]interface{}{
				"field": "boolean",
			},
			valid: true,
		},
		{
			name: "boolean rule passes with 0",
			data: map[string]interface{}{
				"field": 0,
			},
			rules: map[string]interface{}{
				"field": "boolean",
			},
			valid: true,
		},
		{
			name: "boolean rule passes with string 'true'",
			data: map[string]interface{}{
				"field": "true",
			},
			rules: map[string]interface{}{
				"field": "boolean",
			},
			valid: true,
		},
		{
			name: "boolean rule passes with string 'false'",
			data: map[string]interface{}{
				"field": "false",
			},
			rules: map[string]interface{}{
				"field": "boolean",
			},
			valid: true,
		},
		{
			name: "boolean rule fails with string 'hello'",
			data: map[string]interface{}{
				"field": "hello",
			},
			rules: map[string]interface{}{
				"field": "boolean",
			},
			valid: false,
		},
		// Array rule tests
		{
			name: "array rule passes with array",
			data: map[string]interface{}{
				"field": []string{"item1", "item2"},
			},
			rules: map[string]interface{}{
				"field": "array",
			},
			valid: true,
		},
		{
			name: "array rule passes with empty array",
			data: map[string]interface{}{
				"field": []string{},
			},
			rules: map[string]interface{}{
				"field": "array",
			},
			valid: true,
		},
		{
			name: "array rule fails with string",
			data: map[string]interface{}{
				"field": "hello",
			},
			rules: map[string]interface{}{
				"field": "array",
			},
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validator := factory.Make(test.data, test.rules)

			if test.valid {
				if !validator.Passes() {
					t.Errorf("Expected validation to pass, but it failed. Errors: %v", validator.Errors().All())
				}
			} else {
				if validator.Passes() {
					t.Error("Expected validation to fail, but it passed")
				}
			}
		})
	}
}