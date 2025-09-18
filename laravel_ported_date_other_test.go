package validation

import (
	"testing"
)

// TestDateValidationRules tests ported from Laravel's ValidationDateRuleTest.php
func TestDateValidationRules(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		// Date rule tests
		{
			name: "date rule passes with valid ISO date",
			data: map[string]interface{}{
				"field": "2024-01-15",
			},
			rules: map[string]interface{}{
				"field": "date",
			},
			valid: true,
		},
		{
			name: "date rule passes with valid datetime",
			data: map[string]interface{}{
				"field": "2024-01-15 14:30:00",
			},
			rules: map[string]interface{}{
				"field": "date",
			},
			valid: true,
		},
		{
			name: "date rule passes with valid ISO 8601",
			data: map[string]interface{}{
				"field": "2024-01-15T14:30:00Z",
			},
			rules: map[string]interface{}{
				"field": "date",
			},
			valid: true,
		},
		{
			name: "date rule fails with invalid date",
			data: map[string]interface{}{
				"field": "not-a-date",
			},
			rules: map[string]interface{}{
				"field": "date",
			},
			valid: false,
		},
		{
			name: "date rule fails with invalid date format",
			data: map[string]interface{}{
				"field": "01/15/2024",
			},
			rules: map[string]interface{}{
				"field": "date",
			},
			valid: false,
		},
		{
			name: "date rule fails with impossible date",
			data: map[string]interface{}{
				"field": "2024-02-30",
			},
			rules: map[string]interface{}{
				"field": "date",
			},
			valid: false,
		},
		// Date format rule tests
		{
			name: "date_format rule passes with matching format",
			data: map[string]interface{}{
				"field": "15/01/2024",
			},
			rules: map[string]interface{}{
				"field": "date_format:d/m/Y",
			},
			valid: true,
		},
		{
			name: "date_format rule passes with US format",
			data: map[string]interface{}{
				"field": "01/15/2024",
			},
			rules: map[string]interface{}{
				"field": "date_format:m/d/Y",
			},
			valid: true,
		},
		{
			name: "date_format rule passes with time format",
			data: map[string]interface{}{
				"field": "14:30:00",
			},
			rules: map[string]interface{}{
				"field": "date_format:H:i:s",
			},
			valid: true,
		},
		{
			name: "date_format rule fails with wrong format",
			data: map[string]interface{}{
				"field": "01/15/2024",
			},
			rules: map[string]interface{}{
				"field": "date_format:d/m/Y",
			},
			valid: false,
		},
		{
			name: "date_format rule fails with invalid date in format",
			data: map[string]interface{}{
				"field": "30/02/2024",
			},
			rules: map[string]interface{}{
				"field": "date_format:d/m/Y",
			},
			valid: false,
		},
		// After rule tests
		{
			name: "after rule passes with date after reference",
			data: map[string]interface{}{
				"field": "2024-02-01",
			},
			rules: map[string]interface{}{
				"field": "date|after:2024-01-01",
			},
			valid: true,
		},
		{
			name: "after rule fails with date before reference",
			data: map[string]interface{}{
				"field": "2023-12-01",
			},
			rules: map[string]interface{}{
				"field": "date|after:2024-01-01",
			},
			valid: false,
		},
		{
			name: "after rule fails with same date",
			data: map[string]interface{}{
				"field": "2024-01-01",
			},
			rules: map[string]interface{}{
				"field": "date|after:2024-01-01",
			},
			valid: false,
		},
		{
			name: "after rule with field reference passes",
			data: map[string]interface{}{
				"start_date": "2024-01-01",
				"end_date":   "2024-02-01",
			},
			rules: map[string]interface{}{
				"end_date": "date|after:start_date",
			},
			valid: true,
		},
		{
			name: "after rule with field reference fails",
			data: map[string]interface{}{
				"start_date": "2024-02-01",
				"end_date":   "2024-01-01",
			},
			rules: map[string]interface{}{
				"end_date": "date|after:start_date",
			},
			valid: false,
		},
		// Before rule tests
		{
			name: "before rule passes with date before reference",
			data: map[string]interface{}{
				"field": "2023-12-01",
			},
			rules: map[string]interface{}{
				"field": "date|before:2024-01-01",
			},
			valid: true,
		},
		{
			name: "before rule fails with date after reference",
			data: map[string]interface{}{
				"field": "2024-02-01",
			},
			rules: map[string]interface{}{
				"field": "date|before:2024-01-01",
			},
			valid: false,
		},
		{
			name: "before rule fails with same date",
			data: map[string]interface{}{
				"field": "2024-01-01",
			},
			rules: map[string]interface{}{
				"field": "date|before:2024-01-01",
			},
			valid: false,
		},
		// After or equal rule tests
		{
			name: "after_or_equal rule passes with date after reference",
			data: map[string]interface{}{
				"field": "2024-02-01",
			},
			rules: map[string]interface{}{
				"field": "date|after_or_equal:2024-01-01",
			},
			valid: true,
		},
		{
			name: "after_or_equal rule passes with same date",
			data: map[string]interface{}{
				"field": "2024-01-01",
			},
			rules: map[string]interface{}{
				"field": "date|after_or_equal:2024-01-01",
			},
			valid: true,
		},
		{
			name: "after_or_equal rule fails with date before reference",
			data: map[string]interface{}{
				"field": "2023-12-01",
			},
			rules: map[string]interface{}{
				"field": "date|after_or_equal:2024-01-01",
			},
			valid: false,
		},
		// Before or equal rule tests
		{
			name: "before_or_equal rule passes with date before reference",
			data: map[string]interface{}{
				"field": "2023-12-01",
			},
			rules: map[string]interface{}{
				"field": "date|before_or_equal:2024-01-01",
			},
			valid: true,
		},
		{
			name: "before_or_equal rule passes with same date",
			data: map[string]interface{}{
				"field": "2024-01-01",
			},
			rules: map[string]interface{}{
				"field": "date|before_or_equal:2024-01-01",
			},
			valid: true,
		},
		{
			name: "before_or_equal rule fails with date after reference",
			data: map[string]interface{}{
				"field": "2024-02-01",
			},
			rules: map[string]interface{}{
				"field": "date|before_or_equal:2024-01-01",
			},
			valid: false,
		},
		// Date equals rule tests
		{
			name: "date_equals rule passes with same date",
			data: map[string]interface{}{
				"field": "2024-01-01",
			},
			rules: map[string]interface{}{
				"field": "date|date_equals:2024-01-01",
			},
			valid: true,
		},
		{
			name: "date_equals rule fails with different date",
			data: map[string]interface{}{
				"field": "2024-01-02",
			},
			rules: map[string]interface{}{
				"field": "date|date_equals:2024-01-01",
			},
			valid: false,
		},
		{
			name: "date_equals rule with field reference passes",
			data: map[string]interface{}{
				"date1": "2024-01-01",
				"date2": "2024-01-01",
			},
			rules: map[string]interface{}{
				"date2": "date|date_equals:date1",
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
					t.Error("Expected validation to fail, but it passed")
				}
			}
		})
	}
}

