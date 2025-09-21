package validation

import "testing"

func TestListRules(t *testing.T) {
	factory := NewFactory()
	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name:  "in rule passes",
			data:  map[string]interface{}{"status": "active"},
			rules: map[string]interface{}{"status": "in:active,inactive,pending"},
			valid: true,
		},
		{
			name:  "in rule fails",
			data:  map[string]interface{}{"status": "unknown"},
			rules: map[string]interface{}{"status": "in:active,inactive,pending"},
			valid: false,
		},
		{
			name:  "not_in rule passes",
			data:  map[string]interface{}{"status": "active"},
			rules: map[string]interface{}{"status": "not_in:deleted,banned"},
			valid: true,
		},
		{
			name:  "not_in rule fails",
			data:  map[string]interface{}{"status": "deleted"},
			rules: map[string]interface{}{"status": "not_in:deleted,banned"},
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
