package validation

import (
	"testing"
)

// TestEmailValidationRule tests ported from Laravel's ValidationEmailRuleTest.php
func TestEmailValidationRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		email interface{}
		valid bool
	}{
		{
			name:  "valid email passes",
			email: "taylor@laravel.com",
			valid: true,
		},
		{
			name:  "valid email with subdomain passes",
			email: "user@mail.example.com",
			valid: true,
		},
		{
			name:  "valid email with plus sign passes",
			email: "user+tag@example.com",
			valid: true,
		},
		{
			name:  "valid email with numbers passes",
			email: "user123@example123.com",
			valid: true,
		},
		{
			name:  "invalid email without domain fails",
			email: "invalid",
			valid: false,
		},
		{
			name:  "invalid email without at symbol fails",
			email: "invalid.email.com",
			valid: false,
		},
		{
			name:  "invalid email with spaces fails",
			email: "invalid @example.com",
			valid: false,
		},
		{
			name:  "invalid email with multiple at symbols fails",
			email: "invalid@@example.com",
			valid: false,
		},
		{
			name:  "number instead of string fails",
			email: 12345,
			valid: false,
		},
		{
			name:  "empty string fails",
			email: "",
			valid: false,
		},
		{
			name:  "nil value passes (nullable)",
			email: nil,
			valid: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := map[string]interface{}{
				"email": test.email,
			}
			rules := map[string]interface{}{
				"email": "email",
			}

			validator := factory.Make(data, rules)

			if test.valid {
				if !validator.Passes() {
					t.Errorf("Expected email validation to pass, but it failed. Errors: %v", validator.Errors().All())
				}
			} else {
				if validator.Passes() {
					t.Error("Expected email validation to fail, but it passed")
				}
			}
		})
	}
}