// TestTimezoneValidationRule tests ported from Laravel validation tests
func TestTimezoneValidationRule(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name     string
		timezone string
		valid    bool
	}{
		{
			name:     "valid timezone UTC",
			timezone: "UTC",
			valid:    true,
		},
		{
			name:     "valid timezone America/New_York",
			timezone: "America/New_York",
			valid:    true,
		},
		{
			name:     "valid timezone Europe/London",
			timezone: "Europe/London",
			valid:    true,
		},
		{
			name:     "valid timezone Asia/Tokyo",
			timezone: "Asia/Tokyo",
			valid:    true,
		},
		{
			name:     "valid timezone Australia/Sydney",
			timezone: "Australia/Sydney",
			valid:    true,
		},
		{
			name:     "invalid timezone",
			timezone: "Invalid/Timezone",
			valid:    false,
		},
		{
			name:     "empty timezone",
			timezone: "",
			valid:    false,
		},
		{
			name:     "random string",
			timezone: "not-a-timezone",
			valid:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := map[string]interface{}{
				"timezone": test.timezone,
			}
			rules := map[string]interface{}{
				"timezone": "timezone",
			}

			validator := factory.Make(data, rules)

			if test.valid {
				if !validator.Passes() {
					t.Errorf("Expected timezone validation to pass, but it failed. Errors: %v", validator.Errors().All())
				}
			} else {
				if validator.Passes() {
					t.Error("Expected timezone validation to fail, but it passed")
				}
			}
		})
	}
}

