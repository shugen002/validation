package validation

import "fmt"

var defaultNumericRules = []string{"numeric", "integer", "int", "decimal"}

type ValidationContext struct {
	FieldName      string
	FieldValue     string
	Type           string
	HasNumericRule bool
	Raw            map[string]string
	memory         map[string]interface{}
	Rules          []string
	GetValue       func(field string) (float64, error)
	GetStr         func(field string) (string, error)
}

type ValidationRule func(ctx *ValidationContext) (next bool, err error)

type RuleConstructor func(cfg map[string]interface{}, args ...string) (ValidationRule, error)

type ParseResult struct {
	Rules          []ValidationRule
	RuleNames      []string
	HasNumericRule bool
}

type Validator struct {
	rules map[string]ParseResult
}

func (v *Validator) Validate(value map[string]string) error {
	for field, rules := range v.rules {
		ctx := &ValidationContext{
			FieldName:      field,
			FieldValue:     value[field],
			Type:           "string", // assuming string type for simplicity
			Raw:            value,
			memory:         make(map[string]interface{}),
			HasNumericRule: rules.HasNumericRule,
			Rules:          rules.RuleNames,
			GetStr: func(f string) (string, error) {
				if val, exists := value[f]; exists {
					return val, nil
				}
				return "", fmt.Errorf("field %s not found", f)
			},
			GetValue: func(f string) (float64, error) {
				valueStr, exists := value[f]
				if !exists {
					return 0, fmt.Errorf("field %s not found", f)
				}
				valueHasNumericRule := false
				if r, ok := v.rules[f]; ok {
					valueHasNumericRule = r.HasNumericRule
				}
				if valueHasNumericRule && isNumeric(valueStr) {
					var num float64
					fmt.Sscanf(valueStr, "%f", &num)
					return num, nil
				}
				return float64(len(valueStr)), nil
			},
		}
		for i := 0; i < len(rules.Rules); i++ {
			rule := rules.Rules[i]
			next, err := rule(ctx)
			if err != nil {
				return err
			}
			if !next {
				break
			}
		}
	}
	return nil
}

func (ctx *ValidationContext) SetMemory(key string, value interface{}) {
	ctx.memory[key] = value
}

func (ctx *ValidationContext) GetMemory(key string, defaultValue interface{}) interface{} {
	if val, exists := ctx.memory[key]; exists {
		return val
	}
	return defaultValue
}
