package validation

import (
	"strings"
	"testing"
)

func TestDateRules(t *testing.T) {
	factory := NewFactory()
	
	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "date rule passes with valid date",
			data: map[string]interface{}{
				"birthday": "1990-01-15",
			},
			rules: map[string]interface{}{
				"birthday": "date",
			},
			valid: true,
		},
		{
			name: "date rule fails with invalid date",
			data: map[string]interface{}{
				"birthday": "not-a-date",
			},
			rules: map[string]interface{}{
				"birthday": "date",
			},
			valid: false,
		},
		{
			name: "date_format rule passes",
			data: map[string]interface{}{
				"birthday": "15/01/1990",
			},
			rules: map[string]interface{}{
				"birthday": "date_format:d/m/Y",
			},
			valid: true,
		},
		{
			name: "after rule passes",
			data: map[string]interface{}{
				"start_date": "2023-01-01",
				"end_date":   "2023-12-31",
			},
			rules: map[string]interface{}{
				"end_date": "after:start_date",
			},
			valid: true,
		},
		{
			name: "after rule fails",
			data: map[string]interface{}{
				"start_date": "2023-12-31",
				"end_date":   "2023-01-01",
			},
			rules: map[string]interface{}{
				"end_date": "after:start_date",
			},
			valid: false,
		},
		{
			name: "before rule passes",
			data: map[string]interface{}{
				"start_date": "2023-01-01",
				"end_date":   "2023-12-31",
			},
			rules: map[string]interface{}{
				"start_date": "before:end_date",
			},
			valid: true,
		},
		{
			name: "timezone rule passes",
			data: map[string]interface{}{
				"timezone": "America/New_York",
			},
			rules: map[string]interface{}{
				"timezone": "timezone",
			},
			valid: true,
		},
		{
			name: "timezone rule fails",
			data: map[string]interface{}{
				"timezone": "Invalid/Timezone",
			},
			rules: map[string]interface{}{
				"timezone": "timezone",
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

func TestUuidAndOtherRules(t *testing.T) {
	factory := NewFactory()
	
	tests := []struct {
		name  string
		data  map[string]interface{}
		rules map[string]interface{}
		valid bool
	}{
		{
			name: "uuid rule passes with valid UUID",
			data: map[string]interface{}{
				"id": "550e8400-e29b-41d4-a716-446655440000",
			},
			rules: map[string]interface{}{
				"id": "uuid",
			},
			valid: true,
		},
		{
			name: "uuid rule fails with invalid UUID",
			data: map[string]interface{}{
				"id": "not-a-uuid",
			},
			rules: map[string]interface{}{
				"id": "uuid",
			},
			valid: false,
		},
		{
			name: "ulid rule passes with valid ULID",
			data: map[string]interface{}{
				"id": "01ARZ3NDEKTSV4RRFFQ69G5FAV",
			},
			rules: map[string]interface{}{
				"id": "ulid",
			},
			valid: true,
		},
		{
			name: "ulid rule fails with invalid ULID",
			data: map[string]interface{}{
				"id": "not-a-ulid",
			},
			rules: map[string]interface{}{
				"id": "ulid",
			},
			valid: false,
		},
		{
			name: "hex_color rule passes with valid hex color",
			data: map[string]interface{}{
				"color": "#FF5733",
			},
			rules: map[string]interface{}{
				"color": "hex_color",
			},
			valid: true,
		},
		{
			name: "hex_color rule passes with short hex color",
			data: map[string]interface{}{
				"color": "#F57",
			},
			rules: map[string]interface{}{
				"color": "hex_color",
			},
			valid: true,
		},
		{
			name: "hex_color rule fails with invalid hex color",
			data: map[string]interface{}{
				"color": "red",
			},
			rules: map[string]interface{}{
				"color": "hex_color",
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

func TestComplexValidation(t *testing.T) {
	factory := NewFactory()
	
	// Test multiple rules on single field
	t.Run("multiple rules on single field", func(t *testing.T) {
		data := map[string]interface{}{
			"email": "user@example.com",
		}
		rules := map[string]interface{}{
			"email": "required|string|email",
		}
		
		validator := factory.Make(data, rules)
		if !validator.Passes() {
			t.Errorf("Expected validation to pass, but it failed. Errors: %v", validator.Errors().All())
		}
	})
	
	// Test validation with custom messages
	t.Run("custom error messages", func(t *testing.T) {
		data := map[string]interface{}{
			"name": "",
		}
		rules := map[string]interface{}{
			"name": "required",
		}
		messages := map[string]string{
			"name.required": "Please enter your name",
		}
		
		validator := factory.Make(data, rules, messages)
		if validator.Passes() {
			t.Errorf("Expected validation to fail")
		}
		
		errorMessage := validator.Errors().First("name")
		if errorMessage != "Please enter your name" {
			t.Errorf("Expected custom error message, got: %s", errorMessage)
		}
	})
	
	// Test validation stops on first failure
	t.Run("stop on first failure", func(t *testing.T) {
		data := map[string]interface{}{
			"field1": "",
			"field2": "",
		}
		rules := map[string]interface{}{
			"field1": "required",
			"field2": "required",
		}
		
		validator := factory.Make(data, rules).StopOnFirstFailure()
		if validator.Passes() {
			t.Errorf("Expected validation to fail")
		}
		
		// Should only have one error due to stop on first failure
		errorCount := validator.Errors().Count()
		if errorCount != 1 {
			t.Errorf("Expected 1 error due to stop on first failure, got: %d", errorCount)
		}
	})
	
	// Test conditional validation with Sometimes
	t.Run("conditional validation", func(t *testing.T) {
		data := map[string]interface{}{
			"type":     "email",
			"contact":  "user@example.com",
		}
		
		validator := factory.Make(data, map[string]interface{}{})
		validator.Sometimes("contact", []Rule{&EmailRule{}}, func(data map[string]interface{}) bool {
			return data["type"] == "email"
		})
		
		if !validator.Passes() {
			t.Errorf("Expected conditional validation to pass, but it failed. Errors: %v", validator.Errors().All())
		}
	})
}

func TestErrorBag(t *testing.T) {
	errorBag := NewErrorBag()
	
	// Test adding errors
	errorBag.Add("name", "Name is required")
	errorBag.Add("email", "Email is invalid")
	errorBag.Add("email", "Email is already taken")
	
	// Test Has method
	if !errorBag.Has("name") {
		t.Error("Expected error bag to have 'name' error")
	}
	
	if !errorBag.Has("email") {
		t.Error("Expected error bag to have 'email' error")
	}
	
	if errorBag.Has("phone") {
		t.Error("Expected error bag not to have 'phone' error")
	}
	
	// Test Get method
	nameErrors := errorBag.Get("name")
	if len(nameErrors) != 1 || nameErrors[0] != "Name is required" {
		t.Errorf("Expected 1 name error, got: %v", nameErrors)
	}
	
	emailErrors := errorBag.Get("email")
	if len(emailErrors) != 2 {
		t.Errorf("Expected 2 email errors, got: %v", emailErrors)
	}
	
	// Test First method
	firstNameError := errorBag.First("name")
	if firstNameError != "Name is required" {
		t.Errorf("Expected first name error to be 'Name is required', got: %s", firstNameError)
	}
	
	// Test Count method
	totalCount := errorBag.Count()
	if totalCount != 3 {
		t.Errorf("Expected total count to be 3, got: %d", totalCount)
	}
	
	// Test IsEmpty/IsNotEmpty
	if errorBag.IsEmpty() {
		t.Error("Expected error bag not to be empty")
	}
	
	if !errorBag.IsNotEmpty() {
		t.Error("Expected error bag to be not empty")
	}
	
	// Test All method
	allErrors := errorBag.All()
	if len(allErrors) != 2 { // 2 fields with errors
		t.Errorf("Expected 2 fields with errors, got: %d", len(allErrors))
	}
}

func TestValidatorMethods(t *testing.T) {
	factory := NewFactory()
	
	data := map[string]interface{}{
		"name":  "John",
		"email": "invalid-email",
		"age":   25,
	}
	rules := map[string]interface{}{
		"name":  "required|string",
		"email": "required|email",
		"age":   "integer|min:18",
	}
	
	validator := factory.Make(data, rules)
	
	// Test validation fails
	if validator.Passes() {
		t.Error("Expected validation to fail due to invalid email")
	}
	
	if !validator.Fails() {
		t.Error("Expected Fails() to return true")
	}
	
	// Test Valid() returns only valid data
	validData := validator.Valid()
	if len(validData) != 2 { // name and age should be valid
		t.Errorf("Expected 2 valid fields, got: %d", len(validData))
	}
	
	if validData["name"] != "John" {
		t.Error("Expected valid data to contain name")
	}
	
	if validData["age"] != 25 {
		t.Error("Expected valid data to contain age")
	}
	
	// Test Invalid() returns only invalid data
	invalidData := validator.Invalid()
	if len(invalidData) != 1 { // only email should be invalid
		t.Errorf("Expected 1 invalid field, got: %d", len(invalidData))
	}
	
	if invalidData["email"] != "invalid-email" {
		t.Error("Expected invalid data to contain email")
	}
}

func TestFactoryMethods(t *testing.T) {
	factory := NewFactory()
	
	// Test setting custom messages at factory level
	factory.SetCustomMessages(map[string]string{
		"required": "The :attribute field is mandatory",
	})
	
	// Test setting custom attributes at factory level
	factory.SetCustomAttributes(map[string]string{
		"email": "Email Address",
	})
	
	data := map[string]interface{}{
		"email": "",
	}
	rules := map[string]interface{}{
		"email": "required",
	}
	
	validator := factory.Make(data, rules)
	
	if validator.Passes() {
		t.Error("Expected validation to fail")
	}
	
	// Check if custom attribute name is used in error message
	errorMessage := validator.Errors().First("email")
	if !strings.Contains(errorMessage, "Email Address") {
		t.Errorf("Expected error message to contain custom attribute name, got: %s", errorMessage)
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test IsNil function
	if !IsNil(nil) {
		t.Error("Expected IsNil(nil) to be true")
	}
	
	var ptr *string
	if !IsNil(ptr) {
		t.Error("Expected IsNil(nil pointer) to be true")
	}
	
	if IsNil("string") {
		t.Error("Expected IsNil(string) to be false")
	}
	
	// Test ToString function
	if ToString(123) != "123" {
		t.Error("Expected ToString(123) to be '123'")
	}
	
	if ToString("hello") != "hello" {
		t.Error("Expected ToString('hello') to be 'hello'")
	}
	
	// Test ToFloat64 function
	if val, ok := ToFloat64(123); !ok || val != 123.0 {
		t.Error("Expected ToFloat64(123) to return (123.0, true)")
	}
	
	if _, ok := ToFloat64("string"); ok {
		t.Error("Expected ToFloat64('string') to return false")
	}
	
	// Test GetSize function
	if size, ok := GetSize("hello"); !ok || size != 5 {
		t.Error("Expected GetSize('hello') to return (5, true)")
	}
	
	if size, ok := GetSize([]int{1, 2, 3}); !ok || size != 3 {
		t.Error("Expected GetSize([1,2,3]) to return (3, true)")
	}
	
	if size, ok := GetSize(42); !ok || size != 42 {
		t.Error("Expected GetSize(42) to return (42, true)")
	}
}