// TestStringValidationRules tests ported from Laravel validation tests
func TestStringValidationRules(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		// Alpha rule tests
		{
			name: "alpha rule passes with alphabetic characters",
			data: map[string]interface{}{
				"field": "hello",
			},
			rules: map[string]interface{}{
				"field": "alpha",
			},
			valid: true,
		},
		{
			name: "alpha rule passes with mixed case",
			data: map[string]interface{}{
				"field": "HelloWorld",
			},
			rules: map[string]interface{}{
				"field": "alpha",
			},
			valid: true,
		},
		{
			name: "alpha rule fails with numbers",
			data: map[string]interface{}{
				"field": "hello123",
			},
			rules: map[string]interface{}{
				"field": "alpha",
			},
			valid: false,
		},
		{
			name: "alpha rule fails with special characters",
			data: map[string]interface{}{
				"field": "hello@world",
			},
			rules: map[string]interface{}{
				"field": "alpha",
			},
			valid: false,
		},
		// Alpha-numeric rule tests
		{
			name: "alpha_num rule passes with letters and numbers",
			data: map[string]interface{}{
				"field": "hello123",
			},
			rules: map[string]interface{}{
				"field": "alpha_num",
			},
			valid: true,
		},
		{
			name: "alpha_num rule passes with only letters",
			data: map[string]interface{}{
				"field": "hello",
			},
			rules: map[string]interface{}{
				"field": "alpha_num",
			},
			valid: true,
		},
		{
			name: "alpha_num rule passes with only numbers",
			data: map[string]interface{}{
				"field": "123",
			},
			rules: map[string]interface{}{
				"field": "alpha_num",
			},
			valid: true,
		},
		{
			name: "alpha_num rule fails with special characters",
			data: map[string]interface{}{
				"field": "hello@123",
			},
			rules: map[string]interface{}{
				"field": "alpha_num",
			},
			valid: false,
		},
		// Alpha-dash rule tests
		{
			name: "alpha_dash rule passes with letters, numbers, dashes and underscores",
			data: map[string]interface{}{
				"field": "hello_world-123",
			},
			rules: map[string]interface{}{
				"field": "alpha_dash",
			},
			valid: true,
		},
		{
			name: "alpha_dash rule passes with username format",
			data: map[string]interface{}{
				"field": "user_name",
			},
			rules: map[string]interface{}{
				"field": "alpha_dash",
			},
			valid: true,
		},
		{
			name: "alpha_dash rule passes with slug format",
			data: map[string]interface{}{
				"field": "my-blog-post",
			},
			rules: map[string]interface{}{
				"field": "alpha_dash",
			},
			valid: true,
		},
		{
			name: "alpha_dash rule fails with spaces",
			data: map[string]interface{}{
				"field": "hello world",
			},
			rules: map[string]interface{}{
				"field": "alpha_dash",
			},
			valid: false,
		},
		{
			name: "alpha_dash rule fails with special characters",
			data: map[string]interface{}{
				"field": "hello@world",
			},
			rules: map[string]interface{}{
				"field": "alpha_dash",
			},
			valid: false,
		},
		// Regex rule tests
		{
			name: "regex rule passes with matching pattern",
			data: map[string]interface{}{
				"field": "ABC123",
			},
			rules: map[string]interface{}{
				"field": "regex:/^[A-Z]{3}[0-9]{3}$/",
			},
			valid: true,
		},
		{
			name: "regex rule fails with non-matching pattern",
			data: map[string]interface{}{
				"field": "abc123",
			},
			rules: map[string]interface{}{
				"field": "regex:/^[A-Z]{3}[0-9]{3}$/",
			},
			valid: false,
		},
		{
			name: "regex rule passes without delimiters",
			data: map[string]interface{}{
				"field": "ABC123",
			},
			rules: map[string]interface{}{
				"field": "regex:^[A-Z]{3}[0-9]{3}$",
			},
			valid: true,
		},
		// Not regex rule tests
		{
			name: "not_regex rule passes with non-matching pattern",
			data: map[string]interface{}{
				"field": "abc123",
			},
			rules: map[string]interface{}{
				"field": "not_regex:/^[A-Z]{3}[0-9]{3}$/",
			},
			valid: true,
		},
		{
			name: "not_regex rule fails with matching pattern",
			data: map[string]interface{}{
				"field": "ABC123",
			},
			rules: map[string]interface{}{
				"field": "not_regex:/^[A-Z]{3}[0-9]{3}$/",
			},
			valid: false,
		},
		// Starts with rule tests
		{
			name: "starts_with rule passes with matching prefix",
			data: map[string]interface{}{
				"field": "hello world",
			},
			rules: map[string]interface{}{
				"field": "starts_with:hello",
			},
			valid: true,
		},
		{
			name: "starts_with rule passes with one of multiple prefixes",
			data: map[string]interface{}{
				"field": "goodbye world",
			},
			rules: map[string]interface{}{
				"field": "starts_with:hello,goodbye",
			},
			valid: true,
		},
		{
			name: "starts_with rule fails with non-matching prefix",
			data: map[string]interface{}{
				"field": "world hello",
			},
			rules: map[string]interface{}{
				"field": "starts_with:hello",
			},
			valid: false,
		},
		// Ends with rule tests
		{
			name: "ends_with rule passes with matching suffix",
			data: map[string]interface{}{
				"field": "hello world",
			},
			rules: map[string]interface{}{
				"field": "ends_with:world",
			},
			valid: true,
		},
		{
			name: "ends_with rule passes with one of multiple suffixes",
			data: map[string]interface{}{
				"field": "hello universe",
			},
			rules: map[string]interface{}{
				"field": "ends_with:world,universe",
			},
			valid: true,
		},
		{
			name: "ends_with rule fails with non-matching suffix",
			data: map[string]interface{}{
				"field": "world hello",
			},
			rules: map[string]interface{}{
				"field": "ends_with:world",
			},
			valid: false,
		},
		// Doesn't start with rule tests
		{
			name: "doesnt_start_with rule passes with non-matching prefix",
			data: map[string]interface{}{
				"field": "world hello",
			},
			rules: map[string]interface{}{
				"field": "doesnt_start_with:hello",
			},
			valid: true,
		},
		{
			name: "doesnt_start_with rule fails with matching prefix",
			data: map[string]interface{}{
				"field": "hello world",
			},
			rules: map[string]interface{}{
				"field": "doesnt_start_with:hello",
			},
			valid: false,
		},
		// Doesn't end with rule tests
		{
			name: "doesnt_end_with rule passes with non-matching suffix",
			data: map[string]interface{}{
				"field": "hello universe",
			},
			rules: map[string]interface{}{
				"field": "doesnt_end_with:world",
			},
			valid: true,
		},
		{
			name: "doesnt_end_with rule fails with matching suffix",
			data: map[string]interface{}{
				"field": "hello world",
			},
			rules: map[string]interface{}{
				"field": "doesnt_end_with:world",
			},
			valid: false,
		},
		// Case validation tests
		{
			name: "uppercase rule passes with all uppercase",
			data: map[string]interface{}{
				"field": "HELLO WORLD",
			},
			rules: map[string]interface{}{
				"field": "uppercase",
			},
			valid: true,
		},
		{
			name: "uppercase rule fails with mixed case",
			data: map[string]interface{}{
				"field": "Hello World",
			},
			rules: map[string]interface{}{
				"field": "uppercase",
			},
			valid: false,
		},
		{
			name: "lowercase rule passes with all lowercase",
			data: map[string]interface{}{
				"field": "hello world",
			},
			rules: map[string]interface{}{
				"field": "lowercase",
			},
			valid: true,
		},
		{
			name: "lowercase rule fails with mixed case",
			data: map[string]interface{}{
				"field": "Hello World",
			},
			rules: map[string]interface{}{
				"field": "lowercase",
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

// TestAcceptedAndDeclinedRules tests ported from Laravel validation tests
func TestAcceptedAndDeclinedRules(t *testing.T) {
	factory := NewFactory()

	acceptedTests := []struct {
		name  string
		value interface{}
		valid bool
	}{
		{
			name:  "accepts 'yes'",
			value: "yes",
			valid: true,
		},
		{
			name:  "accepts 'on'",
			value: "on",
			valid: true,
		},
		{
			name:  "accepts '1'",
			value: "1",
			valid: true,
		},
		{
			name:  "accepts 1",
			value: 1,
			valid: true,
		},
		{
			name:  "accepts true",
			value: true,
			valid: true,
		},
		{
			name:  "accepts 'true'",
			value: "true",
			valid: true,
		},
		{
			name:  "rejects 'no'",
			value: "no",
			valid: false,
		},
		{
			name:  "rejects 'off'",
			value: "off",
			valid: false,
		},
		{
			name:  "rejects '0'",
			value: "0",
			valid: false,
		},
		{
			name:  "rejects 0",
			value: 0,
			valid: false,
		},
		{
			name:  "rejects false",
			value: false,
			valid: false,
		},
		{
			name:  "rejects 'false'",
			value: "false",
			valid: false,
		},
	}

	for _, test := range acceptedTests {
		t.Run("accepted "+test.name, func(t *testing.T) {
			data := map[string]interface{}{
				"field": test.value,
			}
			rules := map[string]interface{}{
				"field": "accepted",
			}

			validator := factory.Make(data, rules)

			if test.valid {
				if !validator.Passes() {
					t.Errorf("Expected accepted validation to pass, but it failed. Errors: %v", validator.Errors().All())
				}
			} else {
				if validator.Passes() {
					t.Error("Expected accepted validation to fail, but it passed")
				}
			}
		})
	}

	declinedTests := []struct {
		name  string
		value interface{}
		valid bool
	}{
		{
			name:  "accepts 'no'",
			value: "no",
			valid: true,
		},
		{
			name:  "accepts 'off'",
			value: "off",
			valid: true,
		},
		{
			name:  "accepts '0'",
			value: "0",
			valid: true,
		},
		{
			name:  "accepts 0",
			value: 0,
			valid: true,
		},
		{
			name:  "accepts false",
			value: false,
			valid: true,
		},
		{
			name:  "accepts 'false'",
			value: "false",
			valid: true,
		},
		{
			name:  "rejects 'yes'",
			value: "yes",
			valid: false,
		},
		{
			name:  "rejects 'on'",
			value: "on",
			valid: false,
		},
		{
			name:  "rejects '1'",
			value: "1",
			valid: false,
		},
		{
			name:  "rejects 1",
			value: 1,
			valid: false,
		},
		{
			name:  "rejects true",
			value: true,
			valid: false,
		},
		{
			name:  "rejects 'true'",
			value: "true",
			valid: false,
		},
	}

	for _, test := range declinedTests {
		t.Run("declined "+test.name, func(t *testing.T) {
			data := map[string]interface{}{
				"field": test.value,
			}
			rules := map[string]interface{}{
				"field": "declined",
			}

			validator := factory.Make(data, rules)

			if test.valid {
				if !validator.Passes() {
					t.Errorf("Expected declined validation to pass, but it failed. Errors: %v", validator.Errors().All())
				}
			} else {
				if validator.Passes() {
					t.Error("Expected declined validation to fail, but it passed")
				}
			}
		})
	}
}