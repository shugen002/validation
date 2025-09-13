package validation

import (
	"encoding/json"
	"os"
	"testing"
)

func TestBasicRules(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name     string
		data     map[string]interface{}
		rules    map[string]interface{}
		valid    bool
		messages map[string]string
	}{
		{
			name: "required rule passes with value",
			data: map[string]interface{}{
				"name": "John",
			},
			rules: map[string]interface{}{
				"name": "required",
			},
			valid: true,
		},
		{
			name: "required rule fails with empty string",
			data: map[string]interface{}{
				"name": "",
			},
			rules: map[string]interface{}{
				"name": "required",
			},
			valid: false,
		},
		{
			name: "required rule fails with nil",
			data: map[string]interface{}{
				"name": nil,
			},
			rules: map[string]interface{}{
				"name": "required",
			},
			valid: false,
		},
		{
			name: "string rule passes with string",
			data: map[string]interface{}{
				"name": "John",
			},
			rules: map[string]interface{}{
				"name": "string",
			},
			valid: true,
		},
		{
			name: "string rule fails with number",
			data: map[string]interface{}{
				"age": 25,
			},
			rules: map[string]interface{}{
				"age": "string",
			},
			valid: false,
		},
		{
			name: "integer rule passes with int",
			data: map[string]interface{}{
				"age": 25,
			},
			rules: map[string]interface{}{
				"age": "integer",
			},
			valid: true,
		},
		{
			name: "integer rule passes with string number",
			data: map[string]interface{}{
				"age": "25",
			},
			rules: map[string]interface{}{
				"age": "integer",
			},
			valid: true,
		},
		{
			name: "integer strict rule fails with string number",
			data: map[string]interface{}{
				"age": "25",
			},
			rules: map[string]interface{}{
				"age": "integer:strict",
			},
			valid: false,
		},
		{
			name: "numeric rule passes with float",
			data: map[string]interface{}{
				"price": 19.99,
			},
			rules: map[string]interface{}{
				"price": "numeric",
			},
			valid: true,
		},
		{
			name: "boolean rule passes with bool",
			data: map[string]interface{}{
				"active": true,
			},
			rules: map[string]interface{}{
				"active": "boolean",
			},
			valid: true,
		},
		{
			name: "boolean rule passes with string",
			data: map[string]interface{}{
				"active": "true",
			},
			rules: map[string]interface{}{
				"active": "boolean",
			},
			valid: true,
		},
		{
			name: "boolean strict rule fails with string",
			data: map[string]interface{}{
				"active": "true",
			},
			rules: map[string]interface{}{
				"active": "boolean:strict",
			},
			valid: false,
		},
		{
			name: "array rule passes with slice",
			data: map[string]interface{}{
				"tags": []interface{}{"go", "validation"},
			},
			rules: map[string]interface{}{
				"tags": "array",
			},
			valid: true,
		},
		{
			name: "json rule passes with valid JSON",
			data: map[string]interface{}{
				"config": `{"key": "value"}`,
			},
			rules: map[string]interface{}{
				"config": "json",
			},
			valid: true,
		},
		{
			name: "json rule fails with invalid JSON",
			data: map[string]interface{}{
				"config": `{"key": value}`,
			},
			rules: map[string]interface{}{
				"config": "json",
			},
			valid: false,
		},
		// Accepted rule tests
		{
			name: "accepted rule passes with true boolean",
			data: map[string]interface{}{
				"terms": true,
			},
			rules: map[string]interface{}{
				"terms": "accepted",
			},
			valid: true,
		},
		{
			name: "accepted rule passes with 'yes' string",
			data: map[string]interface{}{
				"terms": "yes",
			},
			rules: map[string]interface{}{
				"terms": "accepted",
			},
			valid: true,
		},
		{
			name: "accepted rule passes with 'on' string",
			data: map[string]interface{}{
				"terms": "on",
			},
			rules: map[string]interface{}{
				"terms": "accepted",
			},
			valid: true,
		},
		{
			name: "accepted rule passes with '1' string",
			data: map[string]interface{}{
				"terms": "1",
			},
			rules: map[string]interface{}{
				"terms": "accepted",
			},
			valid: true,
		},
		{
			name: "accepted rule passes with 'true' string",
			data: map[string]interface{}{
				"terms": "true",
			},
			rules: map[string]interface{}{
				"terms": "accepted",
			},
			valid: true,
		},
		{
			name: "accepted rule passes with integer 1",
			data: map[string]interface{}{
				"terms": 1,
			},
			rules: map[string]interface{}{
				"terms": "accepted",
			},
			valid: true,
		},
		{
			name: "accepted rule passes with float 1.0",
			data: map[string]interface{}{
				"terms": 1.0,
			},
			rules: map[string]interface{}{
				"terms": "accepted",
			},
			valid: true,
		},
		{
			name: "accepted rule fails with false boolean",
			data: map[string]interface{}{
				"terms": false,
			},
			rules: map[string]interface{}{
				"terms": "accepted",
			},
			valid: false,
		},
		{
			name: "accepted rule fails with 'no' string",
			data: map[string]interface{}{
				"terms": "no",
			},
			rules: map[string]interface{}{
				"terms": "accepted",
			},
			valid: false,
		},
		{
			name: "accepted rule fails with integer 0",
			data: map[string]interface{}{
				"terms": 0,
			},
			rules: map[string]interface{}{
				"terms": "accepted",
			},
			valid: false,
		},
		// Declined rule tests
		{
			name: "declined rule passes with false boolean",
			data: map[string]interface{}{
				"newsletter": false,
			},
			rules: map[string]interface{}{
				"newsletter": "declined",
			},
			valid: true,
		},
		{
			name: "declined rule passes with 'no' string",
			data: map[string]interface{}{
				"newsletter": "no",
			},
			rules: map[string]interface{}{
				"newsletter": "declined",
			},
			valid: true,
		},
		{
			name: "declined rule passes with 'off' string",
			data: map[string]interface{}{
				"newsletter": "off",
			},
			rules: map[string]interface{}{
				"newsletter": "declined",
			},
			valid: true,
		},
		{
			name: "declined rule passes with '0' string",
			data: map[string]interface{}{
				"newsletter": "0",
			},
			rules: map[string]interface{}{
				"newsletter": "declined",
			},
			valid: true,
		},
		{
			name: "declined rule passes with 'false' string",
			data: map[string]interface{}{
				"newsletter": "false",
			},
			rules: map[string]interface{}{
				"newsletter": "declined",
			},
			valid: true,
		},
		{
			name: "declined rule passes with integer 0",
			data: map[string]interface{}{
				"newsletter": 0,
			},
			rules: map[string]interface{}{
				"newsletter": "declined",
			},
			valid: true,
		},
		{
			name: "declined rule passes with float 0.0",
			data: map[string]interface{}{
				"newsletter": 0.0,
			},
			rules: map[string]interface{}{
				"newsletter": "declined",
			},
			valid: true,
		},
		{
			name: "declined rule fails with true boolean",
			data: map[string]interface{}{
				"newsletter": true,
			},
			rules: map[string]interface{}{
				"newsletter": "declined",
			},
			valid: false,
		},
		{
			name: "declined rule fails with 'yes' string",
			data: map[string]interface{}{
				"newsletter": "yes",
			},
			rules: map[string]interface{}{
				"newsletter": "declined",
			},
			valid: false,
		},
		{
			name: "declined rule fails with integer 1",
			data: map[string]interface{}{
				"newsletter": 1,
			},
			rules: map[string]interface{}{
				"newsletter": "declined",
			},
			valid: false,
		},
		// Accepted If rule tests
		{
			name: "accepted_if rule passes when condition met and field accepted",
			data: map[string]interface{}{
				"payment_method": "card",
				"save_card":      "yes",
			},
			rules: map[string]interface{}{
				"save_card": "accepted_if:payment_method,card",
			},
			valid: true,
		},
		{
			name: "accepted_if rule fails when condition met and field not accepted",
			data: map[string]interface{}{
				"payment_method": "card",
				"save_card":      "no",
			},
			rules: map[string]interface{}{
				"save_card": "accepted_if:payment_method,card",
			},
			valid: false,
		},
		{
			name: "accepted_if rule passes when condition not met",
			data: map[string]interface{}{
				"payment_method": "cash",
				"save_card":      "no",
			},
			rules: map[string]interface{}{
				"save_card": "accepted_if:payment_method,card",
			},
			valid: true,
		},
		{
			name: "accepted_if rule passes when condition field missing",
			data: map[string]interface{}{
				"save_card": "no",
			},
			rules: map[string]interface{}{
				"save_card": "accepted_if:payment_method,card",
			},
			valid: true,
		},
		// Declined If rule tests
		{
			name: "declined_if rule passes when condition met and field declined",
			data: map[string]interface{}{
				"account_type": "guest",
				"newsletter":   "no",
			},
			rules: map[string]interface{}{
				"newsletter": "declined_if:account_type,guest",
			},
			valid: true,
		},
		{
			name: "declined_if rule fails when condition met and field not declined",
			data: map[string]interface{}{
				"account_type": "guest",
				"newsletter":   "yes",
			},
			rules: map[string]interface{}{
				"newsletter": "declined_if:account_type,guest",
			},
			valid: false,
		},
		{
			name: "declined_if rule passes when condition not met",
			data: map[string]interface{}{
				"account_type": "premium",
				"newsletter":   "yes",
			},
			rules: map[string]interface{}{
				"newsletter": "declined_if:account_type,guest",
			},
			valid: true,
		},
		{
			name: "declined_if rule passes when condition field missing",
			data: map[string]interface{}{
				"newsletter": "yes",
			},
			rules: map[string]interface{}{
				"newsletter": "declined_if:account_type,guest",
			},
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

func TestStringRules(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "email rule passes with valid email",
			data: map[string]interface{}{
				"email": "user@example.com",
			},
			rules: map[string]interface{}{
				"email": "email",
			},
			valid: true,
		},
		{
			name: "email rule fails with invalid email",
			data: map[string]interface{}{
				"email": "not-an-email",
			},
			rules: map[string]interface{}{
				"email": "email",
			},
			valid: false,
		},
		{
			name: "alpha rule passes with letters only",
			data: map[string]interface{}{
				"name": "John",
			},
			rules: map[string]interface{}{
				"name": "alpha",
			},
			valid: true,
		},
		{
			name: "alpha rule fails with numbers",
			data: map[string]interface{}{
				"name": "John123",
			},
			rules: map[string]interface{}{
				"name": "alpha",
			},
			valid: false,
		},
		{
			name: "alpha_num rule passes with letters and numbers",
			data: map[string]interface{}{
				"username": "user123",
			},
			rules: map[string]interface{}{
				"username": "alpha_num",
			},
			valid: true,
		},
		{
			name: "alpha_dash rule passes with letters, numbers, dashes and underscores",
			data: map[string]interface{}{
				"username": "user_name-123",
			},
			rules: map[string]interface{}{
				"username": "alpha_dash",
			},
			valid: true,
		},
		{
			name: "starts_with rule passes",
			data: map[string]interface{}{
				"url": "https://example.com",
			},
			rules: map[string]interface{}{
				"url": "starts_with:https://,http://",
			},
			valid: true,
		},
		{
			name: "starts_with rule fails",
			data: map[string]interface{}{
				"url": "ftp://example.com",
			},
			rules: map[string]interface{}{
				"url": "starts_with:https://,http://",
			},
			valid: false,
		},
		{
			name: "ends_with rule passes",
			data: map[string]interface{}{
				"file": "document.pdf",
			},
			rules: map[string]interface{}{
				"file": "ends_with:.pdf,.doc",
			},
			valid: true,
		},
		{
			name: "uppercase rule passes",
			data: map[string]interface{}{
				"code": "ABCD",
			},
			rules: map[string]interface{}{
				"code": "uppercase",
			},
			valid: true,
		},
		{
			name: "lowercase rule passes",
			data: map[string]interface{}{
				"code": "abcd",
			},
			rules: map[string]interface{}{
				"code": "lowercase",
			},
			valid: true,
		},
		{
			name: "regex rule passes",
			data: map[string]interface{}{
				"phone": "123-456-7890",
			},
			rules: map[string]interface{}{
				"phone": `regex:^\d{3}-\d{3}-\d{4}$`,
			},
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

func TestNumericRules(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "min rule passes",
			data: map[string]interface{}{
				"age": 25,
			},
			rules: map[string]interface{}{
				"age": "min:18",
			},
			valid: true,
		},
		{
			name: "min rule fails",
			data: map[string]interface{}{
				"age": 15,
			},
			rules: map[string]interface{}{
				"age": "min:18",
			},
			valid: false,
		},
		{
			name: "max rule passes",
			data: map[string]interface{}{
				"age": 65,
			},
			rules: map[string]interface{}{
				"age": "max:100",
			},
			valid: true,
		},
		{
			name: "max rule fails",
			data: map[string]interface{}{
				"age": 150,
			},
			rules: map[string]interface{}{
				"age": "max:100",
			},
			valid: false,
		},
		{
			name: "between rule passes",
			data: map[string]interface{}{
				"age": 30,
			},
			rules: map[string]interface{}{
				"age": "between:18,65",
			},
			valid: true,
		},
		{
			name: "between rule fails - too low",
			data: map[string]interface{}{
				"age": 15,
			},
			rules: map[string]interface{}{
				"age": "between:18,65",
			},
			valid: false,
		},
		{
			name: "between rule fails - too high",
			data: map[string]interface{}{
				"age": 70,
			},
			rules: map[string]interface{}{
				"age": "between:18,65",
			},
			valid: false,
		},
		{
			name: "size rule passes for string length",
			data: map[string]interface{}{
				"code": "ABCD",
			},
			rules: map[string]interface{}{
				"code": "size:4",
			},
			valid: true,
		},
		{
			name: "size rule passes for number value",
			data: map[string]interface{}{
				"quantity": 10,
			},
			rules: map[string]interface{}{
				"quantity": "size:10",
			},
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

func TestListRules(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "in rule passes",
			data: map[string]interface{}{
				"status": "active",
			},
			rules: map[string]interface{}{
				"status": "in:active,inactive,pending",
			},
			valid: true,
		},
		{
			name: "in rule fails",
			data: map[string]interface{}{
				"status": "unknown",
			},
			rules: map[string]interface{}{
				"status": "in:active,inactive,pending",
			},
			valid: false,
		},
		{
			name: "not_in rule passes",
			data: map[string]interface{}{
				"status": "active",
			},
			rules: map[string]interface{}{
				"status": "not_in:deleted,banned",
			},
			valid: true,
		},
		{
			name: "not_in rule fails",
			data: map[string]interface{}{
				"status": "deleted",
			},
			rules: map[string]interface{}{
				"status": "not_in:deleted,banned",
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

func TestNetworkRules(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "url rule passes with valid URL",
			data: map[string]interface{}{
				"website": "https://example.com",
			},
			rules: map[string]interface{}{
				"website": "url",
			},
			valid: true,
		},
		{
			name: "url rule fails with invalid URL",
			data: map[string]interface{}{
				"website": "not-a-url",
			},
			rules: map[string]interface{}{
				"website": "url",
			},
			valid: false,
		},
		{
			name: "ip rule passes with valid IPv4",
			data: map[string]interface{}{
				"server": "192.168.1.1",
			},
			rules: map[string]interface{}{
				"server": "ip",
			},
			valid: true,
		},
		{
			name: "ip rule passes with valid IPv6",
			data: map[string]interface{}{
				"server": "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			},
			rules: map[string]interface{}{
				"server": "ip",
			},
			valid: true,
		},
		{
			name: "ipv4 rule passes with IPv4",
			data: map[string]interface{}{
				"server": "192.168.1.1",
			},
			rules: map[string]interface{}{
				"server": "ipv4",
			},
			valid: true,
		},
		{
			name: "ipv4 rule fails with IPv6",
			data: map[string]interface{}{
				"server": "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			},
			rules: map[string]interface{}{
				"server": "ipv4",
			},
			valid: false,
		},
		{
			name: "ipv6 rule passes with IPv6",
			data: map[string]interface{}{
				"server": "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			},
			rules: map[string]interface{}{
				"server": "ipv6",
			},
			valid: true,
		},
		{
			name: "ipv6 rule fails with IPv4",
			data: map[string]interface{}{
				"server": "192.168.1.1",
			},
			rules: map[string]interface{}{
				"server": "ipv6",
			},
			valid: false,
		},
		{
			name: "mac_address rule passes",
			data: map[string]interface{}{
				"mac": "00:1B:63:84:45:E6",
			},
			rules: map[string]interface{}{
				"mac": "mac_address",
			},
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

func TestRealWorldRules(t *testing.T) {
	factory := NewFactory()
	// read from tests/variables.json
	file, err := os.ReadFile("tests/variables.json")
	if err != nil {
		t.Fatalf("Failed to read variables.json: %v", err)
	}

	var rules []struct {
		Name    string `json:"name"`
		Key     string `json:"key"`
		Rule    string `json:"rule"`
		Default string `json:"default"`
	}
	if err := json.Unmarshal(file, &rules); err != nil {
		t.Fatalf("Failed to unmarshal variables.json: %v", err)
	}

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{}

	for _, r := range rules {
		test := struct {
			name  string
			data  map[string]interface{}
			rules map[string]interface{}
			valid bool
		}{
			name: r.Name,
			data: map[string]interface{}{
				r.Key: r.Default,
			},
			rules: map[string]interface{}{
				r.Key: r.Rule,
			},
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
