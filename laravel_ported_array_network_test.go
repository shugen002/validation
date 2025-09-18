package validation

import (
	"testing"
)

// TestArrayAndListValidationRules tests ported from Laravel's ValidationArrayRuleTest.php and ValidationInRuleTest.php
func TestArrayAndListValidationRules(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		// In rule tests
		{
			name: "in rule passes with valid value",
			data: map[string]interface{}{
				"field": "apple",
			},
			rules: map[string]interface{}{
				"field": "in:apple,banana,orange",
			},
			valid: true,
		},
		{
			name: "in rule fails with invalid value",
			data: map[string]interface{}{
				"field": "grape",
			},
			rules: map[string]interface{}{
				"field": "in:apple,banana,orange",
			},
			valid: false,
		},
		{
			name: "in rule passes with numeric value",
			data: map[string]interface{}{
				"field": 1,
			},
			rules: map[string]interface{}{
				"field": "in:1,2,3",
			},
			valid: true,
		},
		{
			name: "in rule passes with string numeric value",
			data: map[string]interface{}{
				"field": "1",
			},
			rules: map[string]interface{}{
				"field": "in:1,2,3",
			},
			valid: true,
		},
		{
			name: "in rule with mixed types",
			data: map[string]interface{}{
				"field": "true",
			},
			rules: map[string]interface{}{
				"field": "in:true,false,1,0",
			},
			valid: true,
		},
		{
			name: "in rule with spaces in values",
			data: map[string]interface{}{
				"field": "value with spaces",
			},
			rules: map[string]interface{}{
				"field": "in:value with spaces,another value",
			},
			valid: true,
		},
		// Not in rule tests
		{
			name: "not_in rule passes with value not in list",
			data: map[string]interface{}{
				"field": "grape",
			},
			rules: map[string]interface{}{
				"field": "not_in:apple,banana,orange",
			},
			valid: true,
		},
		{
			name: "not_in rule fails with value in list",
			data: map[string]interface{}{
				"field": "apple",
			},
			rules: map[string]interface{}{
				"field": "not_in:apple,banana,orange",
			},
			valid: false,
		},
		{
			name: "not_in rule with numeric values",
			data: map[string]interface{}{
				"field": 4,
			},
			rules: map[string]interface{}{
				"field": "not_in:1,2,3",
			},
			valid: true,
		},
		{
			name: "not_in rule fails with numeric value in list",
			data: map[string]interface{}{
				"field": 2,
			},
			rules: map[string]interface{}{
				"field": "not_in:1,2,3",
			},
			valid: false,
		},
		// Array rule tests - basic array validation
		{
			name: "array rule passes with array",
			data: map[string]interface{}{
				"field": []string{"item1", "item2"},
			},
			rules: map[string]interface{}{
				"field": "array",
			},
			valid: true,
		},
		{
			name: "array rule passes with empty array",
			data: map[string]interface{}{
				"field": []string{},
			},
			rules: map[string]interface{}{
				"field": "array",
			},
			valid: true,
		},
		{
			name: "array rule passes with mixed type array",
			data: map[string]interface{}{
				"field": []interface{}{"string", 123, true},
			},
			rules: map[string]interface{}{
				"field": "array",
			},
			valid: true,
		},
		{
			name: "array rule passes with map",
			data: map[string]interface{}{
				"field": map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
			},
			rules: map[string]interface{}{
				"field": "array",
			},
			valid: true,
		},
		{
			name: "array rule fails with string",
			data: map[string]interface{}{
				"field": "hello",
			},
			rules: map[string]interface{}{
				"field": "array",
			},
			valid: false,
		},
		{
			name: "array rule fails with number",
			data: map[string]interface{}{
				"field": 123,
			},
			rules: map[string]interface{}{
				"field": "array",
			},
			valid: false,
		},
		{
			name: "array rule fails with boolean",
			data: map[string]interface{}{
				"field": true,
			},
			rules: map[string]interface{}{
				"field": "array",
			},
			valid: false,
		},
		// Array with allowed keys
		{
			name: "array rule with specific keys passes",
			data: map[string]interface{}{
				"field": map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
			},
			rules: map[string]interface{}{
				"field": "array:key1,key2,key3",
			},
			valid: true,
		},
		{
			name: "array rule with specific keys fails with invalid key",
			data: map[string]interface{}{
				"field": map[string]interface{}{
					"key1":    "value1",
					"invalid": "value2",
				},
			},
			rules: map[string]interface{}{
				"field": "array:key1,key2,key3",
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
					t.Error("Expected validation to fail, but it passed")
				}
			}
		})
	}
}

