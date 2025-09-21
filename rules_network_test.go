package validation

import "testing"

func TestNetworkRules(t *testing.T) {
	factory := NewFactory()
	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name:  "ip rule passes with IPv4",
			data:  map[string]interface{}{"ip": "192.168.1.1"},
			rules: map[string]interface{}{"ip": "ip"},
			valid: true,
		},
		{
			name:  "ip rule passes with IPv6",
			data:  map[string]interface{}{"ip": "2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
			rules: map[string]interface{}{"ip": "ip"},
			valid: true,
		},
		{
			name:  "ip rule fails with invalid IP",
			data:  map[string]interface{}{"ip": "not-an-ip"},
			rules: map[string]interface{}{"ip": "ip"},
			valid: false,
		},
		{
			name:  "ipv4 rule passes",
			data:  map[string]interface{}{"ip": "192.168.1.1"},
			rules: map[string]interface{}{"ip": "ipv4"},
			valid: true,
		},
		{
			name:  "ipv4 rule fails with IPv6",
			data:  map[string]interface{}{"ip": "2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
			rules: map[string]interface{}{"ip": "ipv4"},
			valid: false,
		},
		{
			name:  "ipv6 rule passes",
			data:  map[string]interface{}{"ip": "2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
			rules: map[string]interface{}{"ip": "ipv6"},
			valid: true,
		},
		{
			name:  "ipv6 rule fails with IPv4",
			data:  map[string]interface{}{"ip": "192.168.1.1"},
			rules: map[string]interface{}{"ip": "ipv6"},
			valid: false,
		},
		{
			name:  "mac_address rule passes",
			data:  map[string]interface{}{"mac": "00:1B:63:84:45:E6"},
			rules: map[string]interface{}{"mac": "mac_address"},
			valid: true,
		},
		{
			name:  "mac_address rule fails",
			data:  map[string]interface{}{"mac": "invalid-mac"},
			rules: map[string]interface{}{"mac": "mac_address"},
			valid: false,
		},
		// 可在此处继续补充更多网络相关规则用例
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
