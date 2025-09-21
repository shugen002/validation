package validation

import "testing"

func TestNumericRules(t *testing.T) {
	factory := NewFactory()
	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name:  "min rule passes",
			data:  map[string]interface{}{"age": 25},
			rules: map[string]interface{}{"age": "min:18"},
			valid: true,
		},
		{
			name:  "min rule fails",
			data:  map[string]interface{}{"age": 15},
			rules: map[string]interface{}{"age": "min:18"},
			valid: false,
		},
		{
			name:  "max rule passes",
			data:  map[string]interface{}{"age": 65},
			rules: map[string]interface{}{"age": "max:100"},
			valid: true,
		},
		{
			name:  "max rule fails",
			data:  map[string]interface{}{"age": 150},
			rules: map[string]interface{}{"age": "max:100"},
			valid: false,
		},
		{
			name:  "between rule passes",
			data:  map[string]interface{}{"age": 30},
			rules: map[string]interface{}{"age": "between:18,65"},
			valid: true,
		},
		{
			name:  "between rule fails - too low",
			data:  map[string]interface{}{"age": 15},
			rules: map[string]interface{}{"age": "between:18,65"},
			valid: false,
		},
		{
			name:  "between rule fails - too high",
			data:  map[string]interface{}{"age": 70},
			rules: map[string]interface{}{"age": "between:18,65"},
			valid: false,
		},
		{
			name:  "size rule passes for string length",
			data:  map[string]interface{}{"code": "ABCD"},
			rules: map[string]interface{}{"code": "size:4"},
			valid: true,
		},
		{
			name:  "size rule passes for number value",
			data:  map[string]interface{}{"quantity": 10},
			rules: map[string]interface{}{"quantity": "size:10"},
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
