package validation

import (
	"testing"
	"time"
)

// TestCoverageImprovement tests all the uncovered functionality to boost coverage
func TestCoverageImprovement(t *testing.T) {
	factory := NewFactory()

	// Test struct validation
	t.Run("StructValidation", func(t *testing.T) {
		type User struct {
			Name  string `validate:"required|string"`
			Email string `validate:"required|email"`
			Age   int    `validate:"required|integer|min:18"`
		}

		user := User{
			Name:  "John Doe",
			Email: "john@example.com", 
			Age:   25,
		}

		validator, err := factory.ValidateStruct(user)
		if err != nil {
			t.Errorf("ValidateStruct should not return error: %v", err)
		}

		if !validator.Passes() {
			t.Errorf("Validation should pass: %v", validator.Errors().All())
		}

		// Test invalid struct
		invalidUser := User{
			Name:  "",
			Email: "invalid-email",
			Age:   15,
		}

		invalidValidator, err := factory.ValidateStruct(invalidUser)
		if err != nil {
			t.Errorf("ValidateStruct should not return error even for invalid data: %v", err)
		}

		if invalidValidator.Passes() {
			t.Errorf("Validation should fail for invalid user")
		}
	})

	// Test IsImplicit methods
	t.Run("ImplicitRules", func(t *testing.T) {
		// Test basic type rules
		if !(&RequiredRule{}).IsImplicit() {
			t.Error("RequiredRule should be implicit")
		}

		if !(&AcceptedRule{}).IsImplicit() {
			t.Error("AcceptedRule should be implicit")
		}

		if !(&AcceptedIfRule{}).IsImplicit() {
			t.Error("AcceptedIfRule should be implicit")
		}

		if !(&DeclinedRule{}).IsImplicit() {
			t.Error("DeclinedRule should be implicit")
		}

		if !(&DeclinedIfRule{}).IsImplicit() {
			t.Error("DeclinedIfRule should be implicit")
		}

		if !(&FilledRule{}).IsImplicit() {
			t.Error("FilledRule should be implicit")
		}

		if !(&PresentRule{}).IsImplicit() {
			t.Error("PresentRule should be implicit")
		}

		if !(&ProhibitedRule{}).IsImplicit() {
			t.Error("ProhibitedRule should be implicit")
		}

		// Test relationship rules
		if !(&RequiredIfRule{}).IsImplicit() {
			t.Error("RequiredIfRule should be implicit")
		}

		if !(&RequiredUnlessRule{}).IsImplicit() {
			t.Error("RequiredUnlessRule should be implicit")
		}

		if !(&RequiredWithRule{}).IsImplicit() {
			t.Error("RequiredWithRule should be implicit")
		}

		if !(&RequiredWithoutRule{}).IsImplicit() {
			t.Error("RequiredWithoutRule should be implicit")
		}

		if !(&RequiredWithAllRule{}).IsImplicit() {
			t.Error("RequiredWithAllRule should be implicit")
		}

		if !(&RequiredIfAcceptedRule{}).IsImplicit() {
			t.Error("RequiredIfAcceptedRule should be implicit")
		}

		if !(&RequiredIfDeclinedRule{}).IsImplicit() {
			t.Error("RequiredIfDeclinedRule should be implicit")
		}

		// Test special rules
		if !(&MissingRule{}).IsImplicit() {
			t.Error("MissingRule should be implicit")
		}

		if !(&MissingIfRule{}).IsImplicit() {
			t.Error("MissingIfRule should be implicit")
		}

		if !(&MissingUnlessRule{}).IsImplicit() {
			t.Error("MissingUnlessRule should be implicit")
		}

		if !(&MissingWithRule{}).IsImplicit() {
			t.Error("MissingWithRule should be implicit")
		}

		if !(&MissingWithAllRule{}).IsImplicit() {
			t.Error("MissingWithAllRule should be implicit")
		}

		if !(&PresentIfRule{}).IsImplicit() {
			t.Error("PresentIfRule should be implicit")
		}

		if !(&PresentUnlessRule{}).IsImplicit() {
			t.Error("PresentUnlessRule should be implicit")
		}

		if !(&PresentWithRule{}).IsImplicit() {
			t.Error("PresentWithRule should be implicit")
		}

		if !(&PresentWithAllRule{}).IsImplicit() {
			t.Error("PresentWithAllRule should be implicit")
		}

		if !(&ProhibitedIfAcceptedRule{}).IsImplicit() {
			t.Error("ProhibitedIfAcceptedRule should be implicit")
		}

		if !(&ProhibitedIfDeclinedRule{}).IsImplicit() {
			t.Error("ProhibitedIfDeclinedRule should be implicit")
		}

		if !(&ProhibitedUnlessRule{}).IsImplicit() {
			t.Error("ProhibitedUnlessRule should be implicit")
		}

		if !(&ProhibitsRule{}).IsImplicit() {
			t.Error("ProhibitsRule should be implicit")
		}
	})

	// Test date rules
	t.Run("DateRules", func(t *testing.T) {
		now := time.Now()
		yesterday := now.AddDate(0, 0, -1)
		tomorrow := now.AddDate(0, 0, 1)

		data := map[string]interface{}{
			"today":     now.Format("2006-01-02"),
			"yesterday": yesterday.Format("2006-01-02"),
			"tomorrow":  tomorrow.Format("2006-01-02"),
		}

		// Test BeforeRule
		beforeRule := &BeforeRule{Date: "tomorrow"}
		beforeRule.SetData(data)
		if !beforeRule.Passes("today", now.Format("2006-01-02")) {
			t.Error("BeforeRule should pass when date is before tomorrow")
		}
		if beforeRule.Message() == "" {
			t.Error("BeforeRule should have a message")
		}

		// Test AfterRule
		afterRule := &AfterRule{Date: "yesterday"}
		afterRule.SetData(data)
		if !afterRule.Passes("today", now.Format("2006-01-02")) {
			t.Error("AfterRule should pass when date is after yesterday")
		}
		if afterRule.Message() == "" {
			t.Error("AfterRule should have a message")
		}

		// Test AfterOrEqualRule
		afterOrEqualRule := &AfterOrEqualRule{Date: "today"}
		afterOrEqualRule.SetData(data)
		if !afterOrEqualRule.Passes("today", now.Format("2006-01-02")) {
			t.Error("AfterOrEqualRule should pass when date is equal")
		}
		if afterOrEqualRule.Message() == "" {
			t.Error("AfterOrEqualRule should have a message")
		}

		// Test BeforeOrEqualRule
		beforeOrEqualRule := &BeforeOrEqualRule{Date: "today"}
		beforeOrEqualRule.SetData(data)
		if !beforeOrEqualRule.Passes("today", now.Format("2006-01-02")) {
			t.Error("BeforeOrEqualRule should pass when date is equal")
		}
		if beforeOrEqualRule.Message() == "" {
			t.Error("BeforeOrEqualRule should have a message")
		}

		// Test DateEqualsRule
		dateEqualsRule := &DateEqualsRule{Date: "today"}
		dateEqualsRule.SetData(data)
		if !dateEqualsRule.Passes("today", now.Format("2006-01-02")) {
			t.Error("DateEqualsRule should pass when dates are equal")
		}
		if dateEqualsRule.Message() == "" {
			t.Error("DateEqualsRule should have a message")
		}
	})

	// Test special/missing rules
	t.Run("SpecialRules", func(t *testing.T) {
		data := map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
			"accepted_field": "yes",
			"declined_field": "no",
		}

		validator := factory.Make(data, map[string]interface{}{})

		// Test NullableRule and SometimesRule
		nullableRule := &NullableRule{}
		if !nullableRule.Passes("field", nil) {
			t.Error("NullableRule should always pass")
		}
		if nullableRule.Message() != "" {
			t.Error("NullableRule should have empty message")
		}

		sometimesRule := &SometimesRule{}
		if !sometimesRule.Passes("field", "value") {
			t.Error("SometimesRule should always pass")
		}
		if sometimesRule.Message() != "" {
			t.Error("SometimesRule should have empty message")
		}

		// Test MissingIfRule
		missingIfRule := &MissingIfRule{Field: "field1", Value: "value1"}
		missingIfRule.SetData(data)
		missingIfRule.SetValidator(validator)
		if !missingIfRule.Passes("missing_field", "any_value") {
			t.Error("MissingIfRule should pass when field is not actually missing from validator")
		}

		// Test MissingUnlessRule
		missingUnlessRule := &MissingUnlessRule{Field: "field1", Value: "different_value"}
		missingUnlessRule.SetData(data)
		missingUnlessRule.SetValidator(validator)
		if missingUnlessRule.Passes("field1", "value1") {
			t.Error("MissingUnlessRule should fail when condition is not met and field is present")
		}

		// Test MissingWithRule
		missingWithRule := &MissingWithRule{Fields: []string{"field1"}}
		missingWithRule.SetData(data)
		missingWithRule.SetValidator(validator)
		if missingWithRule.Passes("field1", "value1") {
			t.Error("MissingWithRule should fail when trigger field is present and field is present")
		}

		// Test MissingWithAllRule
		missingWithAllRule := &MissingWithAllRule{Fields: []string{"field1", "field2"}}
		missingWithAllRule.SetData(data)
		missingWithAllRule.SetValidator(validator)
		if missingWithAllRule.Passes("field1", "value1") {
			t.Error("MissingWithAllRule should fail when all trigger fields are present and field is present")
		}

		// Test PresentIfRule
		presentIfRule := &PresentIfRule{Field: "field1", Value: "value1"}
		presentIfRule.SetData(data)
		presentIfRule.SetValidator(validator)
		if !presentIfRule.Passes("field1", "value1") {
			t.Error("PresentIfRule should pass when condition is met and field is present")
		}

		// Test PresentUnlessRule
		presentUnlessRule := &PresentUnlessRule{Field: "field1", Value: "different_value"}
		presentUnlessRule.SetData(data)
		presentUnlessRule.SetValidator(validator)
		if !presentUnlessRule.Passes("field1", "value1") {
			t.Error("PresentUnlessRule should pass when condition is not met and field is present")
		}

		// Test PresentWithRule
		presentWithRule := &PresentWithRule{Fields: []string{"field1"}}
		presentWithRule.SetData(data)
		presentWithRule.SetValidator(validator)
		if !presentWithRule.Passes("field1", "value1") {
			t.Error("PresentWithRule should pass when trigger field is present and field is present")
		}

		// Test PresentWithAllRule
		presentWithAllRule := &PresentWithAllRule{Fields: []string{"field1", "field2"}}
		presentWithAllRule.SetData(data)
		presentWithAllRule.SetValidator(validator)
		if !presentWithAllRule.Passes("field1", "value1") {
			t.Error("PresentWithAllRule should pass when all trigger fields are present and field is present")
		}

		// Test ProhibitedIfAcceptedRule
		prohibitedIfAcceptedRule := &ProhibitedIfAcceptedRule{Field: "accepted_field"}
		prohibitedIfAcceptedRule.SetData(data)
		if prohibitedIfAcceptedRule.Passes("prohibited_field", "any_value") {
			t.Error("ProhibitedIfAcceptedRule should fail when trigger field is accepted and field has value")
		}

		// Test ProhibitedIfDeclinedRule
		prohibitedIfDeclinedRule := &ProhibitedIfDeclinedRule{Field: "declined_field"}
		prohibitedIfDeclinedRule.SetData(data)
		if prohibitedIfDeclinedRule.Passes("prohibited_field", "any_value") {
			t.Error("ProhibitedIfDeclinedRule should fail when trigger field is declined and field has value")
		}

		// Test ProhibitedUnlessRule
		prohibitedUnlessRule := &ProhibitedUnlessRule{Field: "field1", Value: "different_value"}
		prohibitedUnlessRule.SetData(data)
		if prohibitedUnlessRule.Passes("prohibited_field", "any_value") {
			t.Error("ProhibitedUnlessRule should fail when condition is not met and field has value")
		}

		// Test ProhibitsRule
		prohibitsRule := &ProhibitsRule{Fields: []string{"field2"}}
		prohibitsRule.SetData(data)
		if prohibitsRule.Passes("field1", "value1") {
			t.Error("ProhibitsRule should fail when prohibited field is present")
		}

		// Test RequiredWithoutRule  
		requiredWithoutRule := &RequiredWithoutRule{Fields: []string{"missing_field"}}
		requiredWithoutRule.SetData(data)
		if requiredWithoutRule.Passes("required_field", "") {
			t.Error("RequiredWithoutRule should fail when trigger field is missing and required field is empty")
		}

		// Test RequiredIfAcceptedRule
		requiredIfAcceptedRule := &RequiredIfAcceptedRule{Field: "accepted_field"}
		requiredIfAcceptedRule.SetData(data)
		if requiredIfAcceptedRule.Passes("required_field", "") {
			t.Error("RequiredIfAcceptedRule should fail when trigger field is accepted and required field is empty")
		}

		// Test RequiredIfDeclinedRule
		requiredIfDeclinedRule := &RequiredIfDeclinedRule{Field: "declined_field"}
		requiredIfDeclinedRule.SetData(data)
		if requiredIfDeclinedRule.Passes("required_field", "") {
			t.Error("RequiredIfDeclinedRule should fail when trigger field is declined and required field is empty")
		}

		// Test RequiredArrayKeysRule
		arrayData := map[string]interface{}{
			"array_field": map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		}
		requiredArrayKeysRule := &RequiredArrayKeysRule{Keys: []string{"key1", "key2"}}
		requiredArrayKeysRule.SetData(arrayData)
		if !requiredArrayKeysRule.Passes("array_field", arrayData["array_field"]) {
			t.Error("RequiredArrayKeysRule should pass when all required keys are present")
		}

		// Test with missing key
		missingKeyData := map[string]interface{}{
			"key1": "value1",
		}
		if requiredArrayKeysRule.Passes("array_field", missingKeyData) {
			t.Error("RequiredArrayKeysRule should fail when required key is missing")
		}
	})

	// Test utility functions with edge cases
	t.Run("UtilityFunctions", func(t *testing.T) {
		// Test comparison rules
		data := map[string]interface{}{
			"number1": "10",
			"number2": "20",
		}

		gtRule := &GreaterThanRule{Field: "number1"}
		gtRule.SetData(data)
		if !gtRule.Passes("number2", "20") {
			t.Error("GreaterThanRule should pass when value is greater")
		}

		gteRule := &GreaterThanOrEqualRule{Field: "number1"}
		gteRule.SetData(data)
		if !gteRule.Passes("number1", "10") {
			t.Error("GreaterThanOrEqualRule should pass when values are equal")
		}

		ltRule := &LessThanRule{Field: "number2"}
		ltRule.SetData(data)
		if !ltRule.Passes("number1", "10") {
			t.Error("LessThanRule should pass when value is less")
		}

		lteRule := &LessThanOrEqualRule{Field: "number1"}
		lteRule.SetData(data)
		if !lteRule.Passes("number1", "10") {
			t.Error("LessThanOrEqualRule should pass when values are equal")
		}
	})

	// Test message methods that weren't covered
	t.Run("MessageMethods", func(t *testing.T) {
		// Test various message methods
		if (&NumericRule{}).Message() == "" {
			t.Error("NumericRule should have a message")
		}
		
		if (&JsonRule{}).Message() == "" {
			t.Error("JsonRule should have a message")
		}

		if (&DateFormatRule{}).Message() == "" {
			t.Error("DateFormatRule should have a message")
		}

		if (&IpRule{}).Message() == "" {
			t.Error("IpRule should have a message")
		}

		if (&MacAddressRule{}).Message() == "" {
			t.Error("MacAddressRule should have a message")
		}

		if (&HexColorRule{}).Message() == "" {
			t.Error("HexColorRule should have a message")
		}

		if (&DigitsRule{Length: 5}).Message() == "" {
			t.Error("DigitsRule should have a message")
		}
	})

	// Test additional uncovered functionality
	t.Run("AdditionalCoverage", func(t *testing.T) {
		// Test ULID validation
		ulidRule := &UlidRule{}
		if !ulidRule.Passes("ulid", "01ARZ3NDEKTSV4RRFFQ69G5FAV") {
			t.Error("UlidRule should pass with valid ULID")
		}
		if ulidRule.Passes("ulid", "invalid-ulid") {
			t.Error("UlidRule should fail with invalid ULID")
		}

		// Test UUID validation with versions
		uuidRule := &UuidRule{Version: 4}
		// Test with a valid UUID v4
		if !uuidRule.Passes("uuid", "f47ac10b-58cc-4372-a567-0e02b2c3d479") {
			t.Error("UuidRule should pass with valid UUID")
		}

		// Test UUID with max version
		uuidMaxRule := &UuidRule{Version: "max"}
		if !uuidMaxRule.Passes("uuid", "f47ac10b-58cc-4372-a567-0e02b2c3d479") {
			t.Error("UuidRule with max version should pass with valid UUID")
		}

		// Test edge cases for date format conversion
		dateFormatRule := &DateFormatRule{Formats: []string{"Y-m-d", "d/m/Y"}}
		// Test different PHP date format conversions
		if !dateFormatRule.Passes("date", "2023-01-15") {
			t.Error("DateFormatRule should pass with Y-m-d format")
		}
		if !dateFormatRule.Passes("date", "15/01/2023") {
			t.Error("DateFormatRule should pass with d/m/Y format")
		}

		// Test error conditions for date rules
		afterRule := &AfterRule{Date: "invalid-date"}
		afterRule.SetData(map[string]interface{}{})
		if afterRule.Passes("date", "2023-01-01") {
			t.Error("AfterRule should fail with invalid comparison date")
		}

		// Test RequiredWithAll with empty trigger
		emptyData := map[string]interface{}{}
		requiredWithAllRule := &RequiredWithAllRule{Fields: []string{"missing1", "missing2"}}
		requiredWithAllRule.SetData(emptyData)
		validator2 := factory.Make(emptyData, map[string]interface{}{})
		if !requiredWithAllRule.Passes("optional_field", "") {
			t.Error("RequiredWithAllRule should pass when trigger fields are missing")
		}

		// Test invalid struct validation
		type InvalidStruct struct {
			Name string `validate:"invalid-rule"`
		}
		invalidStruct := InvalidStruct{Name: "test"}
		_, err := factory.ValidateStruct(invalidStruct)
		// Should not error even with invalid rule - it would just be ignored
		if err != nil {
			t.Errorf("ValidateStruct should handle invalid rules gracefully: %v", err)
		}

		// Test various rule messages for coverage
		rules := []Rule{
			&RequiredWithoutRule{},
			&RequiredIfAcceptedRule{},
			&RequiredIfDeclinedRule{},
			&RequiredArrayKeysRule{},
			&MissingIfRule{},
			&MissingUnlessRule{},
			&MissingWithRule{},
			&MissingWithAllRule{},
			&PresentIfRule{},
			&PresentUnlessRule{},
			&PresentWithRule{},
			&PresentWithAllRule{},
			&ProhibitedIfAcceptedRule{},
			&ProhibitedIfDeclinedRule{},
			&ProhibitedUnlessRule{},
			&ProhibitsRule{},
			&GreaterThanRule{},
			&GreaterThanOrEqualRule{},
			&LessThanRule{},
			&LessThanOrEqualRule{},
		}

		for _, rule := range rules {
			if rule.Message() == "" {
				t.Errorf("Rule %T should have a non-empty message", rule)
			}
		}

		// Test some SetData and SetValidator methods for full coverage
		dataAwareRules := []interface{}{
			&MissingIfRule{},
			&MissingUnlessRule{},
			&MissingWithRule{},
			&MissingWithAllRule{},
			&PresentIfRule{},
			&PresentUnlessRule{},
			&PresentWithRule{},
			&PresentWithAllRule{},
		}

		for _, rule := range dataAwareRules {
			if dataRule, ok := rule.(interface{ SetData(map[string]interface{}) }); ok {
				dataRule.SetData(map[string]interface{}{})
			}
			if validatorRule, ok := rule.(interface{ SetValidator(Validator) }); ok {
				validatorRule.SetValidator(validator2)
			}
		}
	})
}