// TestOtherValidationRules tests ported from Laravel validation tests for UUID, ULID, and hex color
func TestOtherValidationRules(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		// UUID rule tests
		{
			name: "uuid rule passes with valid UUID v4",
			data: map[string]interface{}{
				"field": "550e8400-e29b-41d4-a716-446655440000",
			},
			rules: map[string]interface{}{
				"field": "uuid",
			},
			valid: true,
		},
		{
			name: "uuid rule passes with valid UUID v1",
			data: map[string]interface{}{
				"field": "12345678-1234-1234-1234-123456789012",
			},
			rules: map[string]interface{}{
				"field": "uuid",
			},
			valid: true,
		},
		{
			name: "uuid rule passes with uppercase UUID",
			data: map[string]interface{}{
				"field": "550E8400-E29B-41D4-A716-446655440000",
			},
			rules: map[string]interface{}{
				"field": "uuid",
			},
			valid: true,
		},
		{
			name: "uuid rule fails with invalid UUID",
			data: map[string]interface{}{
				"field": "not-a-uuid",
			},
			rules: map[string]interface{}{
				"field": "uuid",
			},
			valid: false,
		},
		{
			name: "uuid rule fails with short UUID",
			data: map[string]interface{}{
				"field": "550e8400-e29b-41d4-a716",
			},
			rules: map[string]interface{}{
				"field": "uuid",
			},
			valid: false,
		},
		{
			name: "uuid rule fails with long UUID",
			data: map[string]interface{}{
				"field": "550e8400-e29b-41d4-a716-446655440000-extra",
			},
			rules: map[string]interface{}{
				"field": "uuid",
			},
			valid: false,
		},
		// ULID rule tests
		{
			name: "ulid rule passes with valid ULID",
			data: map[string]interface{}{
				"field": "01ARZ3NDEKTSV4RRFFQ69G5FAV",
			},
			rules: map[string]interface{}{
				"field": "ulid",
			},
			valid: true,
		},
		{
			name: "ulid rule passes with lowercase ULID",
			data: map[string]interface{}{
				"field": "01arz3ndektsv4rrffq69g5fav",
			},
			rules: map[string]interface{}{
				"field": "ulid",
			},
			valid: true,
		},
		{
			name: "ulid rule fails with invalid ULID",
			data: map[string]interface{}{
				"field": "not-a-ulid",
			},
			rules: map[string]interface{}{
				"field": "ulid",
			},
			valid: false,
		},
		{
			name: "ulid rule fails with short ULID",
			data: map[string]interface{}{
				"field": "01ARZ3NDEKTSV4RRFFQ69G5",
			},
			rules: map[string]interface{}{
				"field": "ulid",
			},
			valid: false,
		},
		{
			name: "ulid rule fails with long ULID",
			data: map[string]interface{}{
				"field": "01ARZ3NDEKTSV4RRFFQ69G5FAVEXTRA",
			},
			rules: map[string]interface{}{
				"field": "ulid",
			},
			valid: false,
		},
		{
			name: "ulid rule fails with invalid characters",
			data: map[string]interface{}{
				"field": "01ARZ3NDEKTSV4RRFFQ69G5FOU", // O and U are not valid
			},
			rules: map[string]interface{}{
				"field": "ulid",
			},
			valid: false,
		},
		// Hex color rule tests
		{
			name: "hex_color rule passes with 6-digit hex",
			data: map[string]interface{}{
				"field": "#FF5733",
			},
			rules: map[string]interface{}{
				"field": "hex_color",
			},
			valid: true,
		},
		{
			name: "hex_color rule passes with 3-digit hex",
			data: map[string]interface{}{
				"field": "#F57",
			},
			rules: map[string]interface{}{
				"field": "hex_color",
			},
			valid: true,
		},
		{
			name: "hex_color rule passes with lowercase hex",
			data: map[string]interface{}{
				"field": "#ff5733",
			},
			rules: map[string]interface{}{
				"field": "hex_color",
			},
			valid: true,
		},
		{
			name: "hex_color rule passes with mixed case hex",
			data: map[string]interface{}{
				"field": "#Ff5733",
			},
			rules: map[string]interface{}{
				"field": "hex_color",
			},
			valid: true,
		},
		{
			name: "hex_color rule fails without hash",
			data: map[string]interface{}{
				"field": "FF5733",
			},
			rules: map[string]interface{}{
				"field": "hex_color",
			},
			valid: false,
		},
		{
			name: "hex_color rule fails with invalid hex characters",
			data: map[string]interface{}{
				"field": "#GG5733",
			},
			rules: map[string]interface{}{
				"field": "hex_color",
			},
			valid: false,
		},
		{
			name: "hex_color rule fails with wrong length",
			data: map[string]interface{}{
				"field": "#FF57",
			},
			rules: map[string]interface{}{
				"field": "hex_color",
			},
			valid: false,
		},
		{
			name: "hex_color rule fails with too long hex",
			data: map[string]interface{}{
				"field": "#FF5733AA",
			},
			rules: map[string]interface{}{
				"field": "hex_color",
			},
			valid: false,
		},
		{
			name: "hex_color rule fails with empty string",
			data: map[string]interface{}{
				"field": "",
			},
			rules: map[string]interface{}{
				"field": "hex_color",
			},
			valid: false,
		},
		{
			name: "hex_color rule fails with just hash",
			data: map[string]interface{}{
				"field": "#",
			},
			rules: map[string]interface{}{
				"field": "hex_color",
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