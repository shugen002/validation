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
			return false, fmt.Errorf("%s must be a numeric value", ctx.FieldName)
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
			return false, fmt.Errorf("%s must be an integer value", ctx.FieldName)
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
			return nil, fmt.Errorf("invalid minimum decimal places: %v", err)
		}
	}
	if len(args) == 2 {
		_, err := fmt.Sscanf(args[1], "%d", &max)
		if err != nil {
			return nil, fmt.Errorf("invalid maximum decimal places: %v", err)
		}
	}
	if min < 0 || max < 0 || min > max {
		return nil, fmt.Errorf("invalid decimal places range")
	}
	return func(ctx *ValidationContext) (bool, error) {
		if _, ok := ctx.memory["numeric"]; !ok && !isNumeric(ctx.FieldValue) {
			return false, fmt.Errorf("%s must be a numeric value", ctx.FieldName)
		}
		ctx.memory["numeric"] = true
		parts := regexp.MustCompile(`\.`).Split(ctx.FieldValue, -1)
		decimalPlaces := 0
		if len(parts) == 2 {
			decimalPlaces = len(parts[1])
		}
		if decimalPlaces < min {
			return false, fmt.Errorf("%s must have at least %d decimal places", ctx.FieldName, min)
		}
		if decimalPlaces > max {
			return false, fmt.Errorf("%s must have at most %d decimal places", ctx.FieldName, max)
		}
		return true, nil
	}, nil
}

func constructDigitsRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	value := 0
	if len(args) == 1 {
		_, err := fmt.Sscanf(args[0], "%d", &value)
		if err != nil {
			return nil, fmt.Errorf("invalid digit length: %v", err)
		}
		if value < 0 {
			return nil, fmt.Errorf("digit length must be non-negative")
		}
	}
	return func(ctx *ValidationContext) (bool, error) {
		if _, ok := ctx.memory["integer"]; !ok && !isInteger(ctx.FieldValue) {
			return false, fmt.Errorf("%s must be an integer value", ctx.FieldName)
		}
		ctx.memory["integer"] = true
		ctx.memory["numeric"] = true
		if value > 0 && len(ctx.FieldValue) != value {
			return false, fmt.Errorf("%s must be exactly %d digits long", ctx.FieldName, value)
		}
		return true, nil
	}, nil
}

func constructMinDigitsRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	min := 0
	if len(args) != 1 {
		return nil, fmt.Errorf("MinDigits rule requires one parameter")
	}
	_, err := fmt.Sscanf(args[0], "%d", &min)
	if err != nil {
		return nil, fmt.Errorf("invalid minimum digit length: %v", err)
	}
	if min < 0 {
		return nil, fmt.Errorf("minimum digit length must be non-negative")
	}
	return func(ctx *ValidationContext) (bool, error) {
		if _, ok := ctx.memory["integer"]; !ok && !isInteger(ctx.FieldValue) {
			return false, fmt.Errorf("%s must be an integer value", ctx.FieldName)
		}
		ctx.memory["integer"] = true
		ctx.memory["numeric"] = true
		if len(ctx.FieldValue) < min {
			return false, fmt.Errorf("%s must be at least %d digits long", ctx.FieldName, min)
		}
		return true, nil
	}, nil
}

func constructMaxDigitsRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	max := 0
	if len(args) < 1 {
		return nil, fmt.Errorf("MaxDigits rule requires one parameter")
	}
	_, err := fmt.Sscanf(args[0], "%d", &max)
	if err != nil {
		return nil, fmt.Errorf("invalid maximum digit length: %v", err)
	}
	if max < 0 {
		return nil, fmt.Errorf("maximum digit length must be non-negative")
	}
	return func(ctx *ValidationContext) (bool, error) {
		if _, ok := ctx.memory["integer"]; !ok && !isInteger(ctx.FieldValue) {
			return false, fmt.Errorf("%s must be an integer value", ctx.FieldName)
		}
		ctx.memory["integer"] = true
		ctx.memory["numeric"] = true
		if len(ctx.FieldValue) > max {
			return false, fmt.Errorf("%s must be at most %d digits long", ctx.FieldName, max)
		}
		return true, nil
	}, nil
}

func constructDigitsBetweenRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	min := 0
	max := 0
	if len(args) != 2 {
		return nil, fmt.Errorf("DigitsBetween rule requires two parameters")
	}
	_, err := fmt.Sscanf(args[0], "%d", &min)
	if err != nil {
		return nil, fmt.Errorf("invalid minimum digit length: %v", err)
	}
	_, err = fmt.Sscanf(args[1], "%d", &max)
	if err != nil {
		return nil, fmt.Errorf("invalid maximum digit length: %v", err)
	}
	if min < 0 || max < 0 || min > max {
		return nil, fmt.Errorf("invalid digit length range")
	}
	return func(ctx *ValidationContext) (bool, error) {
		if _, ok := ctx.memory["integer"]; !ok && !isInteger(ctx.FieldValue) {
			return false, fmt.Errorf("%s must be an integer value", ctx.FieldName)
		}
		ctx.memory["integer"] = true
		ctx.memory["numeric"] = true
		length := len(ctx.FieldValue)
		if length < min || length > max {
			return false, fmt.Errorf("%s must be between %d and %d digits long", ctx.FieldName, min, max)
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
