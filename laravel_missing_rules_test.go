package validation

import (
	"testing"
)

func TestDigitsRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "digits rule passes with exact digit count",
			data: map[string]interface{}{
				"phone": "123456",
			},
			rules: map[string]interface{}{
				"phone": "digits:6",
			},
			valid: true,
		},
		{
			name: "digits rule fails with fewer digits",
			data: map[string]interface{}{
				"phone": "12345",
			},
			rules: map[string]interface{}{
				"phone": "digits:6",
			},
			valid: false,
		},
		{
			name: "digits rule fails with more digits",
			data: map[string]interface{}{
				"phone": "1234567",
			},
			rules: map[string]interface{}{
				"phone": "digits:6",
			},
			valid: false,
		},
		{
			name: "digits rule fails with non-digits",
			data: map[string]interface{}{
				"phone": "12a456",
			},
			rules: map[string]interface{}{
				"phone": "digits:6",
			},
			valid: false,
		},
		{
			name: "digits rule passes with integer",
			data: map[string]interface{}{
				"number": 123456,
			},
			rules: map[string]interface{}{
				"number": "digits:6",
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
					t.Errorf("Expected validation to fail, but it passed")
				}
			}
		})
	}
}

func TestDigitsBetweenRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "digits_between rule passes with min digits",
			data: map[string]interface{}{
				"code": "123",
			},
			rules: map[string]interface{}{
				"code": "digits_between:3,6",
			},
			valid: true,
		},
		{
			name: "digits_between rule passes with max digits",
			data: map[string]interface{}{
				"code": "123456",
			},
			rules: map[string]interface{}{
				"code": "digits_between:3,6",
			},
			valid: true,
		},
		{
			name: "digits_between rule passes within range",
			data: map[string]interface{}{
				"code": "1234",
			},
			rules: map[string]interface{}{
				"code": "digits_between:3,6",
			},
			valid: true,
		},
		{
			name: "digits_between rule fails below min",
			data: map[string]interface{}{
				"code": "12",
			},
			rules: map[string]interface{}{
				"code": "digits_between:3,6",
			},
			valid: false,
		},
		{
			name: "digits_between rule fails above max",
			data: map[string]interface{}{
				"code": "1234567",
			},
			rules: map[string]interface{}{
				"code": "digits_between:3,6",
			},
			valid: false,
		},
		{
			name: "digits_between rule fails with non-digits",
			data: map[string]interface{}{
				"code": "123a",
			},
			rules: map[string]interface{}{
				"code": "digits_between:3,6",
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
					t.Errorf("Expected validation to fail, but it passed")
				}
			}
		})
	}
}

func TestFilledRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "filled rule passes with non-empty string",
			data: map[string]interface{}{
				"name": "John",
			},
			rules: map[string]interface{}{
				"name": "filled",
			},
			valid: true,
		},
		{
			name: "filled rule passes when field is absent",
			data: map[string]interface{}{},
			rules: map[string]interface{}{
				"name": "filled",
			},
			valid: true,
		},
		{
			name: "filled rule fails with empty string",
			data: map[string]interface{}{
				"name": "",
			},
			rules: map[string]interface{}{
				"name": "filled",
			},
			valid: false,
		},
		{
			name: "filled rule fails with whitespace only",
			data: map[string]interface{}{
				"name": "   ",
			},
			rules: map[string]interface{}{
				"name": "filled",
			},
			valid: false,
		},
		{
			name: "filled rule passes with non-empty array",
			data: map[string]interface{}{
				"items": []string{"item1"},
			},
			rules: map[string]interface{}{
				"items": "filled",
			},
			valid: true,
		},
		{
			name: "filled rule fails with empty array",
			data: map[string]interface{}{
				"items": []string{},
			},
			rules: map[string]interface{}{
				"items": "filled",
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
					t.Errorf("Expected validation to fail, but it passed")
				}
			}
		})
	}
}

func TestPresentRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "present rule passes when field exists with value",
			data: map[string]interface{}{
				"name": "John",
			},
			rules: map[string]interface{}{
				"name": "present",
			},
			valid: true,
		},
		{
			name: "present rule passes when field exists but empty",
			data: map[string]interface{}{
				"name": "",
			},
			rules: map[string]interface{}{
				"name": "present",
			},
			valid: true,
		},
		{
			name: "present rule fails when field is missing",
			data: map[string]interface{}{},
			rules: map[string]interface{}{
				"name": "present",
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
					t.Errorf("Expected validation to fail, but it passed")
				}
			}
		})
	}
}

func TestProhibitedRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "prohibited rule passes when field is absent",
			data: map[string]interface{}{},
			rules: map[string]interface{}{
				"secret": "prohibited",
			},
			valid: true,
		},
		{
			name: "prohibited rule passes when field is empty string",
			data: map[string]interface{}{
				"secret": "",
			},
			rules: map[string]interface{}{
				"secret": "prohibited",
			},
			valid: true,
		},
		{
			name: "prohibited rule passes when field is whitespace only",
			data: map[string]interface{}{
				"secret": "   ",
			},
			rules: map[string]interface{}{
				"secret": "prohibited",
			},
			valid: true,
		},
		{
			name: "prohibited rule fails when field has value",
			data: map[string]interface{}{
				"secret": "value",
			},
			rules: map[string]interface{}{
				"secret": "prohibited",
			},
			valid: false,
		},
		{
			name: "prohibited rule passes when array is empty",
			data: map[string]interface{}{
				"items": []string{},
			},
			rules: map[string]interface{}{
				"items": "prohibited",
			},
			valid: true,
		},
		{
			name: "prohibited rule fails when array has items",
			data: map[string]interface{}{
				"items": []string{"item1"},
			},
			rules: map[string]interface{}{
				"items": "prohibited",
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
					t.Errorf("Expected validation to fail, but it passed")
				}
			}
		})
	}
}

func TestRequiredIfRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "required_if rule passes when condition field equals value and target field is present",
			data: map[string]interface{}{
				"payment_method": "card",
				"card_number":    "1234567890123456",
			},
			rules: map[string]interface{}{
				"card_number": "required_if:payment_method,card",
			},
			valid: true,
		},
		{
			name: "required_if rule fails when condition field equals value but target field is missing",
			data: map[string]interface{}{
				"payment_method": "card",
			},
			rules: map[string]interface{}{
				"card_number": "required_if:payment_method,card",
			},
			valid: false,
		},
		{
			name: "required_if rule fails when condition field equals value but target field is empty",
			data: map[string]interface{}{
				"payment_method": "card",
				"card_number":    "",
			},
			rules: map[string]interface{}{
				"card_number": "required_if:payment_method,card",
			},
			valid: false,
		},
		{
			name: "required_if rule passes when condition field does not equal value",
			data: map[string]interface{}{
				"payment_method": "paypal",
			},
			rules: map[string]interface{}{
				"card_number": "required_if:payment_method,card",
			},
			valid: true,
		},
		{
			name: "required_if rule passes when condition field is missing",
			data: map[string]interface{}{},
			rules: map[string]interface{}{
				"card_number": "required_if:payment_method,card",
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
					t.Errorf("Expected validation to fail, but it passed")
				}
			}
		})
	}
}

func TestRequiredUnlessRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "required_unless rule passes when condition field equals exception value",
			data: map[string]interface{}{
				"account_type": "guest",
			},
			rules: map[string]interface{}{
				"email": "required_unless:account_type,guest",
			},
			valid: true,
		},
		{
			name: "required_unless rule passes when condition field equals exception value and target field present",
			data: map[string]interface{}{
				"account_type": "guest",
				"email":        "test@example.com",
			},
			rules: map[string]interface{}{
				"email": "required_unless:account_type,guest",
			},
			valid: true,
		},
		{
			name: "required_unless rule passes when condition field does not equal exception value and target field present",
			data: map[string]interface{}{
				"account_type": "user",
				"email":        "test@example.com",
			},
			rules: map[string]interface{}{
				"email": "required_unless:account_type,guest",
			},
			valid: true,
		},
		{
			name: "required_unless rule fails when condition field does not equal exception value and target field missing",
			data: map[string]interface{}{
				"account_type": "user",
			},
			rules: map[string]interface{}{
				"email": "required_unless:account_type,guest",
			},
			valid: false,
		},
		{
			name: "required_unless rule fails when condition field does not equal exception value and target field empty",
			data: map[string]interface{}{
				"account_type": "user",
				"email":        "",
			},
			rules: map[string]interface{}{
				"email": "required_unless:account_type,guest",
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
					t.Errorf("Expected validation to fail, but it passed")
				}
			}
		})
	}
}

func TestRequiredWithRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "required_with rule passes when no dependency fields are present",
			data: map[string]interface{}{},
			rules: map[string]interface{}{
				"password_confirmation": "required_with:password",
			},
			valid: true,
		},
		{
			name: "required_with rule passes when dependency field is present and target field is present",
			data: map[string]interface{}{
				"password":              "secret123",
				"password_confirmation": "secret123",
			},
			rules: map[string]interface{}{
				"password_confirmation": "required_with:password",
			},
			valid: true,
		},
		{
			name: "required_with rule fails when dependency field is present but target field is missing",
			data: map[string]interface{}{
				"password": "secret123",
			},
			rules: map[string]interface{}{
				"password_confirmation": "required_with:password",
			},
			valid: false,
		},
		{
			name: "required_with rule fails when dependency field is present but target field is empty",
			data: map[string]interface{}{
				"password":              "secret123",
				"password_confirmation": "",
			},
			rules: map[string]interface{}{
				"password_confirmation": "required_with:password",
			},
			valid: false,
		},
		{
			name: "required_with rule passes with multiple dependencies when none are present",
			data: map[string]interface{}{},
			rules: map[string]interface{}{
				"phone": "required_with:name,email",
			},
			valid: true,
		},
		{
			name: "required_with rule passes with multiple dependencies when one is present and target is present",
			data: map[string]interface{}{
				"name":  "John",
				"phone": "123456789",
			},
			rules: map[string]interface{}{
				"phone": "required_with:name,email",
			},
			valid: true,
		},
		{
			name: "required_with rule fails with multiple dependencies when one is present but target is missing",
			data: map[string]interface{}{
				"name": "John",
			},
			rules: map[string]interface{}{
				"phone": "required_with:name,email",
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
					t.Errorf("Expected validation to fail, but it passed")
				}
			}
		})
	}
}

func TestRequiredWithoutRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "required_without rule passes when all dependency fields are present",
			data: map[string]interface{}{
				"username": "john_doe",
			},
			rules: map[string]interface{}{
				"email": "required_without:username",
			},
			valid: true,
		},
		{
			name: "required_without rule passes when dependency field is missing and target field is present",
			data: map[string]interface{}{
				"email": "john@example.com",
			},
			rules: map[string]interface{}{
				"email": "required_without:username",
			},
			valid: true,
		},
		{
			name: "required_without rule fails when dependency field is missing and target field is missing",
			data: map[string]interface{}{},
			rules: map[string]interface{}{
				"email": "required_without:username",
			},
			valid: false,
		},
		{
			name: "required_without rule fails when dependency field is missing and target field is empty",
			data: map[string]interface{}{
				"email": "",
			},
			rules: map[string]interface{}{
				"email": "required_without:username",
			},
			valid: false,
		},
		{
			name: "required_without rule passes with multiple dependencies when all are present",
			data: map[string]interface{}{
				"username": "john_doe",
				"phone":    "123456789",
			},
			rules: map[string]interface{}{
				"email": "required_without:username,phone",
			},
			valid: true,
		},
		{
			name: "required_without rule passes with multiple dependencies when one is missing and target is present",
			data: map[string]interface{}{
				"username": "john_doe",
				"email":    "john@example.com",
			},
			rules: map[string]interface{}{
				"email": "required_without:username,phone",
			},
			valid: true,
		},
		{
			name: "required_without rule fails with multiple dependencies when one is missing and target is missing",
			data: map[string]interface{}{
				"username": "john_doe",
			},
			rules: map[string]interface{}{
				"email": "required_without:username,phone",
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
					t.Errorf("Expected validation to fail, but it passed")
				}
			}
		})
	}
}