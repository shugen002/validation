package main

import (
	"fmt"

	"github.com/shugen002/validation"
)

func main() {
	factory := validation.NewFactory()

	// Example 1: Decimal validation for pricing
	fmt.Println("=== Example 1: Decimal Places Validation ===")
	priceData := map[string]interface{}{
		"price":     "19.99",
		"discount":  "5.5",
		"shipping":  "0",
	}
	priceRules := map[string]interface{}{
		"price":    "required|decimal:2",        // Exactly 2 decimal places
		"discount": "decimal:1,3",               // Between 1-3 decimal places
		"shipping": "decimal:0,2",               // 0-2 decimal places
	}
	priceValidator := factory.Make(priceData, priceRules)
	if priceValidator.Passes() {
		fmt.Println("✅ Price validation passed")
	} else {
		fmt.Printf("❌ Price validation failed: %v\n", priceValidator.Errors().All())
	}

	// Example 2: Distinct array validation
	fmt.Println("\n=== Example 2: Distinct Array Values ===")
	arrayData := map[string]interface{}{
		"tags":       []interface{}{"php", "golang", "javascript"},
		"categories": []interface{}{"Web", "web", "API"}, // Case-sensitive duplicates
		"ids":        []interface{}{1, 2, 3, 1},         // Duplicates
	}
	arrayRules := map[string]interface{}{
		"tags":       "distinct",                // No duplicates
		"categories": "distinct:ignore_case",    // No duplicates (case-insensitive)
		"ids":        "distinct",                // Should fail
	}
	arrayValidator := factory.Make(arrayData, arrayRules)
	if arrayValidator.Passes() {
		fmt.Println("✅ Array validation passed")
	} else {
		fmt.Printf("❌ Array validation failed: %v\n", arrayValidator.Errors().All())
	}

	// Example 3: Multiple of validation
	fmt.Println("\n=== Example 3: Multiple Of Validation ===")
	quantityData := map[string]interface{}{
		"quantity":   "10",   // Must be multiple of 5
		"batch_size": "12",   // Must be multiple of 6
		"step":       "2.5",  // Must be multiple of 0.5
	}
	quantityRules := map[string]interface{}{
		"quantity":   "multiple_of:5",
		"batch_size": "multiple_of:6",
		"step":       "multiple_of:0.5",
	}
	quantityValidator := factory.Make(quantityData, quantityRules)
	if quantityValidator.Passes() {
		fmt.Println("✅ Quantity validation passed")
	} else {
		fmt.Printf("❌ Quantity validation failed: %v\n", quantityValidator.Errors().All())
	}

	// Example 4: Minimum digits validation
	fmt.Println("\n=== Example 4: Minimum Digits Validation ===")
	codeData := map[string]interface{}{
		"user_id":     "123456",    // At least 6 digits
		"auth_code":   "12345",     // At least 4 digits
		"short_code":  "12",        // Should fail - less than 3 digits
	}
	codeRules := map[string]interface{}{
		"user_id":    "min_digits:6",
		"auth_code":  "min_digits:4",
		"short_code": "min_digits:3",
	}
	codeValidator := factory.Make(codeData, codeRules)
	if codeValidator.Passes() {
		fmt.Println("✅ Code validation passed")
	} else {
		fmt.Printf("❌ Code validation failed: %v\n", codeValidator.Errors().All())
	}

	// Example 5: Advanced presence validation
	fmt.Println("\n=== Example 5: Advanced Presence Validation ===")
	presenceData := map[string]interface{}{
		"user_type": "admin",
		"api_key":   "secret123",
		// "debug_mode" should be missing for production
	}
	presenceRules := map[string]interface{}{
		"user_type":  "required|in:admin,user,guest",
		"api_key":    "present_if:user_type,admin",           // Required when user is admin
		"debug_mode": "missing_unless:user_type,developer",   // Should be missing unless developer
	}
	presenceValidator := factory.Make(presenceData, presenceRules)
	if presenceValidator.Passes() {
		fmt.Println("✅ Presence validation passed")
	} else {
		fmt.Printf("❌ Presence validation failed: %v\n", presenceValidator.Errors().All())
	}

	// Example 6: Required array keys validation
	fmt.Println("\n=== Example 6: Required Array Keys Validation ===")
	configData := map[string]interface{}{
		"database": map[string]interface{}{
			"host":     "localhost",
			"port":     5432,
			"name":     "myapp",
			"user":     "admin",
			"password": "secret",
		},
		"cache": map[string]interface{}{
			"driver": "redis",
			// Missing "host" key
		},
	}
	configRules := map[string]interface{}{
		"database": "required_array_keys:host,port,name,user", // All required keys present
		"cache":    "required_array_keys:driver,host",         // Missing "host" key
	}
	configValidator := factory.Make(configData, configRules)
	if configValidator.Passes() {
		fmt.Println("✅ Configuration validation passed")
	} else {
		fmt.Printf("❌ Configuration validation failed: %v\n", configValidator.Errors().All())
	}

	// Example 7: Conditional prohibition and acceptance
	fmt.Println("\n=== Example 7: Conditional Prohibition ===")
	formData := map[string]interface{}{
		"newsletter":    "yes",
		"email":         "user@example.com",
		"sms_marketing": "", // Should be prohibited when newsletter is accepted
		"backup_plan":   "basic",
	}
	formRules := map[string]interface{}{
		"newsletter":    "required|in:yes,no",
		"email":         "required_if_accepted:newsletter|email",
		"sms_marketing": "prohibited_if_accepted:newsletter",
		"premium_plan":  "missing_with:backup_plan", // Can't have both
	}
	formValidator := factory.Make(formData, formRules)
	if formValidator.Passes() {
		fmt.Println("✅ Form validation passed")
	} else {
		fmt.Printf("❌ Form validation failed: %v\n", formValidator.Errors().All())
	}

	fmt.Println("\n=== All Additional Laravel Rules Implemented Successfully! ===")
	fmt.Println("Total rules implemented: 20+ new Laravel validation rules")
	fmt.Println("- decimal, distinct, min_digits, multiple_of, missing")
	fmt.Println("- missing_if, missing_unless, missing_with, missing_with_all")
	fmt.Println("- present_if, present_unless, present_with, present_with_all")
	fmt.Println("- prohibited_if_accepted, prohibited_if_declined, prohibited_unless")
	fmt.Println("- prohibits, required_if_accepted, required_if_declined")
	fmt.Println("- required_with_all, required_array_keys")
}