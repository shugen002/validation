package validation

import (
	"reflect"
	"strings"
)

// validator is the main validation engine
type validator struct {
	data                   map[string]interface{}
	rules                  map[string][]Rule
	customMessages         map[string]string
	customAttributes       map[string]string
	errorBag               *ErrorBag
	stopOnFirstFailure     bool
	excludeUnvalidatedKeys bool
	failedRules            map[string]map[string]interface{}
}

// NewValidator creates a new validator instance
func NewValidator(data map[string]interface{}, rules map[string][]Rule, customMessages map[string]string) Validator {
	return &validator{
		data:            data,
		rules:           rules,
		customMessages:  customMessages,
		customAttributes: make(map[string]string),
		errorBag:        NewErrorBag(),
		failedRules:     make(map[string]map[string]interface{}),
	}
}

// Validate runs the validation and throws an exception if it fails
func (v *validator) Validate() error {
	if v.Fails() {
		return NewValidationException("The given data was invalid.", v.errorBag)
	}
	return nil
}

// ValidateWithBag runs validation with a custom error bag name
func (v *validator) ValidateWithBag(errorBag string) error {
	// In this implementation, we'll just use the default behavior
	// In a more complete implementation, you might support multiple error bags
	return v.Validate()
}

// Passes returns true if validation passes
func (v *validator) Passes() bool {
	v.runValidation()
	return v.errorBag.IsEmpty()
}

// Fails returns true if validation fails
func (v *validator) Fails() bool {
	return !v.Passes()
}

// Errors returns the error bag
func (v *validator) Errors() *ErrorBag {
	if v.errorBag.IsEmpty() {
		v.runValidation()
	}
	return v.errorBag
}

// Valid returns the valid data
func (v *validator) Valid() map[string]interface{} {
	v.runValidation()
	
	valid := make(map[string]interface{})
	for field := range v.rules {
		if !v.errorBag.Has(field) {
			if value, exists := v.data[field]; exists {
				valid[field] = value
			}
		}
	}
	return valid
}

// Invalid returns the invalid data
func (v *validator) Invalid() map[string]interface{} {
	v.runValidation()
	
	invalid := make(map[string]interface{})
	for field := range v.errorBag.All() {
		if value, exists := v.data[field]; exists {
			invalid[field] = value
		}
	}
	return invalid
}

// AddRule adds a rule to a field
func (v *validator) AddRule(field string, rules ...Rule) Validator {
	if v.rules == nil {
		v.rules = make(map[string][]Rule)
	}
	v.rules[field] = append(v.rules[field], rules...)
	return v
}

// Sometimes adds conditional rules
func (v *validator) Sometimes(field string, rules []Rule, callback func(data map[string]interface{}) bool) Validator {
	if callback(v.data) {
		v.AddRule(field, rules...)
	}
	return v
}

// StopOnFirstFailure configures the validator to stop on first failure
func (v *validator) StopOnFirstFailure() Validator {
	v.stopOnFirstFailure = true
	return v
}

// runValidation executes the validation process
func (v *validator) runValidation() {
	// Clear previous errors
	v.errorBag = NewErrorBag()
	v.failedRules = make(map[string]map[string]interface{})
	
	// Validate each field
	for field, fieldRules := range v.rules {
		value := v.getValue(field)
		
		for _, rule := range fieldRules {
			// Set data for data-aware rules
			if dataAware, ok := rule.(DataAwareRule); ok {
				dataAware.SetData(v.data)
			}
			
			// Set validator for validator-aware rules
			if validatorAware, ok := rule.(ValidatorAwareRule); ok {
				validatorAware.SetValidator(v)
			}
			
			// Check if the field is required or has a value
			if !v.shouldValidateField(field, value, rule) {
				continue
			}
			
			// Run the validation
			if !rule.Passes(field, value) {
				v.addFailure(field, rule)
				
				// Stop on first failure if configured
				if v.stopOnFirstFailure {
					return
				}
				
				// Stop validating this field if it's an implicit rule and failed
				if _, isImplicit := rule.(ImplicitRule); isImplicit {
					break
				}
			}
		}
	}
}

// shouldValidateField determines if a field should be validated
func (v *validator) shouldValidateField(field string, value interface{}, rule Rule) bool {
	// Always validate implicit rules
	if _, isImplicit := rule.(ImplicitRule); isImplicit {
		return true
	}
	
	// Skip validation if the field is missing and not required
	if !v.hasField(field) {
		return false
	}
	
	// Skip validation if the value is nil and nullable
	if IsNil(value) {
		return false
	}
	
	return true
}

// hasField checks if a field exists in the data
func (v *validator) hasField(field string) bool {
	_, exists := v.data[field]
	return exists
}

// getValue gets the value of a field
func (v *validator) getValue(field string) interface{} {
	if strings.Contains(field, ".") {
		return v.getNestedValue(field)
	}
	return v.data[field]
}

// getNestedValue gets a nested value using dot notation
func (v *validator) getNestedValue(field string) interface{} {
	parts := strings.Split(field, ".")
	current := v.data
	
	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part, return the value
			return current[part]
		}
		
		// Navigate deeper
		if next, exists := current[part]; exists {
			if nextMap, ok := next.(map[string]interface{}); ok {
				current = nextMap
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
	
	return nil
}

// addFailure adds a validation failure
func (v *validator) addFailure(field string, rule Rule) {
	// Add to failed rules tracking
	if v.failedRules[field] == nil {
		v.failedRules[field] = make(map[string]interface{})
	}
	
	ruleType := reflect.TypeOf(rule).String()
	v.failedRules[field][ruleType] = nil
	
	// Get error message
	message := v.getErrorMessage(field, rule)
	
	// Add to error bag
	v.errorBag.Add(field, message)
}

// getErrorMessage gets the error message for a failed rule
func (v *validator) getErrorMessage(field string, rule Rule) string {
	// Try custom message first
	ruleType := strings.ToLower(strings.TrimSuffix(reflect.TypeOf(rule).Elem().Name(), "Rule"))
	if customMessage, exists := v.customMessages[field+"."+ruleType]; exists {
		return v.formatMessage(customMessage, field)
	}
	
	// Try just the rule name without field prefix
	if customMessage, exists := v.customMessages[ruleType]; exists {
		return v.formatMessage(customMessage, field)
	}
	
	// Use rule's default message
	message := rule.Message()
	return v.formatMessage(message, field)
}

// formatMessage formats an error message
func (v *validator) formatMessage(message, field string) string {
	// Get attribute name (custom or default)
	attribute := field
	if customAttr, exists := v.customAttributes[field]; exists {
		attribute = customAttr
	}
	
	// Replace placeholders
	message = strings.ReplaceAll(message, ":attribute", attribute)
	message = strings.ReplaceAll(message, ":field", field)
	
	return message
}

// SetCustomAttributes sets custom attribute names
func (v *validator) SetCustomAttributes(attributes map[string]string) {
	if v.customAttributes == nil {
		v.customAttributes = make(map[string]string)
	}
	for key, value := range attributes {
		v.customAttributes[key] = value
	}
}

// SetCustomMessages sets custom error messages
func (v *validator) SetCustomMessages(messages map[string]string) {
	if v.customMessages == nil {
		v.customMessages = make(map[string]string)
	}
	for key, value := range messages {
		v.customMessages[key] = value
	}
}