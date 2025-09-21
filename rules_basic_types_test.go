package validation

import "testing"

func TestBasicRules(t *testing.T) {
	factory := NewFactory()
	tests := []struct {
		name     string
		data     map[string]interface{}
		rules    map[string]interface{}
		valid    bool
		messages map[string]string
	}{
		{
			name:  "required rule passes with value",
			data:  map[string]interface{}{"name": "John"},
			rules: map[string]interface{}{"name": "required"},
			valid: true,
		},
		{
			name:  "required rule fails with empty string",
			data:  map[string]interface{}{"name": ""},
			rules: map[string]interface{}{"name": "required"},
			valid: false,
		},
		{
			name:  "required rule fails with nil",
			data:  map[string]interface{}{"name": nil},
			rules: map[string]interface{}{"name": "required"},
			valid: false,
		},
		{
			name:  "string rule passes with string",
			data:  map[string]interface{}{"name": "John"},
			rules: map[string]interface{}{"name": "string"},
			valid: true,
		},
		{
			name:  "string rule fails with number",
			data:  map[string]interface{}{"age": 25},
			rules: map[string]interface{}{"age": "string"},
			valid: false,
		},
		{
			name:  "integer rule passes with int",
			data:  map[string]interface{}{"age": 25},
			rules: map[string]interface{}{"age": "integer"},
			valid: true,
		},
		{
			name:  "integer rule passes with string number",
			data:  map[string]interface{}{"age": "25"},
			rules: map[string]interface{}{"age": "integer"},
			valid: true,
		},
		{
			name:  "integer strict rule fails with string number",
			data:  map[string]interface{}{"age": "25"},
			rules: map[string]interface{}{"age": "integer:strict"},
			valid: false,
		},
		{
			name:  "numeric rule passes with float",
			data:  map[string]interface{}{"price": 19.99},
			rules: map[string]interface{}{"price": "numeric"},
			valid: true,
		},
		{
			name:  "boolean rule passes with bool",
			data:  map[string]interface{}{"active": true},
			rules: map[string]interface{}{"active": "boolean"},
			valid: true,
		},
		{
			name:  "boolean rule passes with string",
			data:  map[string]interface{}{"active": "true"},
			rules: map[string]interface{}{"active": "boolean"},
			valid: true,
		},
		{
			name:  "boolean strict rule fails with string",
			data:  map[string]interface{}{"active": "true"},
			rules: map[string]interface{}{"active": "boolean:strict"},
			valid: false,
		},
		{
			name:  "array rule passes with slice",
			data:  map[string]interface{}{"tags": []interface{}{"go", "validation"}},
			rules: map[string]interface{}{"tags": "array"},
			valid: true,
		},
		{
			name:  "json rule passes with valid JSON",
			data:  map[string]interface{}{"config": `{"key": "value"}`},
			rules: map[string]interface{}{"config": "json"},
			valid: true,
		},
		{
			name:  "json rule fails with invalid JSON",
			data:  map[string]interface{}{"config": `{"key": value}`},
			rules: map[string]interface{}{"config": "json"},
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
