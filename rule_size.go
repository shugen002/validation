package validation

import "fmt"

func getSize(ctx *ValidationContext) float64 {
	if ctx.HasNumericRule && (ctx.memory["numeric"] == true || isNumeric(ctx.FieldValue)) {
		ctx.memory["numeric"] = true
		var num float64
		fmt.Sscanf(ctx.FieldValue, "%f", &num)
		return num
	} else {
		return float64(len(ctx.FieldValue))
	}
}

func constructSizeRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("size rule requires a size argument")
	}
	var expectedSize float64
	_, err := fmt.Sscanf(args[0], "%f", &expectedSize)
	if err != nil {
		return nil, fmt.Errorf("invalid size argument: %s", args[0])
	}

	return func(ctx *ValidationContext) (bool, error) {
		actualSize := getSize(ctx)
		if actualSize != expectedSize {
			return false, fmt.Errorf("%s must be %v in size", ctx.FieldName, expectedSize)
		}
		return true, nil
	}, nil
}

func constructMinRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("min rule requires a minimum argument")
	}
	var minSize float64
	_, err := fmt.Sscanf(args[0], "%f", &minSize)
	if err != nil {
		return nil, fmt.Errorf("invalid minimum argument: %s", args[0])
	}

	return func(ctx *ValidationContext) (bool, error) {
		actualSize := getSize(ctx)
		if actualSize < minSize {
			return false, fmt.Errorf("%s must be at least %v in size", ctx.FieldName, minSize)
		}
		return true, nil
	}, nil
}

func constructMaxRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("max rule requires a maximum argument")
	}
	var maxSize float64
	_, err := fmt.Sscanf(args[0], "%f", &maxSize)
	if err != nil {
		return nil, fmt.Errorf("invalid maximum argument: %s", args[0])
	}

	return func(ctx *ValidationContext) (bool, error) {
		actualSize := getSize(ctx)
		if actualSize > maxSize {
			return false, fmt.Errorf("%s must be at most %v in size", ctx.FieldName, maxSize)
		}
		return true, nil
	}, nil
}

func constructBetweenRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("between rule requires two arguments")
	}
	var minSize, maxSize float64
	_, err := fmt.Sscanf(args[0], "%f", &minSize)
	if err != nil {
		return nil, fmt.Errorf("invalid minimum argument: %s", args[0])
	}
	_, err = fmt.Sscanf(args[1], "%f", &maxSize)
	if err != nil {
		return nil, fmt.Errorf("invalid maximum argument: %s", args[1])
	}
	if minSize > maxSize {
		return nil, fmt.Errorf("minimum size cannot be greater than maximum size")
	}

	return func(ctx *ValidationContext) (bool, error) {
		actualSize := getSize(ctx)
		if actualSize < minSize || actualSize > maxSize {
			return false, fmt.Errorf("%s must be between %v and %v in size", ctx.FieldName, minSize, maxSize)
		}
		return true, nil
	}, nil
}

func constructGtRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("gt rule requires a minimum argument")
	}

	return func(ctx *ValidationContext) (bool, error) {
		targetField, hasTargetField := ctx.Raw[args[0]]
		fieldSize := getSize(ctx)
		if !hasTargetField && (isNumeric(ctx.FieldValue) && isNumeric(args[0])) {
			argValue := 0.0
			_, err := fmt.Sscanf(args[0], "%f", &argValue)
			if err != nil {
				return false, fmt.Errorf("invalid gt argument: %s", args[0])
			}
			if fieldSize > argValue {
				return true, nil
			} else {
				return false, fmt.Errorf("%s must be greater than %v", ctx.FieldName, argValue)
			}
		}

		if isNumeric(args[0]) {
			return false, fmt.Errorf("gt argument must be a field name or a numeric value")
		}

		if ctx.HasNumericRule && isNumeric(ctx.FieldValue) && isNumeric(targetField) {
			argValue := 0.0
			_, err := fmt.Sscanf(targetField, "%f", &argValue)
			if err != nil {
				return false, fmt.Errorf("invalid gt argument: %s", args[0])
			}
			if fieldSize > argValue {
				return true, nil
			} else {
				return false, fmt.Errorf("%s must be greater than %v", ctx.FieldName, args[0])
			}
		}
		targetSize, err := ctx.GetValue(args[0])
		if err != nil {
			return false, fmt.Errorf("gt argument must be a field name or a numeric value")
		}
		if fieldSize > targetSize {
			return true, nil
		} else {
			return false, fmt.Errorf("%s must be greater than %s", ctx.FieldName, args[0])
		}
	}, nil
}

func constructGteRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("gte rule requires a minimum argument")
	}

	return func(ctx *ValidationContext) (bool, error) {
		targetField, hasTargetField := ctx.Raw[args[0]]
		fieldSize := getSize(ctx)
		if !hasTargetField && (isNumeric(ctx.FieldValue) && isNumeric(args[0])) {
			argValue := 0.0
			_, err := fmt.Sscanf(args[0], "%f", &argValue)
			if err != nil {
				return false, fmt.Errorf("invalid gte argument: %s", args[0])
			}
			if fieldSize >= argValue {
				return true, nil
			} else {
				return false, fmt.Errorf("%s must be greater than or equal to %v", ctx.FieldName, argValue)
			}
		}

		if isNumeric(args[0]) {
			return false, fmt.Errorf("gte argument must be a field name or a numeric value")
		}
		if ctx.HasNumericRule && isNumeric(ctx.FieldValue) && isNumeric(targetField) {
			argValue := 0.0
			_, err := fmt.Sscanf(targetField, "%f", &argValue)
			if err != nil {
				return false, fmt.Errorf("invalid gte argument: %s", args[0])
			}
			if fieldSize >= argValue {
				return true, nil
			} else {
				return false, fmt.Errorf("%s must be greater than or equal to %v", ctx.FieldName, args[0])
			}
		}
		targetSize, err := ctx.GetValue(args[0])
		if err != nil {
			return false, fmt.Errorf("gte argument must be a field name or a numeric value")
		}
		if fieldSize >= targetSize {
			return true, nil
		} else {
			return false, fmt.Errorf("%s must be greater than or equal to %s", ctx.FieldName, args[0])
		}
	}, nil
}

func constructLtRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("lt rule requires a maximum argument")
	}

	return func(ctx *ValidationContext) (bool, error) {
		targetField, hasTargetField := ctx.Raw[args[0]]
		fieldSize := getSize(ctx)
		if !hasTargetField && (isNumeric(ctx.FieldValue) && isNumeric(args[0])) {
			argValue := 0.0
			_, err := fmt.Sscanf(args[0], "%f", &argValue)
			if err != nil {
				return false, fmt.Errorf("invalid lt argument: %s", args[0])
			}
			if fieldSize < argValue {
				return true, nil
			} else {
				return false, fmt.Errorf("%s must be less than %v", ctx.FieldName, argValue)
			}
		}

		if isNumeric(args[0]) {
			return false, fmt.Errorf("lt argument must be a field name or a numeric value")
		}
		if ctx.HasNumericRule && isNumeric(ctx.FieldValue) && isNumeric(targetField) {
			argValue := 0.0
			_, err := fmt.Sscanf(targetField, "%f", &argValue)
			if err != nil {
				return false, fmt.Errorf("invalid lt argument: %s", args[0])
			}
			if fieldSize < argValue {
				return true, nil
			} else {
				return false, fmt.Errorf("%s must be less than %v", ctx.FieldName, args[0])
			}
		}
		targetSize, err := ctx.GetValue(args[0])
		if err != nil {
			return false, fmt.Errorf("lt argument must be a field name or a numeric value")
		}
		if fieldSize < targetSize {
			return true, nil
		} else {
			return false, fmt.Errorf("%s must be less than %s", ctx.FieldName, args[0])
		}
	}, nil
}

func constructLteRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("lte rule requires a maximum argument")
	}

	return func(ctx *ValidationContext) (bool, error) {
		targetField, hasTargetField := ctx.Raw[args[0]]
		fieldSize := getSize(ctx)
		if !hasTargetField && (isNumeric(ctx.FieldValue) && isNumeric(args[0])) {
			argValue := 0.0
			_, err := fmt.Sscanf(args[0], "%f", &argValue)
			if err != nil {
				return false, fmt.Errorf("invalid lte argument: %s", args[0])
			}
			if fieldSize <= argValue {
				return true, nil
			} else {
				return false, fmt.Errorf("%s must be less than or equal to %v", ctx.FieldName, argValue)
			}
		}

		if isNumeric(args[0]) {
			return false, fmt.Errorf("lte argument must be a field name or a numeric value")
		}
		if ctx.HasNumericRule && isNumeric(ctx.FieldValue) && isNumeric(targetField) {
			argValue := 0.0
			_, err := fmt.Sscanf(targetField, "%f", &argValue)
			if err != nil {
				return false, fmt.Errorf("invalid lte argument: %s", args[0])
			}
			if fieldSize <= argValue {
				return true, nil
			} else {
				return false, fmt.Errorf("%s must be less than or equal to %v", ctx.FieldName, args[0])
			}
		}
		targetSize, err := ctx.GetValue(args[0])
		if err != nil {
			return false, fmt.Errorf("lte argument must be a field name or a numeric value")
		}
		if fieldSize <= targetSize {
			return true, nil
		} else {
			return false, fmt.Errorf("%s must be less than or equal to %s", ctx.FieldName, args[0])
		}
	}, nil
}

var embeddedSizeRules = map[string]RuleConstructor{
	"size":    constructSizeRule,
	"min":     constructMinRule,
	"max":     constructMaxRule,
	"gt":      constructGtRule,
	"gte":     constructGteRule,
	"lt":      constructLtRule,
	"lte":     constructLteRule,
	"between": constructBetweenRule,
}
