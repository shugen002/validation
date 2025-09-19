package validation

import (
	"testing"
)

// TestAdditionalCoverageBoost tests remaining uncovered functionality
func TestAdditionalCoverageBoost(t *testing.T) {
	factory := NewFactory()

	// Test all uncovered date functionality
	t.Run("CompleteDateRuleCoverage", func(t *testing.T) {
		data := map[string]interface{}{
			"start_date": "2023-01-01",
			"end_date":   "2023-12-31",
			"invalid":    "not-a-date",
		}

		// Test all date rule failure cases
		beforeRule := &BeforeRule{Date: "start_date"}
		beforeRule.SetData(data)
		
		// Test with invalid date
		if beforeRule.Passes("date", "invalid-date") {
			t.Error("BeforeRule should fail with invalid date")
		}
		
		// Test with field reference that has invalid date
		beforeRule2 := &BeforeRule{Date: "invalid"}
		beforeRule2.SetData(data)
		if beforeRule2.Passes("date", "2023-01-01") {
			t.Error("BeforeRule should fail when comparison field has invalid date")
		}

		// Test AfterRule edge cases
		afterRule := &AfterRule{Date: "end_date"}
		afterRule.SetData(data)
		if afterRule.Passes("date", "invalid-date") {
			t.Error("AfterRule should fail with invalid date")
		}

		// Test AfterOrEqualRule edge cases
		afterOrEqualRule := &AfterOrEqualRule{Date: "invalid"}
		afterOrEqualRule.SetData(data)
		if afterOrEqualRule.Passes("date", "2023-01-01") {
			t.Error("AfterOrEqualRule should fail when comparison field has invalid date")
		}

		// Test BeforeOrEqualRule edge cases
		beforeOrEqualRule := &BeforeOrEqualRule{Date: "invalid"}
		beforeOrEqualRule.SetData(data)
		if beforeOrEqualRule.Passes("date", "2023-01-01") {
			t.Error("BeforeOrEqualRule should fail when comparison field has invalid date")
		}

		// Test DateEqualsRule edge cases
		dateEqualsRule := &DateEqualsRule{Date: "invalid"}
		dateEqualsRule.SetData(data)
		if dateEqualsRule.Passes("date", "2023-01-01") {
			t.Error("DateEqualsRule should fail when comparison field has invalid date")
		}

		// Test with missing field reference
		dateEqualsRule2 := &DateEqualsRule{Date: "missing_field"}
		dateEqualsRule2.SetData(data)
		if dateEqualsRule2.Passes("date", "2023-01-01") {
			t.Error("DateEqualsRule should fail when comparison field doesn't exist")
		}
	})

	// Test all missing/present/prohibited rule combinations
	t.Run("CompleteSpecialRulesCoverage", func(t *testing.T) {
		data := map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
			"field3": "value3",
			"accepted": "yes",
			"declined": "no",
		}

		validator := factory.Make(data, map[string]interface{}{})

		// Test MissingUnlessRule - field missing, condition not met (should pass)
		missingUnlessRule := &MissingUnlessRule{Field: "missing_field", Value: "any_value"}
		missingUnlessRule.SetData(data)
		missingUnlessRule.SetValidator(validator)
		if !missingUnlessRule.Passes("another_missing_field", "") {
			t.Error("MissingUnlessRule should pass when trigger field doesn't exist")
		}

		// Test MissingWithRule - no trigger fields present (should pass)
		missingWithRule := &MissingWithRule{Fields: []string{"missing1", "missing2"}}
		missingWithRule.SetData(data)
		missingWithRule.SetValidator(validator)
		if !missingWithRule.Passes("test_field", "value") {
			t.Error("MissingWithRule should pass when no trigger fields are present")
		}

		// Test MissingWithAllRule - not all trigger fields present (should pass)
		missingWithAllRule := &MissingWithAllRule{Fields: []string{"field1", "missing_field"}}
		missingWithAllRule.SetData(data)
		missingWithAllRule.SetValidator(validator)
		if !missingWithAllRule.Passes("test_field", "value") {
			t.Error("MissingWithAllRule should pass when not all trigger fields are present")
		}

		// Test PresentIfRule - condition not met (should pass)
		presentIfRule := &PresentIfRule{Field: "field1", Value: "different_value"}
		presentIfRule.SetData(data)
		presentIfRule.SetValidator(validator)
		if !presentIfRule.Passes("test_field", "value") {
			t.Error("PresentIfRule should pass when condition is not met")
		}

		// Test PresentIfRule - condition met but field missing (should fail)
		presentIfRule2 := &PresentIfRule{Field: "field1", Value: "value1"}
		presentIfRule2.SetData(data)
		presentIfRule2.SetValidator(validator)
		if presentIfRule2.Passes("missing_field", "") {
			t.Error("PresentIfRule should fail when condition is met but field is missing")
		}

		// Test PresentUnlessRule - condition met (should pass)
		presentUnlessRule := &PresentUnlessRule{Field: "field1", Value: "value1"}
		presentUnlessRule.SetData(data)
		presentUnlessRule.SetValidator(validator)
		if !presentUnlessRule.Passes("test_field", "value") {
			t.Error("PresentUnlessRule should pass when condition is met")
		}

		// Test PresentUnlessRule - condition not met and field missing (should fail)
		presentUnlessRule2 := &PresentUnlessRule{Field: "missing_field", Value: "any_value"}
		presentUnlessRule2.SetData(data)
		presentUnlessRule2.SetValidator(validator)
		if presentUnlessRule2.Passes("missing_field2", "") {
			t.Error("PresentUnlessRule should fail when condition field doesn't exist and target field is missing")
		}

		// Test PresentWithRule - no trigger fields present (should pass)
		presentWithRule := &PresentWithRule{Fields: []string{"missing1", "missing2"}}
		presentWithRule.SetData(data)
		presentWithRule.SetValidator(validator)
		if !presentWithRule.Passes("test_field", "value") {
			t.Error("PresentWithRule should pass when no trigger fields are present")
		}

		// Test PresentWithRule - trigger field present but target field missing (should fail)
		presentWithRule2 := &PresentWithRule{Fields: []string{"field1"}}
		presentWithRule2.SetData(data)
		presentWithRule2.SetValidator(validator)
		if presentWithRule2.Passes("missing_field", "") {
			t.Error("PresentWithRule should fail when trigger field is present but target field is missing")
		}

		// Test PresentWithAllRule - not all trigger fields present (should pass)
		presentWithAllRule := &PresentWithAllRule{Fields: []string{"field1", "missing_field"}}
		presentWithAllRule.SetData(data)
		presentWithAllRule.SetValidator(validator)
		if !presentWithAllRule.Passes("test_field", "value") {
			t.Error("PresentWithAllRule should pass when not all trigger fields are present")
		}

		// Test PresentWithAllRule - all trigger fields present but target field missing (should fail)
		presentWithAllRule2 := &PresentWithAllRule{Fields: []string{"field1", "field2"}}
		presentWithAllRule2.SetData(data)
		presentWithAllRule2.SetValidator(validator)
		if presentWithAllRule2.Passes("missing_field", "") {
			t.Error("PresentWithAllRule should fail when all trigger fields are present but target field is missing")
		}

		// Test ProhibitedIfAcceptedRule - trigger not accepted (should pass)
		prohibitedIfAcceptedRule := &ProhibitedIfAcceptedRule{Field: "declined"}
		prohibitedIfAcceptedRule.SetData(data)
		if !prohibitedIfAcceptedRule.Passes("test_field", "value") {
			t.Error("ProhibitedIfAcceptedRule should pass when trigger field is not accepted")
		}

		// Test ProhibitedIfAcceptedRule - trigger missing (should pass)
		prohibitedIfAcceptedRule2 := &ProhibitedIfAcceptedRule{Field: "missing_field"}
		prohibitedIfAcceptedRule2.SetData(data)
		if !prohibitedIfAcceptedRule2.Passes("test_field", "value") {
			t.Error("ProhibitedIfAcceptedRule should pass when trigger field doesn't exist")
		}

		// Test ProhibitedIfDeclinedRule - trigger not declined (should pass)
		prohibitedIfDeclinedRule := &ProhibitedIfDeclinedRule{Field: "accepted"}
		prohibitedIfDeclinedRule.SetData(data)
		if !prohibitedIfDeclinedRule.Passes("test_field", "value") {
			t.Error("ProhibitedIfDeclinedRule should pass when trigger field is not declined")
		}

		// Test ProhibitedUnlessRule - condition met (should pass)
		prohibitedUnlessRule := &ProhibitedUnlessRule{Field: "field1", Value: "value1"}
		prohibitedUnlessRule.SetData(data)
		if !prohibitedUnlessRule.Passes("test_field", "") {
			t.Error("ProhibitedUnlessRule should pass when condition is met (empty value)")
		}

		// Test ProhibitedUnlessRule - trigger field missing and field has value (should fail)
		prohibitedUnlessRule2 := &ProhibitedUnlessRule{Field: "missing_field", Value: "any_value"}
		prohibitedUnlessRule2.SetData(data)
		if prohibitedUnlessRule2.Passes("test_field", "value") {
			t.Error("ProhibitedUnlessRule should fail when trigger field doesn't exist and target field has value")
		}

		// Test ProhibitsRule - no prohibited fields present (should pass)
		prohibitsRule := &ProhibitsRule{Fields: []string{"missing1", "missing2"}}
		prohibitsRule.SetData(data)
		if !prohibitsRule.Passes("test_field", "value") {
			t.Error("ProhibitsRule should pass when no prohibited fields are present")
		}

		// Test all comparison rules with missing fields
		gtRule := &GreaterThanRule{Field: "missing_field"}
		gtRule.SetData(data)
		if gtRule.Passes("test", "10") {
			t.Error("GreaterThanRule should fail when comparison field doesn't exist")
		}

		// Test with invalid numeric values
		gtRule2 := &GreaterThanRule{Field: "field1"}  // field1 = "value1" (non-numeric)
		gtRule2.SetData(data)
		if gtRule2.Passes("test", "not-a-number") {
			t.Error("GreaterThanRule should fail with non-numeric values")
		}

		gteRule := &GreaterThanOrEqualRule{Field: "missing_field"}
		gteRule.SetData(data)
		if gteRule.Passes("test", "10") {
			t.Error("GreaterThanOrEqualRule should fail when comparison field doesn't exist")
		}

		ltRule := &LessThanRule{Field: "missing_field"}
		ltRule.SetData(data)
		if ltRule.Passes("test", "10") {
			t.Error("LessThanRule should fail when comparison field doesn't exist")
		}

		lteRule := &LessThanOrEqualRule{Field: "missing_field"}
		lteRule.SetData(data)
		if lteRule.Passes("test", "10") {
			t.Error("LessThanOrEqualRule should fail when comparison field doesn't exist")
		}

		// Test RequiredArrayKeysRule with non-map value
		requiredArrayKeysRule := &RequiredArrayKeysRule{Keys: []string{"key1"}}
		if requiredArrayKeysRule.Passes("field", "not-a-map") {
			t.Error("RequiredArrayKeysRule should fail with non-map value")
		}
	})

	// Test edge cases for existing functionality
	t.Run("EdgeCaseCoverage", func(t *testing.T) {
		// Test DateFormatRule with empty format list
		dateFormatRule := &DateFormatRule{Formats: []string{}}
		if dateFormatRule.Passes("date", "2023-01-01") {
			t.Error("DateFormatRule should fail with empty format list")
		}

		// Test DateFormatRule with empty value
		dateFormatRule2 := &DateFormatRule{Formats: []string{"Y-m-d"}}
		if dateFormatRule2.Passes("date", "") {
			t.Error("DateFormatRule should fail with empty value")
		}

		// Test UuidRule with invalid version
		uuidRule := &UuidRule{Version: 99}
		// This should fail because version 99 doesn't exist
		if uuidRule.Passes("uuid", "f47ac10b-58cc-4372-a567-0e02b2c3d479") {
			t.Error("UuidRule with invalid version should fail")
		}

		// Test UlidRule with wrong length
		ulidRule := &UlidRule{}
		if ulidRule.Passes("ulid", "TOO-SHORT") {
			t.Error("UlidRule should fail with wrong length")
		}

		// Test UlidRule with invalid characters
		if ulidRule.Passes("ulid", "01ARZ3NDEKTSV4RRFFQ69G5FA!") {
			t.Error("UlidRule should fail with invalid characters")
		}

		// Test various edge cases for existing rules
		data := map[string]interface{}{
			"str_field": "test",
		}

		// Test RequiredWithoutRule with all fields present (should pass)
		requiredWithoutRule := &RequiredWithoutRule{Fields: []string{"str_field"}}
		requiredWithoutRule.SetData(data)
		if !requiredWithoutRule.Passes("optional_field", "") {
			t.Error("RequiredWithoutRule should pass when all trigger fields are present")
		}

		// Test RequiredIfAcceptedRule with non-accepted field (should pass)
		requiredIfAcceptedRule := &RequiredIfAcceptedRule{Field: "str_field"}
		requiredIfAcceptedRule.SetData(data)
		if !requiredIfAcceptedRule.Passes("optional_field", "") {
			t.Error("RequiredIfAcceptedRule should pass when trigger field is not accepted")
		}

		// Test RequiredIfDeclinedRule with non-declined field (should pass)
		requiredIfDeclinedRule := &RequiredIfDeclinedRule{Field: "str_field"}
		requiredIfDeclinedRule.SetData(data)
		if !requiredIfDeclinedRule.Passes("optional_field", "") {
			t.Error("RequiredIfDeclinedRule should pass when trigger field is not declined")
		}
	})

	// Test struct validation with complex scenarios
	t.Run("StructValidationEdgeCases", func(t *testing.T) {
		// Test struct with no validation tags
		type SimpleStruct struct {
			Name string
			Age  int
		}

		simple := SimpleStruct{Name: "test", Age: 25}
		validator, err := factory.ValidateStruct(simple)
		if err != nil {
			t.Errorf("ValidateStruct should handle struct with no validation tags: %v", err)
		}
		if !validator.Passes() {
			t.Error("Validation should pass for struct with no validation rules")
		}

		// Test struct with pointer fields
		type PointerStruct struct {
			Name *string `validate:"required"`
		}

		name := "test"
		pointer := PointerStruct{Name: &name}
		validator2, err := factory.ValidateStruct(pointer)
		if err != nil {
			t.Errorf("ValidateStruct should handle pointer fields: %v", err)
		}
		if !validator2.Passes() {
			t.Error("Validation should pass for valid pointer struct")
		}

		// Test with nil pointer
		nilPointer := PointerStruct{Name: nil}
		validator3, err := factory.ValidateStruct(nilPointer)
		if err != nil {
			t.Errorf("ValidateStruct should handle nil pointer fields: %v", err)
		}
		if validator3.Passes() {
			t.Error("Validation should fail for nil required pointer field")
		}
	})
}