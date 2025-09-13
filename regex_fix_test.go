package validation

import (
	"testing"
)

// TestRegexRuleWithCommas tests that regex patterns containing commas are parsed correctly
func TestRegexRuleWithCommas(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name    string
		data    map[string]interface{}
		rules   map[string]interface{}
		valid   bool
		desc    string
	}{
		{
			name: "regex with quantifier containing comma should pass",
			data: map[string]interface{}{
				"map_name": "ctf_2fort",
			},
			rules: map[string]interface{}{
				"map_name": "required|regex:/^(\\w{1,20})$/",
			},
			valid: true,
			desc:  "Map name matching regex with {1,20} quantifier",
		},
		{
			name: "regex with quantifier containing comma should fail invalid input",
			data: map[string]interface{}{
				"map_name": "this_is_way_too_long_for_the_regex_pattern_to_match",
			},
			rules: map[string]interface{}{
				"map_name": "required|regex:/^(\\w{1,20})$/",
			},
			valid: false,
			desc:  "Map name not matching regex with {1,20} quantifier",
		},
		{
			name: "regex without slash delimiters should work",
			data: map[string]interface{}{
				"code": "ABC123",
			},
			rules: map[string]interface{}{
				"code": "required|regex:^[A-Z]{1,3}[0-9]{1,5}$",
			},
			valid: true,
			desc:  "Code matching regex without slash delimiters",
		},
		{
			name: "not_regex with quantifier containing comma should pass",
			data: map[string]interface{}{
				"password": "verylongpassword",
			},
			rules: map[string]interface{}{
				"password": "required|not_regex:/^.{1,5}$/",
			},
			valid: true,
			desc:  "Password should not match short pattern (should pass because it's long)",
		},
		{
			name: "complex regex with multiple commas should work",
			data: map[string]interface{}{
				"complex": "test123",
			},
			rules: map[string]interface{}{
				"complex": "required|regex:/^[a-z]{1,10}[0-9]{1,5}$/",
			},
			valid: true,
			desc:  "Complex pattern with multiple quantifiers",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validator := factory.Make(test.data, test.rules)

			if test.valid {
				if !validator.Passes() {
					t.Errorf("%s: Expected validation to pass, but it failed. Errors: %v", test.desc, validator.Errors().All())
				}
			} else {
				if validator.Passes() {
					t.Errorf("%s: Expected validation to fail, but it passed", test.desc)
				}
			}
		})
	}
}