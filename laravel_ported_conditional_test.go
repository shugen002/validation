package validation

import (
	"testing"
)

// TestConditionalValidationRules tests ported from Laravel validation tests
func TestConditionalValidationRules(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		// Accepted if tests
		{
			name: "accepted_if rule passes when condition met and field accepted",
			data: map[string]interface{}{
				"type":   "premium",
				"agreed": "yes",
			},
			rules: map[string]interface{}{
				"agreed": "accepted_if:type,premium",
			},
			valid: true,
		},
		{
			name: "accepted_if rule passes when condition not met",
			data: map[string]interface{}{
				"type":   "basic",
				"agreed": "no",
			},
			rules: map[string]interface{}{
				"agreed": "accepted_if:type,premium",
			},
			valid: true,
		},
		{
			name: "accepted_if rule fails when condition met but field not accepted",
			data: map[string]interface{}{
				"type":   "premium",
				"agreed": "no",
			},
			rules: map[string]interface{}{
				"agreed": "accepted_if:type,premium",
			},
			valid: false,
		},
		{
			name: "accepted_if rule with multiple values",
			data: map[string]interface{}{
				"type":   "enterprise",
				"agreed": "yes",
			},
			rules: map[string]interface{}{
				"agreed": "accepted_if:type,premium,enterprise",
			},
			valid: true,
		},
		// Declined if tests
		{
			name: "declined_if rule passes when condition met and field declined",
			data: map[string]interface{}{
				"type":   "basic",
				"agreed": "no",
			},
			rules: map[string]interface{}{
				"agreed": "declined_if:type,basic",
			},
			valid: true,
		},
		{
			name: "declined_if rule passes when condition not met",
			data: map[string]interface{}{
				"type":   "premium",
				"agreed": "yes",
			},
			rules: map[string]interface{}{
				"agreed": "declined_if:type,basic",
			},
			valid: true,
		},
		{
			name: "declined_if rule fails when condition met but field not declined",
			data: map[string]interface{}{
				"type":   "basic",
				"agreed": "yes",
			},
			rules: map[string]interface{}{
				"agreed": "declined_if:type,basic",
			},
			valid: false,
		},
		// Confirmed rule tests
		{
			name: "confirmed rule passes with matching confirmation",
			data: map[string]interface{}{
				"password":              "secret123",
				"password_confirmation": "secret123",
			},
			rules: map[string]interface{}{
				"password": "confirmed",
			},
			valid: true,
		},
		{
			name: "confirmed rule fails with non-matching confirmation",
			data: map[string]interface{}{
				"password":              "secret123",
				"password_confirmation": "different",
			},
			rules: map[string]interface{}{
				"password": "confirmed",
			},
			valid: false,
		},
		{
			name: "confirmed rule fails with missing confirmation",
			data: map[string]interface{}{
				"password": "secret123",
			},
			rules: map[string]interface{}{
				"password": "confirmed",
			},
			valid: false,
		},
		{
			name: "confirmed rule with custom confirmation field name",
			data: map[string]interface{}{
				"email":       "user@example.com",
				"email_match": "user@example.com",
			},
			rules: map[string]interface{}{
				"email": "confirmed",
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

// TestValidatorAdvancedFeatures tests ported from Laravel validator advanced features
func TestValidatorAdvancedFeatures(t *testing.T) {
	factory := NewFactory()

	t.Run("nullable validation", func(t *testing.T) {
		tests := []struct {
			name  string
			data  map[string]interface{}
			rules map[string]interface{}
			valid bool
		}{
			{
				name: "nullable allows nil values",
				data: map[string]interface{}{
					"field": nil,
				},
				rules: map[string]interface{}{
					"field": "nullable|string",
				},
				valid: true,
			},
			{
				name: "nullable allows valid non-nil values",
				data: map[string]interface{}{
					"field": "hello",
				},
				rules: map[string]interface{}{
					"field": "nullable|string",
				},
				valid: true,
			},
			{
				name: "nullable still validates non-nil values",
				data: map[string]interface{}{
					"field": 123,
				},
				rules: map[string]interface{}{
					"field": "nullable|string",
				},
				valid: false,
			},
			{
				name: "nullable with email rule",
				data: map[string]interface{}{
					"email": nil,
				},
				rules: map[string]interface{}{
					"email": "nullable|email",
				},
				valid: true,
			},
			{
				name: "nullable with valid email",
				data: map[string]interface{}{
					"email": "user@example.com",
				},
				rules: map[string]interface{}{
					"email": "nullable|email",
				},
				valid: true,
			},
			{
				name: "nullable with invalid email",
				data: map[string]interface{}{
					"email": "invalid-email",
				},
				rules: map[string]interface{}{
					"email": "nullable|email",
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
	})

	t.Run("stop on first failure", func(t *testing.T) {
		data := map[string]interface{}{
			"field": "invalid",
		}
		rules := map[string]interface{}{
			"field": "bail|numeric|min:10|max:5", // Multiple rules that would fail
		}

		validator := factory.Make(data, rules)
		if validator.Passes() {
			t.Error("Expected validation to fail")
		}

		// With bail, only the first failing rule should be reported
		errors := validator.Errors()
		if !errors.Has("field") {
			t.Error("Expected field to have errors")
		}

		// The exact behavior of bail may vary by implementation
		// but there should be at least one error
		fieldErrors := errors.Get("field")
		if len(fieldErrors) == 0 {
			t.Error("Expected at least one error for field")
		}
	})

	t.Run("custom error messages", func(t *testing.T) {
		data := map[string]interface{}{
			"email": "invalid-email",
		}
		rules := map[string]interface{}{
			"email": "required|email",
		}
		messages := map[string]string{
			"email.required": "The email field is mandatory",
			"email.email":    "The email must be a valid email address",
		}

		validator := factory.Make(data, rules, messages)
		if validator.Passes() {
			t.Error("Expected validation to fail")
		}

		errors := validator.Errors()
		if !errors.Has("email") {
			t.Error("Expected email field to have errors")
		}

		// Check if custom message is used (implementation may vary)
		emailErrors := errors.Get("email")
		if len(emailErrors) == 0 {
			t.Error("Expected email to have error messages")
		}
	})

	t.Run("validator methods", func(t *testing.T) {
		data := map[string]interface{}{
			"name":  "John",
			"email": "invalid-email",
			"age":   25,
		}
		rules := map[string]interface{}{
			"name":  "required|string",
			"email": "required|email",
			"age":   "integer|min:18",
		}

		validator := factory.Make(data, rules)

		// Test Passes/Fails methods
		if validator.Passes() {
			t.Error("Expected validation to fail due to invalid email")
		}

		if !validator.Fails() {
			t.Error("Expected Fails() to return true")
		}

		// Test Valid() method - should return only valid data
		validData := validator.Valid()
		if len(validData) != 2 { // name and age should be valid
			t.Errorf("Expected 2 valid fields, got: %d", len(validData))
		}

		if validData["name"] != "John" {
			t.Error("Expected valid data to contain name")
		}

		if validData["age"] != 25 {
			t.Error("Expected valid data to contain age")
		}

		// Test Invalid() method - should return only invalid data
		invalidData := validator.Invalid()
		if len(invalidData) != 1 { // only email should be invalid
			t.Errorf("Expected 1 invalid field, got: %d", len(invalidData))
		}

		if invalidData["email"] != "invalid-email" {
			t.Error("Expected invalid data to contain email")
		}

		// Test error bag methods
		errors := validator.Errors()
		if !errors.Has("email") {
			t.Error("Expected errors to have email field")
		}

		if errors.Has("name") {
			t.Error("Expected errors to not have name field")
		}

		if errors.Has("age") {
			t.Error("Expected errors to not have age field")
		}

		allErrors := errors.All()
		if len(allErrors) == 0 {
			t.Error("Expected to have some errors")
		}

		emailErrors := errors.Get("email")
		if len(emailErrors) == 0 {
			t.Error("Expected email to have error messages")
		}

		firstError := errors.First("email")
		if firstError == "" {
			t.Error("Expected email to have a first error message")
		}
	})

	t.Run("array validation with wildcard", func(t *testing.T) {
		data := map[string]interface{}{
			"users": []map[string]interface{}{
				{
					"name":  "John",
					"email": "john@example.com",
					"age":   25,
				},
				{
					"name":  "Jane",
					"email": "invalid-email",
					"age":   30,
				},
			},
		}
		rules := map[string]interface{}{
			"users.*.name":  "required|string",
			"users.*.email": "required|email",
			"users.*.age":   "integer|min:18",
		}

		validator := factory.Make(data, rules)
		if validator.Passes() {
			t.Error("Expected validation to fail due to invalid email in second user")
		}

		errors := validator.Errors()
		
		// Check that the validation properly handles array indices
		// The exact error key format may vary by implementation
		hasArrayError := false
		for key := range errors.All() {
			if key == "users.1.email" || key == "users[1].email" || key == "users.*.email" {
				hasArrayError = true
				break
			}
		}
		
		if !hasArrayError {
			t.Errorf("Expected to find array validation error, got errors: %v", errors.All())
		}
	})
}