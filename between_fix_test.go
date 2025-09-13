package validation

import (
	"testing"
)

// TestBetweenRuleWithNumericStrings tests the specific fix for issue #4
// where BetweenRule should handle numeric strings correctly
func TestBetweenRuleWithNumericStrings(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
		desc  string
	}{
		{
			name: "SSH Port validation - numeric string in valid range",
			data: map[string]interface{}{
				"SSH_PORT": "2020",
			},
			rules: map[string]interface{}{
				"SSH_PORT": "required|integer|between:1024,65535",
			},
			valid: true,
			desc:  "Should pass because 2020 is between 1024 and 65535",
		},
		{
			name: "SSH Port validation - numeric string out of range",
			data: map[string]interface{}{
				"SSH_PORT": "100",
			},
			rules: map[string]interface{}{
				"SSH_PORT": "required|integer|between:1024,65535",
			},
			valid: false,
			desc:  "Should fail because 100 is not between 1024 and 65535",
		},
		{
			name: "String length validation still works",
			data: map[string]interface{}{
				"code": "ABCD",
			},
			rules: map[string]interface{}{
				"code": "string|between:3,5",
			},
			valid: true,
			desc:  "Should pass because string length 4 is between 3 and 5",
		},
		{
			name: "Non-numeric string uses length",
			data: map[string]interface{}{
				"name": "toolong",
			},
			rules: map[string]interface{}{
				"name": "string|between:3,5",
			},
			valid: false,
			desc:  "Should fail because string length 7 is not between 3 and 5",
		},
		{
			name: "Numeric value works unchanged",
			data: map[string]interface{}{
				"count": 50,
			},
			rules: map[string]interface{}{
				"count": "integer|between:1,100",
			},
			valid: true,
			desc:  "Should pass because 50 is between 1 and 100",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validator := factory.Make(test.data, test.rules)

			if test.valid {
				if !validator.Passes() {
					t.Errorf("Expected validation to pass, but it failed. %s. Errors: %v", test.desc, validator.Errors().All())
				}
			} else {
				if validator.Passes() {
					t.Errorf("Expected validation to fail, but it passed. %s", test.desc)
				}
			}
		})
	}
}