// TestNetworkValidationRules tests ported from Laravel validation tests
func TestNetworkValidationRules(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		// URL rule tests
		{
			name: "url rule passes with valid http URL",
			data: map[string]interface{}{
				"field": "http://example.com",
			},
			rules: map[string]interface{}{
				"field": "url",
			},
			valid: true,
		},
		{
			name: "url rule passes with valid https URL",
			data: map[string]interface{}{
				"field": "https://example.com",
			},
			rules: map[string]interface{}{
				"field": "url",
			},
			valid: true,
		},
		{
			name: "url rule passes with URL with path",
			data: map[string]interface{}{
				"field": "https://example.com/path/to/page",
			},
			rules: map[string]interface{}{
				"field": "url",
			},
			valid: true,
		},
		{
			name: "url rule passes with URL with query parameters",
			data: map[string]interface{}{
				"field": "https://example.com/search?q=test&page=1",
			},
			rules: map[string]interface{}{
				"field": "url",
			},
			valid: true,
		},
		{
			name: "url rule passes with URL with port",
			data: map[string]interface{}{
				"field": "http://example.com:8080",
			},
			rules: map[string]interface{}{
				"field": "url",
			},
			valid: true,
		},
		{
			name: "url rule fails with invalid URL",
			data: map[string]interface{}{
				"field": "not-a-url",
			},
			rules: map[string]interface{}{
				"field": "url",
			},
			valid: false,
		},
		{
			name: "url rule fails with URL without protocol",
			data: map[string]interface{}{
				"field": "example.com",
			},
			rules: map[string]interface{}{
				"field": "url",
			},
			valid: false,
		},
		// IP rule tests
		{
			name: "ip rule passes with valid IPv4",
			data: map[string]interface{}{
				"field": "192.168.1.1",
			},
			rules: map[string]interface{}{
				"field": "ip",
			},
			valid: true,
		},
		{
			name: "ip rule passes with valid IPv6",
			data: map[string]interface{}{
				"field": "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			},
			rules: map[string]interface{}{
				"field": "ip",
			},
			valid: true,
		},
		{
			name: "ip rule passes with compressed IPv6",
			data: map[string]interface{}{
				"field": "2001:db8:85a3::8a2e:370:7334",
			},
			rules: map[string]interface{}{
				"field": "ip",
			},
			valid: true,
		},
		{
			name: "ip rule passes with localhost IPv4",
			data: map[string]interface{}{
				"field": "127.0.0.1",
			},
			rules: map[string]interface{}{
				"field": "ip",
			},
			valid: true,
		},
		{
			name: "ip rule passes with localhost IPv6",
			data: map[string]interface{}{
				"field": "::1",
			},
			rules: map[string]interface{}{
				"field": "ip",
			},
			valid: true,
		},
		{
			name: "ip rule fails with invalid IP",
			data: map[string]interface{}{
				"field": "256.256.256.256",
			},
			rules: map[string]interface{}{
				"field": "ip",
			},
			valid: false,
		},
		{
			name: "ip rule fails with domain name",
			data: map[string]interface{}{
				"field": "example.com",
			},
			rules: map[string]interface{}{
				"field": "ip",
			},
			valid: false,
		},
		// IPv4 rule tests
		{
			name: "ipv4 rule passes with valid IPv4",
			data: map[string]interface{}{
				"field": "192.168.1.1",
			},
			rules: map[string]interface{}{
				"field": "ipv4",
			},
			valid: true,
		},
		{
			name: "ipv4 rule passes with localhost",
			data: map[string]interface{}{
				"field": "127.0.0.1",
			},
			rules: map[string]interface{}{
				"field": "ipv4",
			},
			valid: true,
		},
		{
			name: "ipv4 rule fails with IPv6",
			data: map[string]interface{}{
				"field": "2001:db8:85a3::8a2e:370:7334",
			},
			rules: map[string]interface{}{
				"field": "ipv4",
			},
			valid: false,
		},
		{
			name: "ipv4 rule fails with invalid IPv4",
			data: map[string]interface{}{
				"field": "256.256.256.256",
			},
			rules: map[string]interface{}{
				"field": "ipv4",
			},
			valid: false,
		},
		// IPv6 rule tests
		{
			name: "ipv6 rule passes with valid IPv6",
			data: map[string]interface{}{
				"field": "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			},
			rules: map[string]interface{}{
				"field": "ipv6",
			},
			valid: true,
		},
		{
			name: "ipv6 rule passes with compressed IPv6",
			data: map[string]interface{}{
				"field": "2001:db8:85a3::8a2e:370:7334",
			},
			rules: map[string]interface{}{
				"field": "ipv6",
			},
			valid: true,
		},
		{
			name: "ipv6 rule passes with localhost",
			data: map[string]interface{}{
				"field": "::1",
			},
			rules: map[string]interface{}{
				"field": "ipv6",
			},
			valid: true,
		},
		{
			name: "ipv6 rule fails with IPv4",
			data: map[string]interface{}{
				"field": "192.168.1.1",
			},
			rules: map[string]interface{}{
				"field": "ipv6",
			},
			valid: false,
		},
		{
			name: "ipv6 rule fails with invalid IPv6",
			data: map[string]interface{}{
				"field": "gggg::1",
			},
			rules: map[string]interface{}{
				"field": "ipv6",
			},
			valid: false,
		},
		// MAC address rule tests
		{
			name: "mac_address rule passes with valid MAC",
			data: map[string]interface{}{
				"field": "00:11:22:33:44:55",
			},
			rules: map[string]interface{}{
				"field": "mac_address",
			},
			valid: true,
		},
		{
			name: "mac_address rule passes with uppercase MAC",
			data: map[string]interface{}{
				"field": "AA:BB:CC:DD:EE:FF",
			},
			rules: map[string]interface{}{
				"field": "mac_address",
			},
			valid: true,
		},
		{
			name: "mac_address rule passes with mixed case MAC",
			data: map[string]interface{}{
				"field": "aA:bB:cC:dD:eE:fF",
			},
			rules: map[string]interface{}{
				"field": "mac_address",
			},
			valid: true,
		},
		{
			name: "mac_address rule passes with dash separator",
			data: map[string]interface{}{
				"field": "00-11-22-33-44-55",
			},
			rules: map[string]interface{}{
				"field": "mac_address",
			},
			valid: true,
		},
		{
			name: "mac_address rule fails with invalid MAC",
			data: map[string]interface{}{
				"field": "00:11:22:33:44",
			},
			rules: map[string]interface{}{
				"field": "mac_address",
			},
			valid: false,
		},
		{
			name: "mac_address rule fails with invalid characters",
			data: map[string]interface{}{
				"field": "00:11:22:33:44:GG",
			},
			rules: map[string]interface{}{
				"field": "mac_address",
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
					t.Error("Expected validation to fail, but it passed")
				}
			}
		})
	}
}