package validation

import "testing"

func TestDateRules(t *testing.T) {
	factory := NewFactory()
	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name:  "date rule passes with valid date",
			data:  map[string]interface{}{"birthday": "2022-01-01"},
			rules: map[string]interface{}{"birthday": "date"},
			valid: true,
		},
		{
			name:  "date rule fails with invalid date",
			data:  map[string]interface{}{"birthday": "not-a-date"},
			rules: map[string]interface{}{"birthday": "date"},
			valid: false,
		},
		{
			name:  "before rule passes",
			data:  map[string]interface{}{"event": "2020-01-01"},
			rules: map[string]interface{}{"event": "before:2021-01-01"},
			valid: true,
		},
		{
			name:  "before rule fails",
			data:  map[string]interface{}{"event": "2022-01-01"},
			rules: map[string]interface{}{"event": "before:2021-01-01"},
			valid: false,
		},
		{
			name:  "after rule passes",
			data:  map[string]interface{}{"event": "2022-01-01"},
			rules: map[string]interface{}{"event": "after:2021-01-01"},
			valid: true,
		},
		{
			name:  "after rule fails",
			data:  map[string]interface{}{"event": "2020-01-01"},
			rules: map[string]interface{}{"event": "after:2021-01-01"},
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
