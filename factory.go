package validation

import (
	"reflect"
	"strconv"
	"strings"
)

// Factory creates validators
type Factory struct {
	customMessages   map[string]string
	customAttributes map[string]string
}

// NewFactory creates a new validator factory
func NewFactory() *Factory {
	return &Factory{
		customMessages:   make(map[string]string),
		customAttributes: make(map[string]string),
	}
}

// Make creates a new validator with the given data and rules
func (f *Factory) Make(data map[string]interface{}, rules map[string]interface{}, messages ...map[string]string) Validator {
	parsedRules := f.parseRules(rules)
	
	customMessages := make(map[string]string)
	
	// Add factory-level custom messages
	for key, value := range f.customMessages {
		customMessages[key] = value
	}
	
	// Add instance-level custom messages
	if len(messages) > 0 {
		for key, value := range messages[0] {
			customMessages[key] = value
		}
	}
	
	validator := NewValidator(data, parsedRules, customMessages).(*validator)
	validator.SetCustomAttributes(f.customAttributes)
	
	return validator
}

// parseRules converts string rules to Rule objects
func (f *Factory) parseRules(rules map[string]interface{}) map[string][]Rule {
	parsed := make(map[string][]Rule)
	
	for field, fieldRules := range rules {
		parsed[field] = f.parseFieldRules(fieldRules)
	}
	
	return parsed
}

// parseFieldRules parses rules for a single field
func (f *Factory) parseFieldRules(rules interface{}) []Rule {
	var result []Rule
	
	switch r := rules.(type) {
	case string:
		// Parse pipe-separated rules like "required|string|min:3"
		ruleStrings := strings.Split(r, "|")
		for _, ruleStr := range ruleStrings {
			if rule := f.parseRule(strings.TrimSpace(ruleStr)); rule != nil {
				result = append(result, rule)
			}
		}
	case []string:
		// Parse array of rule strings
		for _, ruleStr := range r {
			if rule := f.parseRule(strings.TrimSpace(ruleStr)); rule != nil {
				result = append(result, rule)
			}
		}
	case []interface{}:
		// Parse mixed array (strings and Rule objects)
		for _, ruleItem := range r {
			switch item := ruleItem.(type) {
			case string:
				if rule := f.parseRule(strings.TrimSpace(item)); rule != nil {
					result = append(result, rule)
				}
			case Rule:
				result = append(result, item)
			}
		}
	case Rule:
		// Single rule object
		result = append(result, r)
	}
	
	return result
}

