package validation

import (
	"testing"
)

func TestMinMaxRuleWithNumericStrings(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name     string
		data     map[string]interface{}
		rules    map[string]interface{}
		expected bool
	}{
		{
			name: "SSH Port validation - integer min rule with numeric string",
			data: map[string]interface{}{
				"SSH_PORT": "2020",
			},
			rules: map[string]interface{}{
				"SSH_PORT": "required|integer|min:1025|max:65535",
			},
			expected: true,
		},
		{
			name: "Creature health modifier - numeric min rule with numeric string",
			data: map[string]interface{}{
				"HEALTH": "100",
			},
			rules: map[string]interface{}{
				"HEALTH": "required|numeric|min:20|max:300",
			},
			expected: true,
		},
		{
			name: "String length validation still works - string max rule",
			data: map[string]interface{}{
				"PORT": "8377",
			},
			rules: map[string]interface{}{
				"PORT": "required|string|max:10",
			},
			expected: true, // "8377" has 4 characters, which is <= 10
		},
		{
			name: "Large numeric string with string rule - should validate length",
			data: map[string]interface{}{
				"LARGE_PORT": "27015",
			},
			rules: map[string]interface{}{
				"LARGE_PORT": "required|string|max:100",
			},
			expected: true, // "27015" has 5 characters, which is <= 100
		},
		{
			name: "Numeric string out of range with integer rule",
			data: map[string]interface{}{
				"SMALL_PORT": "500",
			},
			rules: map[string]interface{}{
				"SMALL_PORT": "required|integer|min:1025|max:65535",
			},
			expected: false, // 500 < 1025
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validator := factory.Make(test.data, test.rules)

			if test.expected {
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