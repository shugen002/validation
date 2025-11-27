package validation

import (
	"fmt"
	"regexp"
)

const numericRegex = `^-?\d+(\.\d+)?$`

var numbeicRegexp *regexp.Regexp

func isNumeric(str string) bool {
	if numbeicRegexp == nil {
		numbeicRegexp = regexp.MustCompile(numericRegex)
	}
	return numbeicRegexp.MatchString(str)
}

func constructNumericRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		if !isNumeric(ctx.FieldValue) {
			return false, newErrValidationFailed(fmt.Sprintf("%s must be a numeric value", ctx.FieldName))
		}
		ctx.memory["numeric"] = true
		return true, nil
	}, nil
}

const integerRegex = `^-?\d+$`

var integerRegexp *regexp.Regexp

func isInteger(str string) bool {
	if integerRegexp == nil {
		integerRegexp = regexp.MustCompile(integerRegex)
	}
	return integerRegexp.MatchString(str)
}

func constructIntergerRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		if !isInteger(ctx.FieldValue) {
			return false, newErrValidationFailed(fmt.Sprintf("%s must be an integer value", ctx.FieldName))
		}
		ctx.memory["integer"] = true
		ctx.memory["numeric"] = true
		return true, nil
	}, nil
}

// decimal:min,max
// decimal:value
// The field under validation must be numeric and must contain the specified number of decimal places:
func constructDecimalRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	min := 0
	max := 0
	if len(args) >= 1 {
		_, err := fmt.Sscanf(args[0], "%d", &min)
		max = min
		if err != nil {
			return nil, &ErrParsingRules{Reason: fmt.Sprintf("invalid minimum decimal places: %v", err)}
		}
	}
	if len(args) == 2 {
		_, err := fmt.Sscanf(args[1], "%d", &max)
		if err != nil {
			return nil, &ErrParsingRules{Reason: fmt.Sprintf("invalid maximum decimal places: %v", err)}
		}
	}
	if min < 0 || max < 0 || min > max {
		return nil, &ErrParsingRules{Reason: "invalid decimal places range"}
	}
	return func(ctx *ValidationContext) (bool, error) {
		if _, ok := ctx.memory["numeric"]; !ok && !isNumeric(ctx.FieldValue) {
			return false, newErrValidationFailed(fmt.Sprintf("%s must be a numeric value", ctx.FieldName))
		}
		ctx.memory["numeric"] = true
		parts := regexp.MustCompile(`\.`).Split(ctx.FieldValue, -1)
		decimalPlaces := 0
		if len(parts) == 2 {
			decimalPlaces = len(parts[1])
		}
		if decimalPlaces < min {
			return false, newErrValidationFailed(fmt.Sprintf("%s must have at least %d decimal places", ctx.FieldName, min))
		}
		if decimalPlaces > max {
			return false, newErrValidationFailed(fmt.Sprintf("%s must have at most %d decimal places", ctx.FieldName, max))
		}
		return true, nil
	}, nil
}

func constructDigitsRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	value := 0
	if len(args) == 1 {
		_, err := fmt.Sscanf(args[0], "%d", &value)
		if err != nil {
			return nil, &ErrParsingRules{Reason: fmt.Sprintf("invalid digit length: %v", err)}
		}
		if value < 0 {
			return nil, &ErrParsingRules{Reason: "digit length must be non-negative"}
		}
	}
	return func(ctx *ValidationContext) (bool, error) {
		if _, ok := ctx.memory["integer"]; !ok && !isInteger(ctx.FieldValue) {
			return false, newErrValidationFailed(fmt.Sprintf("%s must be an integer value", ctx.FieldName))
		}
		ctx.memory["integer"] = true
		ctx.memory["numeric"] = true
		if value > 0 && len(ctx.FieldValue) != value {
			return false, newErrValidationFailed(fmt.Sprintf("%s must be exactly %d digits long", ctx.FieldName, value))
		}
		return true, nil
	}, nil
}

func constructMinDigitsRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	min := 0
	if len(args) != 1 {
		return nil, &ErrParsingRules{Reason: "MinDigits rule requires one parameter"}
	}
	_, err := fmt.Sscanf(args[0], "%d", &min)
	if err != nil {
		return nil, &ErrParsingRules{Reason: fmt.Sprintf("invalid minimum digit length: %v", err)}
	}
	if min < 0 {
		return nil, &ErrParsingRules{Reason: "minimum digit length must be non-negative"}
	}
	return func(ctx *ValidationContext) (bool, error) {
		if _, ok := ctx.memory["integer"]; !ok && !isInteger(ctx.FieldValue) {
			return false, newErrValidationFailed(fmt.Sprintf("%s must be an integer value", ctx.FieldName))
		}
		ctx.memory["integer"] = true
		ctx.memory["numeric"] = true
		if len(ctx.FieldValue) < min {
			return false, newErrValidationFailed(fmt.Sprintf("%s must be at least %d digits long", ctx.FieldName, min))
		}
		return true, nil
	}, nil
}

func constructMaxDigitsRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	max := 0
	if len(args) < 1 {
		return nil, &ErrParsingRules{Reason: "MaxDigits rule requires one parameter"}
	}
	_, err := fmt.Sscanf(args[0], "%d", &max)
	if err != nil {
		return nil, &ErrParsingRules{Reason: fmt.Sprintf("invalid maximum digit length: %v", err)}
	}
	if max < 0 {
		return nil, &ErrParsingRules{Reason: "maximum digit length must be non-negative"}
	}
	return func(ctx *ValidationContext) (bool, error) {
		if _, ok := ctx.memory["integer"]; !ok && !isInteger(ctx.FieldValue) {
			return false, newErrValidationFailed(fmt.Sprintf("%s must be an integer value", ctx.FieldName))
		}
		ctx.memory["integer"] = true
		ctx.memory["numeric"] = true
		if len(ctx.FieldValue) > max {
			return false, newErrValidationFailed(fmt.Sprintf("%s must be at most %d digits long", ctx.FieldName, max))
		}
		return true, nil
	}, nil
}

func constructDigitsBetweenRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	min := 0
	max := 0
	if len(args) != 2 {
		return nil, &ErrParsingRules{Reason: "DigitsBetween rule requires two parameters"}
	}
	_, err := fmt.Sscanf(args[0], "%d", &min)
	if err != nil {
		return nil, &ErrParsingRules{Reason: fmt.Sprintf("invalid minimum digit length: %v", err)}
	}
	_, err = fmt.Sscanf(args[1], "%d", &max)
	if err != nil {
		return nil, &ErrParsingRules{Reason: fmt.Sprintf("invalid maximum digit length: %v", err)}
	}
	if min < 0 || max < 0 || min > max {
		return nil, &ErrParsingRules{Reason: "invalid digit length range"}
	}
	return func(ctx *ValidationContext) (bool, error) {
		if _, ok := ctx.memory["integer"]; !ok && !isInteger(ctx.FieldValue) {
			return false, newErrValidationFailed(fmt.Sprintf("%s must be an integer value", ctx.FieldName))
		}
		ctx.memory["integer"] = true
		ctx.memory["numeric"] = true
		length := len(ctx.FieldValue)
		if length < min || length > max {
			return false, newErrValidationFailed(fmt.Sprintf("%s must be between %d and %d digits long", ctx.FieldName, min, max))
		}
		return true, nil
	}, nil
}

var embeddedNumberRules = map[string]RuleConstructor{
	"numeric":        constructNumericRule,
	"integer":        constructIntergerRule,
	"int":            constructIntergerRule,
	"decimal":        constructDecimalRule,
	"digits":         constructDigitsRule,
	"digits_between": constructDigitsBetweenRule,
	"min_digits":     constructMinDigitsRule,
	"max_digits":     constructMaxDigitsRule,
}
