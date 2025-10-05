package validation

import (
	"fmt"
	"strings"
)

func Nullable(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		if strings.TrimSpace(ctx.FieldValue) == "" {
			return false, nil
		}
		return true, nil
	}, nil
}

func Required(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		if strings.TrimSpace(ctx.FieldValue) == "" {
			return false, fmt.Errorf("%s is required", ctx.FieldName)
		}
		return true, nil
	}, nil
}

func Missing(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		if ctx.FieldValue != "" {
			return false, fmt.Errorf("%s must be missing", ctx.FieldName)
		}
		return true, nil
	}, nil
}

var embeddedUtilitiesRules = map[string]RuleConstructor{
	"nullable": Nullable,
	"required": Required,
	"missing":  Missing,
}
