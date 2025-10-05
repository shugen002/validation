package validation

import (
	"encoding/json"
	"os"
	"testing"
)

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
		data  map[string]string
		rules map[string]string
		valid bool
	}{}
	for _, r := range files {
		rules := map[string]string{}
		data := map[string]string{}
		for _, v := range r.Variables {
			rules[v.Key] = v.Rule
			data[v.Key] = v.Default
		}
		test := struct {
			name  string
			data  map[string]string
			rules map[string]string
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
			validator, err := factory.Parse(test.rules)
			if err != nil {
				t.Fatalf("Failed to parse rules: %v", err)
			}
			err = validator.Validate(test.data)
			if (err == nil) != test.valid {
				t.Errorf("Validation result mismatch. Expected valid: %v, got error: %v", test.valid, err)
			}
		})
	}
}
