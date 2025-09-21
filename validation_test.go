package validation

import (
	"encoding/json"
	"os"
	"testing"
)

// 关系型规则集成测试
func TestFieldRelationshipRules(t *testing.T) {
	factory := NewFactory()
	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "same rule passes",
			data: map[string]interface{}{
				"password":         "secret123",
				"password_confirm": "secret123",
			},
			rules: map[string]interface{}{
				"password_confirm": "same:password",
			},
			valid: true,
		},
		{
			name: "same rule fails",
			data: map[string]interface{}{
				"password":         "secret123",
				"password_confirm": "different",
			},
			rules: map[string]interface{}{
				"password_confirm": "same:password",
			},
			valid: false,
		},
		{
			name: "confirmed rule passes",
			data: map[string]interface{}{
				"password":              "secret123",
				"password_confirmation": "secret123",
			},
			rules: map[string]interface{}{
				"password": "confirmed",
			},
			valid: true,
		},
		{
			name: "confirmed rule fails",
			data: map[string]interface{}{
				"password":              "secret123",
				"password_confirmation": "different",
			},
			rules: map[string]interface{}{
				"password": "confirmed",
			},
			valid: false,
		},
		{
			name: "different rule passes",
			data: map[string]interface{}{
				"username": "john",
				"email":    "john@example.com",
			},
			rules: map[string]interface{}{
				"username": "different:email",
			},
			valid: true,
		},
		{
			name: "different rule fails",
			data: map[string]interface{}{
				"field1": "same_value",
				"field2": "same_value",
			},
			rules: map[string]interface{}{
				"field1": "different:field2",
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
					t.Errorf("Expected validation to fail, but it passed")
				}
			}
		})
	}
}

// 真实场景规则集成测试
func TestRealWorldRules(t *testing.T) {
	factory := NewFactory()
	file, err := os.ReadFile("testdata/eggs.json")
	if err != nil {
		t.Fatalf("Failed to read variables.json: %v", err)
	}
	var files []struct {
		File      string `json:"file"`
		Variables []struct {
			Name    string `json:"name"`
			Key     string `json:"key"`
			Rule    string `json:"rule"`
			Default string `json:"default"`
		}
	}
	if err := json.Unmarshal(file, &files); err != nil {
		t.Fatalf("Failed to unmarshal variables.json: %v", err)
	}
	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{}
	for _, r := range files {
		rules := map[string]interface{}{}
		data := map[string]interface{}{}
		for _, v := range r.Variables {
			rules[v.Key] = v.Rule
			data[v.Key] = v.Default
		}
		test := struct {
			name  string
			data  map[string]interface{}
			rules map[string]interface{}
			valid bool
		}{
			name:  r.File,
			data:  data,
			rules: rules,
			valid: true,
		}
		tests = append(tests, test)
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
