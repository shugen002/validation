package validation

import (
	"testing"
)

// TestNumericValidationRules tests ported from Laravel's ValidationNumericRuleTest.php
func TestNumericValidationRules(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		// Min rule tests
		{
			name: "min rule passes with value above minimum",
			data: map[string]interface{}{
				"field": 10,
			},
			rules: map[string]interface{}{
				"field": "min:5",
			},
			valid: true,
		},
		{
			name: "min rule passes with value equal to minimum",
			data: map[string]interface{}{
				"field": 5,
			},
			rules: map[string]interface{}{
				"field": "min:5",
			},
			valid: true,
		},
		{
			name: "min rule fails with value below minimum",
			data: map[string]interface{}{
				"field": 3,
			},
			rules: map[string]interface{}{
				"field": "min:5",
			},
			valid: false,
		},
		{
			name: "min rule with string checks length",
			data: map[string]interface{}{
				"field": "hello world",
			},
			rules: map[string]interface{}{
				"field": "min:5",
			},
			valid: true,
		},
		{
			name: "min rule with short string fails",
			data: map[string]interface{}{
				"field": "hi",
			},
			rules: map[string]interface{}{
				"field": "min:5",
			},
			valid: false,
		},
		// Max rule tests
		{
			name: "max rule passes with value below maximum",
			data: map[string]interface{}{
				"field": 5,
			},
			rules: map[string]interface{}{
				"field": "max:10",
			},
			valid: true,
		},
		{
			name: "max rule passes with value equal to maximum",
			data: map[string]interface{}{
				"field": 10,
			},
			rules: map[string]interface{}{
				"field": "max:10",
			},
			valid: true,
		},
		{
			name: "max rule fails with value above maximum",
			data: map[string]interface{}{
				"field": 15,
			},
			rules: map[string]interface{}{
				"field": "max:10",
			},
			valid: false,
		},
		{
			name: "max rule with string checks length",
			data: map[string]interface{}{
				"field": "hello",
			},
			rules: map[string]interface{}{
				"field": "max:10",
			},
			valid: true,
		},
		{
			name: "max rule with long string fails",
			data: map[string]interface{}{
				"field": "hello world this is a very long string",
			},
			rules: map[string]interface{}{
				"field": "max:10",
			},
			valid: false,
		},
		// Between rule tests
		{
			name: "between rule passes with value in range",
			data: map[string]interface{}{
				"field": 7,
			},
			rules: map[string]interface{}{
				"field": "between:5,10",
			},
			valid: true,
		},
		{
			name: "between rule passes with value at lower bound",
			data: map[string]interface{}{
				"field": 5,
			},
			rules: map[string]interface{}{
				"field": "between:5,10",
			},
			valid: true,
		},
		{
			name: "between rule passes with value at upper bound",
			data: map[string]interface{}{
				"field": 10,
			},
			rules: map[string]interface{}{
				"field": "between:5,10",
			},
			valid: true,
		},
		{
			name: "between rule fails with value below range",
			data: map[string]interface{}{
				"field": 3,
			},
			rules: map[string]interface{}{
				"field": "between:5,10",
			},
			valid: false,
		},
		{
			name: "between rule fails with value above range",
			data: map[string]interface{}{
				"field": 15,
			},
			rules: map[string]interface{}{
				"field": "between:5,10",
			},
			valid: false,
		},
		{
			name: "between rule with string checks length",
			data: map[string]interface{}{
				"field": "hello",
			},
			rules: map[string]interface{}{
				"field": "between:3,10",
			},
			valid: true,
		},
		{
			name: "between rule with numeric string",
			data: map[string]interface{}{
				"field": "7",
			},
			rules: map[string]interface{}{
				"field": "numeric|between:5,10",
			},
			valid: true,
		},
		// Size rule tests
		{
			name: "size rule passes with exact numeric value",
			data: map[string]interface{}{
				"field": 10,
			},
			rules: map[string]interface{}{
				"field": "size:10",
			},
			valid: true,
		},
		{
			name: "size rule fails with different numeric value",
			data: map[string]interface{}{
				"field": 5,
			},
			rules: map[string]interface{}{
				"field": "size:10",
			},
			valid: false,
		},
		{
			name: "size rule passes with exact string length",
			data: map[string]interface{}{
				"field": "hello",
			},
			rules: map[string]interface{}{
				"field": "size:5",
			},
			valid: true,
		},
		{
			name: "size rule fails with different string length",
			data: map[string]interface{}{
				"field": "hello world",
			},
			rules: map[string]interface{}{
				"field": "size:5",
			},
			valid: false,
		},
		// Greater than rule tests
		{
			name: "gt rule passes with value greater than reference",
			data: map[string]interface{}{
				"field": 10,
			},
			rules: map[string]interface{}{
				"field": "gt:5",
			},
			valid: true,
		},
		{
			name: "gt rule fails with value equal to reference",
			data: map[string]interface{}{
				"field": 5,
			},
			rules: map[string]interface{}{
				"field": "gt:5",
			},
			valid: false,
		},
		{
			name: "gt rule fails with value less than reference",
			data: map[string]interface{}{
				"field": 3,
			},
			rules: map[string]interface{}{
				"field": "gt:5",
			},
			valid: false,
		},
		// Greater than or equal rule tests
		{
			name: "gte rule passes with value greater than reference",
			data: map[string]interface{}{
				"field": 10,
			},
			rules: map[string]interface{}{
				"field": "gte:5",
			},
			valid: true,
		},
		{
			name: "gte rule passes with value equal to reference",
			data: map[string]interface{}{
				"field": 5,
			},
			rules: map[string]interface{}{
				"field": "gte:5",
			},
			valid: true,
		},
		{
			name: "gte rule fails with value less than reference",
			data: map[string]interface{}{
				"field": 3,
			},
			rules: map[string]interface{}{
				"field": "gte:5",
			},
			valid: false,
		},
		// Less than rule tests
		{
			name: "lt rule passes with value less than reference",
			data: map[string]interface{}{
				"field": 3,
			},
			rules: map[string]interface{}{
				"field": "lt:5",
			},
			valid: true,
		},
		{
			name: "lt rule fails with value equal to reference",
			data: map[string]interface{}{
				"field": 5,
			},
			rules: map[string]interface{}{
				"field": "lt:5",
			},
			valid: false,
		},
		{
			name: "lt rule fails with value greater than reference",
			data: map[string]interface{}{
				"field": 10,
			},
			rules: map[string]interface{}{
				"field": "lt:5",
			},
			valid: false,
		},
		// Less than or equal rule tests
		{
			name: "lte rule passes with value less than reference",
			data: map[string]interface{}{
				"field": 3,
			},
			rules: map[string]interface{}{
				"field": "lte:5",
			},
			valid: true,
		},
		{
			name: "lte rule passes with value equal to reference",
			data: map[string]interface{}{
				"field": 5,
			},
			rules: map[string]interface{}{
				"field": "lte:5",
			},
			valid: true,
		},
		{
			name: "lte rule fails with value greater than reference",
			data: map[string]interface{}{
				"field": 10,
			},
			rules: map[string]interface{}{
				"field": "lte:5",
			},
			valid: false,
		},
		// Float/decimal tests
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
			name: "numeric rule passes with negative number",
			data: map[string]interface{}{
				"field": -123,
			},
			rules: map[string]interface{}{
				"field": "numeric",
			},
			valid: true,
		},
		{
			name: "min rule with float",
			data: map[string]interface{}{
				"field": 5.5,
			},
			rules: map[string]interface{}{
				"field": "numeric|min:5.0",
			},
			valid: true,
		},
		{
			name: "max rule with float",
			data: map[string]interface{}{
				"field": 4.5,
			},
			rules: map[string]interface{}{
				"field": "numeric|max:5.0",
			},
			valid: true,
		},
		{
			name: "between rule with floats",
			data: map[string]interface{}{
				"field": 7.5,
			},
			rules: map[string]interface{}{
				"field": "numeric|between:5.0,10.0",
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

// TestFieldComparisonRules tests ported from Laravel validation tests
func TestFieldComparisonRules(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		// Greater than field tests
		{
			name: "gt field rule passes when field is greater",
			data: map[string]interface{}{
				"field1": 10,
				"field2": 5,
			},
			rules: map[string]interface{}{
				"field1": "gt:field2",
			},
			valid: true,
		},
		{
			name: "gt field rule fails when field is equal",
			data: map[string]interface{}{
				"field1": 5,
				"field2": 5,
			},
			rules: map[string]interface{}{
				"field1": "gt:field2",
			},
			valid: false,
		},
		{
			name: "gt field rule fails when field is less",
			data: map[string]interface{}{
				"field1": 3,
				"field2": 5,
			},
			rules: map[string]interface{}{
				"field1": "gt:field2",
			},
			valid: false,
		},
		// Greater than or equal field tests
		{
			name: "gte field rule passes when field is greater",
			data: map[string]interface{}{
				"field1": 10,
				"field2": 5,
			},
			rules: map[string]interface{}{
				"field1": "gte:field2",
			},
			valid: true,
		},
		{
			name: "gte field rule passes when field is equal",
			data: map[string]interface{}{
				"field1": 5,
				"field2": 5,
			},
			rules: map[string]interface{}{
				"field1": "gte:field2",
			},
			valid: true,
		},
		{
			name: "gte field rule fails when field is less",
			data: map[string]interface{}{
				"field1": 3,
				"field2": 5,
			},
			rules: map[string]interface{}{
				"field1": "gte:field2",
			},
			valid: false,
		},
		// Less than field tests
		{
			name: "lt field rule passes when field is less",
			data: map[string]interface{}{
				"field1": 3,
				"field2": 5,
			},
			rules: map[string]interface{}{
				"field1": "lt:field2",
			},
			valid: true,
		},
		{
			name: "lt field rule fails when field is equal",
			data: map[string]interface{}{
				"field1": 5,
				"field2": 5,
			},
			rules: map[string]interface{}{
				"field1": "lt:field2",
			},
			valid: false,
		},
		{
			name: "lt field rule fails when field is greater",
			data: map[string]interface{}{
				"field1": 10,
				"field2": 5,
			},
			rules: map[string]interface{}{
				"field1": "lt:field2",
			},
			valid: false,
		},
		// Less than or equal field tests
		{
			name: "lte field rule passes when field is less",
			data: map[string]interface{}{
				"field1": 3,
				"field2": 5,
			},
			rules: map[string]interface{}{
				"field1": "lte:field2",
			},
			valid: true,
		},
		{
			name: "lte field rule passes when field is equal",
			data: map[string]interface{}{
				"field1": 5,
				"field2": 5,
			},
			rules: map[string]interface{}{
				"field1": "lte:field2",
			},
			valid: true,
		},
		{
			name: "lte field rule fails when field is greater",
			data: map[string]interface{}{
				"field1": 10,
				"field2": 5,
			},
			rules: map[string]interface{}{
				"field1": "lte:field2",
			},
			valid: false,
		},
		// Same field tests
		{
			name: "same rule passes with identical values",
			data: map[string]interface{}{
				"password":         "secret123",
				"password_confirm": "secret123",
			},
			rules: map[string]interface{}{
				"password": "same:password_confirm",
			},
			valid: true,
		},
		{
			name: "same rule fails with different values",
			data: map[string]interface{}{
				"password":         "secret123",
				"password_confirm": "different",
			},
			rules: map[string]interface{}{
				"password": "same:password_confirm",
			},
			valid: false,
		},
		{
			name: "same rule with numbers",
			data: map[string]interface{}{
				"field1": 123,
				"field2": 123,
			},
			rules: map[string]interface{}{
				"field1": "same:field2",
			},
			valid: true,
		},
		// Different field tests
		{
			name: "different rule passes with different values",
			data: map[string]interface{}{
				"username": "john",
				"password": "secret123",
			},
			rules: map[string]interface{}{
				"username": "different:password",
			},
			valid: true,
		},
		{
			name: "different rule fails with same values",
			data: map[string]interface{}{
				"username": "secret123",
				"password": "secret123",
			},
			rules: map[string]interface{}{
				"username": "different:password",
			},
			valid: false,
		},
		{
			name: "different rule with multiple fields",
			data: map[string]interface{}{
				"field1": "value1",
				"field2": "value2",
				"field3": "value3",
			},
			rules: map[string]interface{}{
				"field1": "different:field2,field3",
			},
			valid: true,
		},
		{
			name: "different rule fails when same as one of multiple fields",
			data: map[string]interface{}{
				"field1": "value1",
				"field2": "value2",
				"field3": "value1",
			},
			rules: map[string]interface{}{
				"field1": "different:field2,field3",
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