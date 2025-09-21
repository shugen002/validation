package validation

import "testing"

func TestStringRules(t *testing.T) {
	factory := NewFactory()
	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name:  "email rule passes with valid email",
			data:  map[string]interface{}{"email": "user@example.com"},
			rules: map[string]interface{}{"email": "email"},
			valid: true,
		},
		{
			name:  "email rule fails with invalid email",
			data:  map[string]interface{}{"email": "not-an-email"},
			rules: map[string]interface{}{"email": "email"},
			valid: false,
		},
		{
			name:  "alpha rule passes with letters only",
			data:  map[string]interface{}{"name": "John"},
			rules: map[string]interface{}{"name": "alpha"},
			valid: true,
		},
		{
			name:  "alpha rule fails with numbers",
			data:  map[string]interface{}{"name": "John123"},
			rules: map[string]interface{}{"name": "alpha"},
			valid: false,
		},
		{
			name:  "alpha_num rule passes with letters and numbers",
			data:  map[string]interface{}{"username": "user123"},
			rules: map[string]interface{}{"username": "alpha_num"},
			valid: true,
		},
		{
			name:  "alpha_dash rule passes with letters, numbers, dashes and underscores",
			data:  map[string]interface{}{"username": "user_name-123"},
			rules: map[string]interface{}{"username": "alpha_dash"},
			valid: true,
		},
		{
			name:  "starts_with rule passes",
			data:  map[string]interface{}{"url": "https://example.com"},
			rules: map[string]interface{}{"url": "starts_with:https://,http://"},
			valid: true,
		},
		{
			name:  "starts_with rule fails",
			data:  map[string]interface{}{"url": "ftp://example.com"},
			rules: map[string]interface{}{"url": "starts_with:https://,http://"},
			valid: false,
		},
		{
			name:  "ends_with rule passes",
			data:  map[string]interface{}{"file": "document.pdf"},
			rules: map[string]interface{}{"file": "ends_with:.pdf,.doc"},
			valid: true,
		},
		{
			name:  "uppercase rule passes",
			data:  map[string]interface{}{"code": "ABCD"},
			rules: map[string]interface{}{"code": "uppercase"},
			valid: true,
		},
		{
			name:  "lowercase rule passes",
			data:  map[string]interface{}{"code": "abcd"},
			rules: map[string]interface{}{"code": "lowercase"},
			valid: true,
		},
		{
			name:  "regex rule passes",
			data:  map[string]interface{}{"phone": "123-456-7890"},
			rules: map[string]interface{}{"phone": `regex:^\d{3}-\d{3}-\d{4}$`},
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
