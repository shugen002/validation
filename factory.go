package validation

import (
	"maps"
	"slices"
	"strings"
)

type Factory struct {
	rules        map[string]RuleConstructor
	config       map[string]interface{}
	numericRules []string
}

func NewFactory() *Factory {
	embeddedRules := []map[string]RuleConstructor{
		embeddedBooleanRules,
		embeddedStringRules,
		embeddedUtilitiesRules,
		embeddedNumberRules,
		embeddedSizeRules,
		// Add other embedded rule maps here as needed
	}

	// copy embedded rules to avoid external modification
	embeddedRulesCopy := maps.Clone(embeddedRules[0])
	for _, ruleMap := range embeddedRules[1:] {
		maps.Copy(embeddedRulesCopy, ruleMap)
	}
	numericRules := make([]string, len(defaultNumericRules))
	copy(numericRules, defaultNumericRules)
	return &Factory{
		rules:        embeddedRulesCopy,
		config:       make(map[string]interface{}),
		numericRules: numericRules,
	}
}

func (f *Factory) RegisterRule(name string, constructor RuleConstructor) {
	f.rules[name] = constructor
}

func (f *Factory) SetConfig(key string, value interface{}) {
	f.config[key] = value
}

func (f *Factory) UnsetConfig(key string) {
	delete(f.config, key)
}

func (f *Factory) Parse(structRules map[string]string) (*Validator, error) {
	parsedRules := make(map[string]ParseResult, len(structRules))
	for field, ruleStr := range structRules {
		ruleStrs := strings.Split(ruleStr, "|")
		rules := make([]ValidationRule, 0, len(ruleStrs))
		ruleNames := make([]string, 0, len(ruleStrs))
		hasNumeric := false
		for _, r := range ruleStrs {
			parts := strings.SplitN(r, ":", 2)
			ruleName := parts[0]
			ruleName = strings.ToLower(strings.TrimSpace(ruleName))
			if ruleName == "" {
				continue
			}
			if !hasNumeric && slices.Contains(f.numericRules, ruleName) {
				hasNumeric = true
			}
			ruleNames = append(ruleNames, ruleName)
			var args []string
			if len(parts) > 1 {
				args = strings.Split(parts[1], ",")
			}
			constructor, exists := f.rules[ruleName]
			if !exists {
				return nil, &ErrUnknownRule{Rule: ruleName}
			}
			rule, err := constructor(f.config, args...)
			if err != nil {
				return nil, err
			}
			rules = append(rules, rule)
		}

		parsedRules[field] = ParseResult{Rules: rules, RuleNames: ruleNames, HasNumericRule: hasNumeric}
	}
	return &Validator{rules: parsedRules}, nil
}