// parseRule parses a single rule string into a Rule object
func (f *Factory) parseRule(ruleStr string) Rule {
	if ruleStr == "" {
		return nil
	}
	
	// Split rule name and parameters
	parts := strings.SplitN(ruleStr, ":", 2)
	ruleName := strings.ToLower(parts[0])
	var params []string
	
	if len(parts) > 1 {
		params = f.parseParameters(parts[1])
	}
	
	// Create rule based on name
	switch ruleName {
	case "required":
		return &RequiredRule{}
	case "string":
		return &StringRule{}
	case "integer", "int":
		strict := len(params) > 0 && params[0] == "strict"
		return &IntegerRule{Strict: strict}
	case "numeric":
		strict := len(params) > 0 && params[0] == "strict"
		return &NumericRule{Strict: strict}
	case "boolean", "bool":
		strict := len(params) > 0 && params[0] == "strict"
		return &BooleanRule{Strict: strict}
	case "accepted":
		return &AcceptedRule{}
	case "accepted_if":
		if len(params) >= 2 {
			return &AcceptedIfRule{Field: params[0], Value: params[1]}
		}
	case "declined":
		return &DeclinedRule{}
	case "declined_if":
		if len(params) >= 2 {
			return &DeclinedIfRule{Field: params[0], Value: params[1]}
		}
	case "array":
		return &ArrayRule{AllowedKeys: params}
	case "json":
		return &JsonRule{}
	case "email":
		return &EmailRule{Validations: params}
	case "alpha":
		ascii := len(params) > 0 && params[0] == "ascii"
		return &AlphaRule{ASCII: ascii}
	case "alpha_num":
		ascii := len(params) > 0 && params[0] == "ascii"
		return &AlphaNumRule{ASCII: ascii}
	case "alpha_dash":
		ascii := len(params) > 0 && params[0] == "ascii"
		return &AlphaDashRule{ASCII: ascii}
	case "min":
		if len(params) > 0 {
			if min, err := strconv.ParseFloat(params[0], 64); err == nil {
				return &MinRule{Min: min}
			}
		}
	case "max":
		if len(params) > 0 {
			if max, err := strconv.ParseFloat(params[0], 64); err == nil {
				return &MaxRule{Max: max}
			}
		}
	case "between":
		if len(params) >= 2 {
			if min, err1 := strconv.ParseFloat(params[0], 64); err1 == nil {
				if max, err2 := strconv.ParseFloat(params[1], 64); err2 == nil {
					return &BetweenRule{Min: min, Max: max}
				}
			}
		}
	case "size":
		if len(params) > 0 {
			if size, err := strconv.ParseFloat(params[0], 64); err == nil {
				return &SizeRule{Size: size}
			}
		}
	case "in":
		return &InRule{Values: params}
	case "not_in":
		return &NotInRule{Values: params}
	case "regex":
		if len(parts) > 1 {
			// For regex, take the entire parameter string as the pattern (don't split on commas)
			pattern := strings.TrimSpace(parts[1])
			// Remove surrounding forward slashes if present (common regex delimiter syntax)
			if len(pattern) >= 2 && pattern[0] == '/' && pattern[len(pattern)-1] == '/' {
				pattern = pattern[1 : len(pattern)-1]
			}
			return &RegexRule{Pattern: pattern}
		}
	case "not_regex":
		if len(parts) > 1 {
			// For not_regex, take the entire parameter string as the pattern (don't split on commas)
			pattern := strings.TrimSpace(parts[1])
			// Remove surrounding forward slashes if present (common regex delimiter syntax)
			if len(pattern) >= 2 && pattern[0] == '/' && pattern[len(pattern)-1] == '/' {
				pattern = pattern[1 : len(pattern)-1]
			}
			return &NotRegexRule{Pattern: pattern}
		}
	case "same":
		if len(params) > 0 {
			return &SameRule{Field: params[0]}
		}
	case "different":
		return &DifferentRule{Fields: params}
	case "confirmed":
		field := ""
		if len(params) > 0 {
			field = params[0]
		}
		return &ConfirmedRule{Field: field}
	case "url":
		return &UrlRule{Protocols: params}
	case "ip":
		return &IpRule{}
	case "ipv4":
		return &Ipv4Rule{}
	case "ipv6":
		return &Ipv6Rule{}
	case "mac_address":
		return &MacAddressRule{}
	case "uuid":
		var version interface{}
		if len(params) > 0 {
			if params[0] == "max" {
				version = "max"
			} else if v, err := strconv.Atoi(params[0]); err == nil {
				version = v
			}
		}
		return &UuidRule{Version: version}
	case "ulid":
		return &UlidRule{}
	case "hex_color":
		return &HexColorRule{}
	case "date":
		return &DateRule{}
	case "date_format":
		return &DateFormatRule{Formats: params}
	case "after":
		if len(params) > 0 {
			return &AfterRule{Date: params[0]}
		}
	case "after_or_equal":
		if len(params) > 0 {
			return &AfterOrEqualRule{Date: params[0]}
		}
	case "before":
		if len(params) > 0 {
			return &BeforeRule{Date: params[0]}
		}
	case "before_or_equal":
		if len(params) > 0 {
			return &BeforeOrEqualRule{Date: params[0]}
		}
	case "date_equals":
		if len(params) > 0 {
			return &DateEqualsRule{Date: params[0]}
		}
	case "timezone":
		group := ""
		country := ""
		if len(params) > 0 {
			group = params[0]
		}
		if len(params) > 1 {
			country = params[1]
		}
		return &TimezoneRule{Group: group, Country: country}
	case "starts_with":
		return &StartsWithRule{Prefixes: params}
	case "ends_with":
		return &EndsWithRule{Suffixes: params}
	case "doesnt_start_with":
		return &DoesntStartWithRule{Prefixes: params}
	case "doesnt_end_with":
		return &DoesntEndWithRule{Suffixes: params}
	case "uppercase":
		return &UppercaseRule{}
	case "lowercase":
		return &LowercaseRule{}
	case "ascii":
		return &AsciiRule{}
	case "nullable":
		return &NullableRule{}
	case "sometimes":
		return &SometimesRule{}
	case "bail":
		return &BailRule{}
	case "digits":
		if len(params) > 0 {
			if length, err := strconv.Atoi(params[0]); err == nil {
				return &DigitsRule{Length: length}
			}
		}
	case "digits_between":
		if len(params) >= 2 {
			if min, err1 := strconv.Atoi(params[0]); err1 == nil {
				if max, err2 := strconv.Atoi(params[1]); err2 == nil {
					return &DigitsBetweenRule{Min: min, Max: max}
				}
			}
		}
	case "filled":
		return &FilledRule{}
	case "present":
		return &PresentRule{}
	case "prohibited":
		return &ProhibitedRule{}
	case "required_if":
		if len(params) >= 2 {
			return &RequiredIfRule{Field: params[0], Value: params[1]}
		}
	case "required_unless":
		if len(params) >= 2 {
			return &RequiredUnlessRule{Field: params[0], Value: params[1]}
		}
	case "required_with":
		return &RequiredWithRule{Fields: params}
	case "required_without":
		return &RequiredWithoutRule{Fields: params}
	}
	
	return nil
}

