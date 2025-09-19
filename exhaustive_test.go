package validation

import (
	"testing"
)

// TestExhaustiveCoverage tests every remaining uncovered line
func TestExhaustiveCoverage(t *testing.T) {
	factory := NewFactory()

	// Test all timezone validation
	t.Run("TimezoneValidation", func(t *testing.T) {
		timezoneRule := &TimezoneRule{}
		
		// Test valid timezone
		if !timezoneRule.Passes("tz", "America/New_York") {
			t.Error("TimezoneRule should pass with valid timezone")
		}
		
		// Test invalid timezone
		if timezoneRule.Passes("tz", "Invalid/Timezone") {
			t.Error("TimezoneRule should fail with invalid timezone")
		}
		
		// Test empty timezone
		if timezoneRule.Passes("tz", "") {
			t.Error("TimezoneRule should fail with empty timezone")
		}

		if timezoneRule.Message() == "" {
			t.Error("TimezoneRule should have a message")
		}
	})

	// Test all remaining factory functionality
	t.Run("FactoryCompleteCoverage", func(t *testing.T) {
		// Test struct parsing with complex scenarios
		type ComplexStruct struct {
			Name        string  `validate:"required|string|min:2"`
			Email       string  `validate:"required|email"`
			Age         int     `validate:"required|integer|min:18|max:120"`
			Score       float64 `validate:"required|numeric|between:0,100"`
			Active      bool    `validate:"required|boolean"`
			Tags        []string `validate:"array"`
			Meta        map[string]string `validate:"required"`
			OptionalTag string  `validate:"sometimes|string"`
		}

		complex := ComplexStruct{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   30,
			Score: 85.5,
			Active: true,
			Tags:  []string{"tag1", "tag2"},
			Meta:  map[string]string{"key": "value"},
		}

		validator, err := factory.ValidateStruct(complex)
		if err != nil {
			t.Errorf("ValidateStruct should handle complex struct: %v", err)
		}
		if !validator.Passes() {
			t.Errorf("Complex struct validation should pass: %v", validator.Errors().All())
		}

		// Test with invalid complex struct
		invalidComplex := ComplexStruct{
			Name:  "J", // too short
			Email: "invalid-email",
			Age:   15, // too young
			Score: 150, // too high
			Active: true,
			Tags:  []string{},
			Meta:  nil, // required but nil
		}

		invalidValidator, err := factory.ValidateStruct(invalidComplex)
		if err != nil {
			t.Errorf("ValidateStruct should not error even with invalid data: %v", err)
		}
		if invalidValidator.Passes() {
			t.Error("Invalid complex struct should fail validation")
		}

		// Verify specific error messages exist
		errors := invalidValidator.Errors().All()
		if len(errors) == 0 {
			t.Error("Should have validation errors")
		}
	})

	// Test all remaining rule edge cases
	t.Run("RuleEdgeCases", func(t *testing.T) {
		// Test all message methods that might not be covered
		rules := []Rule{
			&BooleanRule{},
			&ArrayRule{AllowedKeys: []string{"key1", "key2"}},
			&FilledRule{},
			&PresentRule{},
			&ProhibitedRule{},
		}

		for _, rule := range rules {
			if rule.Message() == "" {
				t.Errorf("Rule %T should have a non-empty message", rule)
			}
		}

		// Test array rule with allowed keys
		arrayRule := &ArrayRule{AllowedKeys: []string{"key1", "key2"}}
		testMap := map[string]interface{}{"key1": "value1", "key2": "value2"}
		if !arrayRule.Passes("field", testMap) {
			t.Error("ArrayRule should pass with valid keys")
		}

		// Test with invalid key
		invalidMap := map[string]interface{}{"key1": "value1", "invalid_key": "value2"}
		if arrayRule.Passes("field", invalidMap) {
			t.Error("ArrayRule should fail with invalid key")
		}

		// Test boolean rule strict mode
		boolRuleStrict := &BooleanRule{Strict: true}
		if boolRuleStrict.Passes("field", "true") {
			t.Error("Boolean rule in strict mode should fail with string")
		}
		if !boolRuleStrict.Passes("field", true) {
			t.Error("Boolean rule in strict mode should pass with boolean")
		}

		// Test integer rule strict mode
		intRuleStrict := &IntegerRule{Strict: true}
		if intRuleStrict.Passes("field", "123") {
			t.Error("Integer rule in strict mode should fail with string")
		}
		if !intRuleStrict.Passes("field", 123) {
			t.Error("Integer rule in strict mode should pass with integer")
		}

		// Test numeric rule strict mode
		numRuleStrict := &NumericRule{Strict: true}
		if numRuleStrict.Passes("field", "123.45") {
			t.Error("Numeric rule in strict mode should fail with string")
		}
		if !numRuleStrict.Passes("field", 123.45) {
			t.Error("Numeric rule in strict mode should pass with float")
		}
	})

	// Test SetValidator methods
	t.Run("SetValidatorMethods", func(t *testing.T) {
		data := map[string]interface{}{
			"field1": "value1",
		}
		validator := factory.Make(data, map[string]interface{}{})

		// Test all rules that implement SetValidator
		rules := []interface{}{
			&FilledRule{},
			&PresentRule{},
			&MissingRule{},
			&MissingIfRule{},
			&MissingUnlessRule{},
			&MissingWithRule{},
			&MissingWithAllRule{},
			&PresentIfRule{},
			&PresentUnlessRule{},
			&PresentWithRule{},
			&PresentWithAllRule{},
		}

		for _, rule := range rules {
			if setValidatorRule, ok := rule.(interface{ SetValidator(Validator) }); ok {
				setValidatorRule.SetValidator(validator)
			}
		}

		// Actually test some of them to ensure they work
		missingRule := &MissingRule{}
		missingRule.SetValidator(validator)
		if missingRule.Passes("field1", "value") {
			t.Error("MissingRule should fail when field is present in validator")
		}
		if !missingRule.Passes("missing_field", "value") {
			t.Error("MissingRule should pass when field is not present in validator")
		}

		filledRule := &FilledRule{}
		filledRule.SetValidator(validator)
		if !filledRule.Passes("missing_field", "") {
			t.Error("FilledRule should pass when field is not in validator")
		}
		if filledRule.Passes("field1", "") {
			t.Error("FilledRule should fail when field is present but empty")
		}

		presentRule := &PresentRule{}
		presentRule.SetValidator(validator)
		if !presentRule.Passes("field1", "value") {
			t.Error("PresentRule should pass when field is present in validator")
		}
		if presentRule.Passes("missing_field", "value") {
			t.Error("PresentRule should fail when field is not present in validator")
		}
	})

	// Test additional validation scenarios
	t.Run("AdditionalValidationScenarios", func(t *testing.T) {
		// Test all rule types with various data scenarios
		testData := map[string]interface{}{
			"string_field":  "test_value",
			"number_field":  42,
			"float_field":   3.14,
			"bool_field":    true,
			"array_field":   []string{"item1", "item2"},
			"map_field":     map[string]string{"key": "value"},
			"empty_string":  "",
			"zero_number":   0,
			"false_bool":    false,
			"nil_field":     nil,
		}

		// Test each rule type with different data types
		rules := map[string]Rule{
			"string":    &StringRule{},
			"integer":   &IntegerRule{},
			"numeric":   &NumericRule{},
			"boolean":   &BooleanRule{},
			"array":     &ArrayRule{},
			"required":  &RequiredRule{},
			"nullable":  &NullableRule{},
			"sometimes": &SometimesRule{},
			"bail":      &BailRule{},
		}

		for ruleName, rule := range rules {
			for fieldName, value := range testData {
				// Just exercise the code paths
				result := rule.Passes(fieldName, value)
				_ = result // Use the result to avoid unused variable
				
				message := rule.Message()
				if ruleName != "nullable" && ruleName != "sometimes" && ruleName != "bail" && message == "" {
					t.Errorf("Rule %s should have a non-empty message", ruleName)
				}
			}
		}
	})

	// Test validator error handling edge cases
	t.Run("ValidatorErrorHandling", func(t *testing.T) {
		// Test validator with complex rule combinations
		data := map[string]interface{}{
			"name":  "John",
			"email": "john@example.com",
			"age":   25,
		}

		rules := map[string]interface{}{
			"name":  "required|string|min:2|max:50",
			"email": "required|email",
			"age":   "required|integer|between:18,65",
			"optional": "sometimes|string",
		}

		validator := factory.Make(data, rules)
		
		// Test all validator methods
		if !validator.Passes() {
			t.Error("Validator should pass with valid data")
		}
		
		if validator.Fails() {
			t.Error("Validator should not fail with valid data")
		}

		valid := validator.Valid()
		if len(valid) != 3 { // name, email, age
			t.Errorf("Valid() should return 3 fields, got %d", len(valid))
		}

		invalid := validator.Invalid()
		if len(invalid) != 0 {
			t.Errorf("Invalid() should return 0 fields for valid data, got %d", len(invalid))
		}

		// Test with invalid data
		invalidData := map[string]interface{}{
			"name":  "J", // too short
			"email": "invalid",
			"age":   15, // too young
		}

		invalidValidator := factory.Make(invalidData, rules)
		
		if invalidValidator.Passes() {
			t.Error("Validator should fail with invalid data")
		}
		
		if !invalidValidator.Fails() {
			t.Error("Validator should fail with invalid data")
		}

		invalidFields := invalidValidator.Invalid()
		if len(invalidFields) == 0 {
			t.Error("Invalid() should return failed fields")
		}

		validFields := invalidValidator.Valid()
		// Should be empty since all fields failed
		if len(validFields) != 0 {
			t.Errorf("Valid() should return 0 fields for invalid data, got %d", len(validFields))
		}

		// Test error messages
		errors := invalidValidator.Errors()
		if errors.IsEmpty() {
			t.Error("Should have validation errors")
		}
		
		if !errors.IsNotEmpty() {
			t.Error("Errors should not be empty")
		}

		if errors.Count() == 0 {
			t.Error("Error count should be greater than 0")
		}

		// Test HasField
		if !invalidValidator.HasField("name") {
			t.Error("HasField should return true for existing field")
		}
		
		if invalidValidator.HasField("nonexistent") {
			t.Error("HasField should return false for non-existing field")
		}
	})
}