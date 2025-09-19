package validation

import (
	"testing"
)

func TestDecimalRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "decimal rule passes with exact decimal places",
			data: map[string]interface{}{
				"price": "9.99",
			},
			rules: map[string]interface{}{
				"price": "decimal:2",
			},
			valid: true,
		},
		{
			name: "decimal rule fails with wrong decimal places",
			data: map[string]interface{}{
				"price": "9.9",
			},
			rules: map[string]interface{}{
				"price": "decimal:2",
			},
			valid: false,
		},
		{
			name: "decimal rule passes with range",
			data: map[string]interface{}{
				"price": "9.123",
			},
			rules: map[string]interface{}{
				"price": "decimal:2,4",
			},
			valid: true,
		},
		{
			name: "decimal rule fails outside range",
			data: map[string]interface{}{
				"price": "9.12345",
			},
			rules: map[string]interface{}{
				"price": "decimal:2,4",
			},
			valid: false,
		},
		{
			name: "decimal rule passes with integer when range includes 0",
			data: map[string]interface{}{
				"price": "9",
			},
			rules: map[string]interface{}{
				"price": "decimal:0,2",
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

func TestDistinctRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "distinct rule passes with unique values",
			data: map[string]interface{}{
				"items": []interface{}{"a", "b", "c"},
			},
			rules: map[string]interface{}{
				"items": "distinct",
			},
			valid: true,
		},
		{
			name: "distinct rule fails with duplicate values",
			data: map[string]interface{}{
				"items": []interface{}{"a", "b", "a"},
			},
			rules: map[string]interface{}{
				"items": "distinct",
			},
			valid: false,
		},
		{
			name: "distinct rule with strict mode",
			data: map[string]interface{}{
				"items": []interface{}{"1", 1},
			},
			rules: map[string]interface{}{
				"items": "distinct:strict",
			},
			valid: true,
		},
		{
			name: "distinct rule without strict mode treats different types as same",
			data: map[string]interface{}{
				"items": []interface{}{"1", 1},
			},
			rules: map[string]interface{}{
				"items": "distinct",
			},
			valid: false,
		},
		{
			name: "distinct rule with ignore_case",
			data: map[string]interface{}{
				"items": []interface{}{"A", "a"},
			},
			rules: map[string]interface{}{
				"items": "distinct:ignore_case",
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

func TestMinDigitsRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "min_digits rule passes with enough digits",
			data: map[string]interface{}{
				"code": "12345",
			},
			rules: map[string]interface{}{
				"code": "min_digits:3",
			},
			valid: true,
		},
		{
			name: "min_digits rule fails with too few digits",
			data: map[string]interface{}{
				"code": "12",
			},
			rules: map[string]interface{}{
				"code": "min_digits:3",
			},
			valid: false,
		},
		{
			name: "min_digits rule fails with non-numeric",
			data: map[string]interface{}{
				"code": "abc",
			},
			rules: map[string]interface{}{
				"code": "min_digits:3",
			},
			valid: false,
		},
		{
			name: "min_digits rule passes with negative number",
			data: map[string]interface{}{
				"code": "-12345",
			},
			rules: map[string]interface{}{
				"code": "min_digits:3",
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

func TestMultipleOfRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "multiple_of rule passes with exact multiple",
			data: map[string]interface{}{
				"quantity": "10",
			},
			rules: map[string]interface{}{
				"quantity": "multiple_of:5",
			},
			valid: true,
		},
		{
			name: "multiple_of rule fails with non-multiple",
			data: map[string]interface{}{
				"quantity": "7",
			},
			rules: map[string]interface{}{
				"quantity": "multiple_of:5",
			},
			valid: false,
		},
		{
			name: "multiple_of rule passes with decimal",
			data: map[string]interface{}{
				"price": "1.5",
			},
			rules: map[string]interface{}{
				"price": "multiple_of:0.5",
			},
			valid: true,
		},
		{
			name: "multiple_of rule fails with decimal",
			data: map[string]interface{}{
				"price": "1.7",
			},
			rules: map[string]interface{}{
				"price": "multiple_of:0.5",
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

func TestMissingRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "missing rule passes when field is absent",
			data: map[string]interface{}{},
			rules: map[string]interface{}{
				"secret": "missing",
			},
			valid: true,
		},
		{
			name: "missing rule fails when field is present",
			data: map[string]interface{}{
				"secret": "value",
			},
			rules: map[string]interface{}{
				"secret": "missing",
			},
			valid: false,
		},
		{
			name: "missing rule fails when field is present but empty",
			data: map[string]interface{}{
				"secret": "",
			},
			rules: map[string]interface{}{
				"secret": "missing",
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

func TestMissingIfRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "missing_if rule passes when condition not met",
			data: map[string]interface{}{
				"type":   "public",
				"secret": "value",
			},
			rules: map[string]interface{}{
				"secret": "missing_if:type,private",
			},
			valid: true,
		},
		{
			name: "missing_if rule passes when condition met and field missing",
			data: map[string]interface{}{
				"type": "private",
			},
			rules: map[string]interface{}{
				"secret": "missing_if:type,private",
			},
			valid: true,
		},
		{
			name: "missing_if rule fails when condition met and field present",
			data: map[string]interface{}{
				"type":   "private",
				"secret": "value",
			},
			rules: map[string]interface{}{
				"secret": "missing_if:type,private",
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

func TestRequiredArrayKeysRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "required_array_keys rule passes with all required keys",
			data: map[string]interface{}{
				"config": map[string]interface{}{
					"name":    "test",
					"version": "1.0",
					"author":  "user",
				},
			},
			rules: map[string]interface{}{
				"config": "required_array_keys:name,version",
			},
			valid: true,
		},
		{
			name: "required_array_keys rule fails with missing key",
			data: map[string]interface{}{
				"config": map[string]interface{}{
					"name": "test",
				},
			},
			rules: map[string]interface{}{
				"config": "required_array_keys:name,version",
			},
			valid: false,
		},
		{
			name: "required_array_keys rule fails with non-array",
			data: map[string]interface{}{
				"config": "not an array",
			},
			rules: map[string]interface{}{
				"config": "required_array_keys:name,version",
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