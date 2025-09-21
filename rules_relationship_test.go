package validation

import "testing"

func TestRelationshipRules(t *testing.T) {
	factory := NewFactory()
	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name:  "required_with rule passes",
			data:  map[string]interface{}{"a": "foo", "b": "bar"},
			rules: map[string]interface{}{"a": "required_with:b"},
			valid: true,
		},
		{
			name:  "required_with rule fails",
			data:  map[string]interface{}{"b": "bar"},
			rules: map[string]interface{}{"a": "required_with:b"},
			valid: false,
		},
		{
			name:  "required_without rule passes",
			data:  map[string]interface{}{"a": "foo"},
			rules: map[string]interface{}{"a": "required_without:b"},
			valid: true,
		},
		{
			name:  "required_without rule passes 2",
			data:  map[string]interface{}{"b": "bar"},
			rules: map[string]interface{}{"a": "required_without:b"},
			valid: true,
		},
		{
			name:  "required_without rule fails",
			data:  map[string]interface{}{},
			rules: map[string]interface{}{"a": "required_without:b"},
			valid: false,
		},
		{
			name:  "same rule passes",
			data:  map[string]interface{}{"password": "123", "password_confirm": "123"},
			rules: map[string]interface{}{"password_confirm": "same:password"},
			valid: true,
		},
		{
			name:  "same rule fails",
			data:  map[string]interface{}{"password": "123", "password_confirm": "456"},
			rules: map[string]interface{}{"password_confirm": "same:password"},
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
