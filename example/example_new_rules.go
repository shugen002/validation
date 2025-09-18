package main

import (
	"fmt"

	"github.com/shugen002/validation"
)

func main() {
	factory := validation.NewFactory()

	// Example 1: Using digits rules for phone validation
	fmt.Println("=== Example 1: Phone Number Validation ===")
	phoneData := map[string]interface{}{
		"phone": "1234567890",
	}
	phoneRules := map[string]interface{}{
		"phone": "required|digits:10",
	}
	phoneValidator := factory.Make(phoneData, phoneRules)
	if phoneValidator.Passes() {
		fmt.Println("✅ Phone number is valid")
	} else {
		fmt.Printf("❌ Phone validation failed: %v\n", phoneValidator.Errors().All())
	}

	// Example 2: Using filled rule for optional fields
	fmt.Println("\n=== Example 2: Optional Field with Filled Rule ===")
	profileData := map[string]interface{}{
		"name":        "John Doe",
		"description": "", // Empty but present
	}
	profileRules := map[string]interface{}{
		"name":        "required|string",
		"description": "filled|min:10", // Must have content if present
	}
	profileValidator := factory.Make(profileData, profileRules)
	if !profileValidator.Passes() {
		fmt.Printf("❌ Profile validation failed: %v\n", profileValidator.Errors().All())
	}

	// Example 3: Using conditional required rules for payment
	fmt.Println("\n=== Example 3: Conditional Payment Validation ===")
	paymentData := map[string]interface{}{
		"payment_method": "card",
		"card_number":    "4111111111111111",
		// paypal_email not needed for card payment
	}
	paymentRules := map[string]interface{}{
		"payment_method": "required|in:card,paypal,cash",
		"card_number":    "required_if:payment_method,card|digits:16",
		"paypal_email":   "required_if:payment_method,paypal|email",
	}
	paymentValidator := factory.Make(paymentData, paymentRules)
	if paymentValidator.Passes() {
		fmt.Println("✅ Payment information is valid")
	} else {
		fmt.Printf("❌ Payment validation failed: %v\n", paymentValidator.Errors().All())
	}

	// Example 4: Using required_with for user registration
	fmt.Println("\n=== Example 4: User Registration with Dependencies ===")
	userData := map[string]interface{}{
		"username":              "johndoe",
		"password":              "secret123",
		"password_confirmation": "secret123", // Now provided
	}
	userRules := map[string]interface{}{
		"username":              "required|alpha_dash|min:3",
		"password":              "required|min:6",
		"password_confirmation": "required_with:password|same:password",
	}
	userValidator := factory.Make(userData, userRules)
	if userValidator.Passes() {
		fmt.Println("✅ User registration is valid")
	} else {
		fmt.Printf("❌ User registration failed: %v\n", userValidator.Errors().All())
	}

	// Example 5: Using present and prohibited rules
	fmt.Println("\n=== Example 5: API Request Validation ===")
	apiData := map[string]interface{}{
		"action": "update",
		"id":     "123",
		// "create_fields" should be prohibited for update action
	}
	apiRules := map[string]interface{}{
		"action":        "required|in:create,update,delete",
		"id":            "required_unless:action,create|integer",
		"create_fields": "prohibited", // Should not be present
	}
	apiValidator := factory.Make(apiData, apiRules)
	if apiValidator.Passes() {
		fmt.Println("✅ API request is valid")
	} else {
		fmt.Printf("❌ API validation failed: %v\n", apiValidator.Errors().All())
	}

	// Example 6: Discord bot configuration (like the real-world test case)
	fmt.Println("\n=== Example 6: Discord Bot Owner ID Validation ===")
	botData := map[string]interface{}{
		"bot_owner": "123456789012345678", // Valid 18-digit Discord ID
	}
	botRules := map[string]interface{}{
		"bot_owner": "required|digits_between:17,18",
	}
	botValidator := factory.Make(botData, botRules)
	if botValidator.Passes() {
		fmt.Println("✅ Discord bot configuration is valid")
	} else {
		fmt.Printf("❌ Bot validation failed: %v\n", botValidator.Errors().All())
	}

	fmt.Println("\n=== All New Laravel Rules Implemented Successfully! ===")
}