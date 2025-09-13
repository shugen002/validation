package main

import (
	"fmt"
	"log"

	"github.com/shugen002/validation"
)

// User struct for demonstration
type User struct {
	Name     string `json:"name" validate:"required|string|min:2"`
	Email    string `json:"email" validate:"required|email"`
	Age      int    `json:"age" validate:"integer|min:18|max:100"`
	Website  string `json:"website" validate:"url"`
	Password string `json:"password" validate:"required|min:8"`
	PasswordConfirmation string `json:"password_confirmation" validate:"same:password"`
}

func main() {
	// Example 1: Basic validation using map
	fmt.Println("=== Example 1: Basic Validation ===")
	
	factory := validation.NewFactory()
	
	data := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   25,
	}
	
	rules := map[string]interface{}{
		"name":  "required|string|min:2",
		"email": "required|email",
		"age":   "integer|min:18|max:100",
	}
	
	validator := factory.Make(data, rules)
	
	if validator.Passes() {
		fmt.Println("✓ Validation passed!")
		fmt.Printf("Valid data: %+v\n", validator.Valid())
	} else {
		fmt.Println("✗ Validation failed!")
		for field, errors := range validator.Errors().All() {
			fmt.Printf("  %s: %v\n", field, errors)
		}
	}
	
	fmt.Println()
	
	// Example 2: Validation with custom messages
	fmt.Println("=== Example 2: Custom Messages ===")
	
	invalidData := map[string]interface{}{
		"name":  "",
		"email": "invalid-email",
		"age":   15,
	}
	
	customMessages := map[string]string{
		"name.required":  "Please provide your name",
		"email.email":    "Please provide a valid email address",
		"age.min":        "You must be at least 18 years old",
	}
	
	validator2 := factory.Make(invalidData, rules, customMessages)
	
	if validator2.Fails() {
		fmt.Println("Validation errors with custom messages:")
		for field, errors := range validator2.Errors().All() {
			fmt.Printf("  %s: %s\n", field, errors[0])
		}
	}
	
	fmt.Println()
	
	// Example 3: Field relationships
	fmt.Println("=== Example 3: Field Relationships ===")
	
	passwordData := map[string]interface{}{
		"password":              "secret123",
		"password_confirmation": "secret123",
	}
	
	passwordRules := map[string]interface{}{
		"password":              "required|min:8",
		"password_confirmation": "required|same:password",
	}
	
	validator3 := factory.Make(passwordData, passwordRules)
	
	if validator3.Passes() {
		fmt.Println("✓ Password confirmation matches!")
	} else {
		fmt.Println("✗ Password validation failed!")
		for field, errors := range validator3.Errors().All() {
			fmt.Printf("  %s: %v\n", field, errors)
		}
	}
	
	fmt.Println()
	
	// Example 4: Network validation
	fmt.Println("=== Example 4: Network Validation ===")
	
	networkData := map[string]interface{}{
		"website":    "https://example.com",
		"ip_address": "192.168.1.1",
		"mac":        "00:1B:63:84:45:E6",
		"user_id":    "550e8400-e29b-41d4-a716-446655440000",
	}
	
	networkRules := map[string]interface{}{
		"website":    "url",
		"ip_address": "ip",
		"mac":        "mac_address",
		"user_id":    "uuid",
	}
	
	validator4 := factory.Make(networkData, networkRules)
	
	if validator4.Passes() {
		fmt.Println("✓ All network data is valid!")
	} else {
		fmt.Println("✗ Network validation failed!")
		for field, errors := range validator4.Errors().All() {
			fmt.Printf("  %s: %v\n", field, errors)
		}
	}
	
	fmt.Println()
	
	// Example 5: Date validation
	fmt.Println("=== Example 5: Date Validation ===")
	
	dateData := map[string]interface{}{
		"start_date": "2023-01-01",
		"end_date":   "2023-12-31",
		"timezone":   "America/New_York",
	}
	
	dateRules := map[string]interface{}{
		"start_date": "required|date",
		"end_date":   "required|date|after:start_date",
		"timezone":   "timezone",
	}
	
	validator5 := factory.Make(dateData, dateRules)
	
	if validator5.Passes() {
		fmt.Println("✓ All dates are valid!")
	} else {
		fmt.Println("✗ Date validation failed!")
		for field, errors := range validator5.Errors().All() {
			fmt.Printf("  %s: %v\n", field, errors)
		}
	}
	
	fmt.Println()
	
	// Example 6: Conditional validation
	fmt.Println("=== Example 6: Conditional Validation ===")
	
	conditionalData := map[string]interface{}{
		"contact_type": "email",
		"contact":      "user@example.com",
	}
	
	validator6 := factory.Make(conditionalData, map[string]interface{}{
		"contact_type": "required|in:email,phone",
	})
	
	// Add conditional validation for contact field
	validator6.Sometimes("contact", []validation.Rule{&validation.EmailRule{}}, func(data map[string]interface{}) bool {
		return data["contact_type"] == "email"
	})
	
	if validator6.Passes() {
		fmt.Println("✓ Conditional validation passed!")
	} else {
		fmt.Println("✗ Conditional validation failed!")
		for field, errors := range validator6.Errors().All() {
			fmt.Printf("  %s: %v\n", field, errors)
		}
	}
	
	fmt.Println()
	
	// Example 7: Struct validation with tags
	fmt.Println("=== Example 7: Struct Validation ===")
	
	user := User{
		Name:     "Jane Smith",
		Email:    "jane@example.com",
		Age:      28,
		Website:  "https://janesmith.com",
		Password: "securepass123",
		PasswordConfirmation: "securepass123",
	}
	
	validator7, err := factory.ValidateStruct(user)
	if err != nil {
		log.Fatal("Error creating validator:", err)
	}
	
	if validator7.Passes() {
		fmt.Println("✓ Struct validation passed!")
		fmt.Printf("Valid user data: %+v\n", validator7.Valid())
	} else {
		fmt.Println("✗ Struct validation failed!")
		for field, errors := range validator7.Errors().All() {
			fmt.Printf("  %s: %v\n", field, errors)
		}
	}
	
	fmt.Println()
	
	// Example 8: Stop on first failure
	fmt.Println("=== Example 8: Stop on First Failure ===")
	
	badData := map[string]interface{}{
		"field1": "",
		"field2": "",
		"field3": "",
	}
	
	badRules := map[string]interface{}{
		"field1": "required",
		"field2": "required", 
		"field3": "required",
	}
	
	validator8 := factory.Make(badData, badRules).StopOnFirstFailure()
	
	if validator8.Fails() {
		fmt.Printf("Validation stopped after first failure. Error count: %d\n", validator8.Errors().Count())
		fmt.Printf("First error: %s\n", validator8.Errors().First("field1"))
	}
	
	fmt.Println()
	
	// Example 9: Boolean validation rules
	fmt.Println("=== Example 9: Boolean Validation Rules ===")
	
	// Accepted rule example
	acceptData := map[string]interface{}{
		"terms_of_service": "yes",
		"privacy_policy":   true,
		"newsletter":       1,
	}
	
	acceptRules := map[string]interface{}{
		"terms_of_service": "required|accepted",
		"privacy_policy":   "required|accepted",
		"newsletter":       "accepted",
	}
	
	acceptValidator := factory.Make(acceptData, acceptRules)
	
	if acceptValidator.Passes() {
		fmt.Println("✓ All accepted fields validated successfully!")
	} else {
		fmt.Println("✗ Accepted validation failed!")
		for field, errors := range acceptValidator.Errors().All() {
			fmt.Printf("  %s: %v\n", field, errors)
		}
	}
	
	// Declined rule example
	declineData := map[string]interface{}{
		"spam_emails":     "no",
		"data_sharing":    false,
		"marketing_calls": 0,
	}
	
	declineRules := map[string]interface{}{
		"spam_emails":     "required|declined",
		"data_sharing":    "declined",
		"marketing_calls": "declined",
	}
	
	declineValidator := factory.Make(declineData, declineRules)
	
	if declineValidator.Passes() {
		fmt.Println("✓ All declined fields validated successfully!")
	} else {
		fmt.Println("✗ Declined validation failed!")
		for field, errors := range declineValidator.Errors().All() {
			fmt.Printf("  %s: %v\n", field, errors)
		}
	}
	
	// Conditional boolean validation
	conditionalBoolData := map[string]interface{}{
		"payment_method": "credit_card",
		"save_card":      "yes",
		"account_type":   "guest",
		"newsletter":     "no",
	}
	
	conditionalBoolRules := map[string]interface{}{
		"save_card":  "accepted_if:payment_method,credit_card",
		"newsletter": "declined_if:account_type,guest",
	}
	
	conditionalBoolValidator := factory.Make(conditionalBoolData, conditionalBoolRules)
	
	if conditionalBoolValidator.Passes() {
		fmt.Println("✓ Conditional boolean validation passed!")
	} else {
		fmt.Println("✗ Conditional boolean validation failed!")
		for field, errors := range conditionalBoolValidator.Errors().All() {
			fmt.Printf("  %s: %v\n", field, errors)
		}
	}
	
	fmt.Println()
	fmt.Println("=== All Examples Complete ===")
}