// parseParameters parses rule parameters from a string
func (f *Factory) parseParameters(paramStr string) []string {
	var params []string
	var current strings.Builder
	var inQuotes bool
	var quoteChar rune
	
	for _, char := range paramStr {
		switch {
		case (char == '"' || char == '\'') && !inQuotes:
			inQuotes = true
			quoteChar = char
		case char == quoteChar && inQuotes:
			inQuotes = false
		case char == ',' && !inQuotes:
			if current.Len() > 0 {
				params = append(params, strings.TrimSpace(current.String()))
				current.Reset()
			}
		default:
			current.WriteRune(char)
		}
	}
	
	if current.Len() > 0 {
		params = append(params, strings.TrimSpace(current.String()))
	}
	
	return params
}

// SetCustomMessages sets factory-level custom messages
func (f *Factory) SetCustomMessages(messages map[string]string) {
	for key, value := range messages {
		f.customMessages[key] = value
	}
}

// SetCustomAttributes sets factory-level custom attributes
func (f *Factory) SetCustomAttributes(attributes map[string]string) {
	for key, value := range attributes {
		f.customAttributes[key] = value
	}
}

// Convenience methods for common validation scenarios

// ValidateStruct validates a struct using reflection and tags
func (f *Factory) ValidateStruct(s interface{}) (Validator, error) {
	data, rules, err := f.parseStruct(s)
	if err != nil {
		return nil, err
	}
	
	return f.Make(data, rules), nil
}

// parseStruct extracts data and rules from a struct using reflection
func (f *Factory) parseStruct(s interface{}) (map[string]interface{}, map[string]interface{}, error) {
	value := reflect.ValueOf(s)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	
	if value.Kind() != reflect.Struct {
		return nil, nil, &ValidationException{Message: "input must be a struct"}
	}
	
	typ := value.Type()
	data := make(map[string]interface{})
	rules := make(map[string]interface{})
	
	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := value.Field(i)
		
		// Skip unexported fields
		if !fieldValue.CanInterface() {
			continue
		}
		
		// Get field name (use json tag if available)
		fieldName := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			if idx := strings.Index(jsonTag, ","); idx != -1 {
				fieldName = jsonTag[:idx]
			} else {
				fieldName = jsonTag
			}
		}
		
		// Get validation rules from validate tag
		if validateTag := field.Tag.Get("validate"); validateTag != "" && validateTag != "-" {
			rules[fieldName] = validateTag
		}
		
		// Set field value
		data[fieldName] = fieldValue.Interface()
	}
	
	return data, rules, nil
}