package validation

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

// alpha
// The field under validation must be entirely Unicode alphabetic characters contained in [\p{L}] and [\p{M}].
func constructAlphaRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	ascii := false
	if len(args) > 0 && args[0] == "ascii" {
		ascii = true
	}
	return func(ctx *ValidationContext) (bool, error) {
		for _, r := range ctx.FieldValue {
			if ascii {
				if r > 127 || !unicode.IsLetter(r) {
					return false, fmt.Errorf("the %s field must be entirely alphabetic characters (ASCII)", ctx.FieldName)
				}
			} else {
				if !unicode.IsLetter(r) && !unicode.IsMark(r) {
					return false, fmt.Errorf("the %s field must be entirely alphabetic characters", ctx.FieldName)
				}
			}
		}
		return true, nil
	}, nil
}

// alpha_dash
// The field under validation must be entirely Unicode alpha-numeric characters contained in [\p{L}], [\p{M}], [\p{N}], as well as ASCII dashes (-) and ASCII underscores (_).
func constructAlphaDashRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	ascii := false
	if len(args) > 0 && args[0] == "ascii" {
		ascii = true
	}
	return func(ctx *ValidationContext) (bool, error) {
		for _, r := range ctx.FieldValue {
			if ascii {
				if r > 127 || (!unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '-' && r != '_') {
					return false, fmt.Errorf("the %s field must be entirely alpha-numeric characters, dashes, or underscores (ASCII)", ctx.FieldName)
				}
			} else {
				if !unicode.IsLetter(r) && !unicode.IsMark(r) && !unicode.IsNumber(r) && r != '-' && r != '_' {
					return false, fmt.Errorf("the %s field must be entirely alpha-numeric characters, dashes, or underscores", ctx.FieldName)
				}
			}
		}
		return true, nil
	}, nil
}

// alpha_num
// The field under validation must be entirely Unicode alpha-numeric characters contained in [\p{L}], [\p{M}], and [\p{N}].
func constructAlphaNumRule(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	ascii := false
	if len(args) > 0 && args[0] == "ascii" {
		ascii = true
	}
	return func(ctx *ValidationContext) (bool, error) {
		for _, r := range ctx.FieldValue {
			if ascii {
				if r > 127 || (!unicode.IsLetter(r) && !unicode.IsNumber(r)) {
					return false, fmt.Errorf("the %s field must be entirely alpha-numeric characters (ASCII)", ctx.FieldName)
				}
			} else {
				if !unicode.IsLetter(r) && !unicode.IsMark(r) && !unicode.IsNumber(r) {
					return false, fmt.Errorf("the %s field must be entirely alpha-numeric characters", ctx.FieldName)
				}
			}
		}
		return true, nil
	}, nil
}

// ascii
// The field under validation must be entirely 7-bit ASCII characters.
func constructAsciiRule(_cfg map[string]interface{}, _args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		for _, r := range ctx.FieldValue {
			if r > 127 {
				return false, fmt.Errorf("the %s field must be entirely ASCII characters", ctx.FieldName)
			}
		}
		return true, nil
	}, nil
}

// confirmed
// The field under validation must have a matching field of {field}_confirmation. For example, if the field under validation is password, a matching password_confirmation field must be present in the input.
func constructConfirmed(_cfg map[string]interface{}, _args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		confirmationField := ctx.FieldName + "_confirmation"
		if confirmationValue, ok := ctx.Raw[confirmationField]; !ok || confirmationValue != ctx.FieldValue {
			return false, fmt.Errorf("the %s field must be confirmed", ctx.FieldName)
		}
		return true, nil
	}, nil
}

// different:field
// The field under validation must have a different value than field.
func constructDifferent(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("different rule requires at least 1 argument")
	}
	otherField := args[0]
	return func(ctx *ValidationContext) (bool, error) {
		if otherValue, ok := ctx.Raw[otherField]; ok && otherValue == ctx.FieldValue {
			return false, fmt.Errorf("the %s field must be different from %s", ctx.FieldName, otherField)
		}
		return true, nil
	}, nil
}

// email
// The field under validation must be formatted as an email address.
func constructEmail(_cfg map[string]interface{}, _args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		// Simple email regex for demonstration
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(ctx.FieldValue) {
			return false, fmt.Errorf("the %s field must be a valid email address", ctx.FieldName)
		}
		// Note: For full Laravel compatibility, would need more complex validation based on mode
		return true, nil
	}, nil
}

// ends_with:foo,bar,...
// The field under validation must end with one of the given values.
func constructEndsWith(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("ends_with rule requires at least 1 argument")
	}
	return func(ctx *ValidationContext) (bool, error) {
		for _, suffix := range args {
			if strings.HasSuffix(ctx.FieldValue, suffix) {
				return true, nil
			}
		}
		return false, fmt.Errorf("the %s field must end with one of the following: %s", ctx.FieldName, strings.Join(args, ", "))
	}, nil
}

// hex_color
// The field under validation must contain a valid color value in hexadecimal format.
func constructHexColor(_cfg map[string]interface{}, _args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		hexRegex := regexp.MustCompile(`^#([a-fA-F0-9]{3}|[a-fA-F0-9]{6}|[a-fA-F0-9]{8})$`)
		if !hexRegex.MatchString(ctx.FieldValue) {
			return false, fmt.Errorf("the %s field must be a valid hexadecimal color", ctx.FieldName)
		}
		return true, nil
	}, nil
}

// in:foo,bar,...
// The field under validation must be included in the given list of values.
func constructIn(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("in rule requires at least 1 argument")
	}
	return func(ctx *ValidationContext) (bool, error) {
		for _, allowed := range args {
			if ctx.FieldValue == allowed {
				return true, nil
			}
		}
		return false, fmt.Errorf("the %s field must be one of the following: %s", ctx.FieldName, strings.Join(args, ", "))
	}, nil
}

// json
// The field under validation must be a valid constructJSON string.
func constructJSON(_cfg map[string]interface{}, _args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		var js interface{}
		if err := json.Unmarshal([]byte(ctx.FieldValue), &js); err != nil {
			return false, fmt.Errorf("the %s field must be a valid JSON string", ctx.FieldName)
		}
		return true, nil
	}, nil
}

// ip
// The field under validation must be an constructIP address.
func constructIP(_cfg map[string]interface{}, _args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		if net.ParseIP(ctx.FieldValue) == nil {
			return false, fmt.Errorf("the %s field must be a valid IP address", ctx.FieldName)
		}
		return true, nil
	}, nil
}

// ipv4
// The field under validation must be an constructIPv4 address.
func constructIPv4(_cfg map[string]interface{}, _args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		if ip := net.ParseIP(ctx.FieldValue); ip == nil || ip.To4() == nil {
			return false, fmt.Errorf("the %s field must be a valid IPv4 address", ctx.FieldName)
		}
		return true, nil
	}, nil
}

// ipv6
// The field under validation must be an constructIPv6 address.
func constructIPv6(_cfg map[string]interface{}, _args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		if ip := net.ParseIP(ctx.FieldValue); ip == nil || ip.To4() != nil {
			return false, fmt.Errorf("the %s field must be a valid IPv6 address", ctx.FieldName)
		}
		return true, nil
	}, nil
}

// mac_address
// The field under validation must be a MAC address.
func constructMACAddress(_cfg map[string]interface{}, _args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		if _, err := net.ParseMAC(ctx.FieldValue); err != nil {
			return false, fmt.Errorf("the %s field must be a valid MAC address", ctx.FieldName)
		}
		return true, nil
	}, nil
}

// ulid
// The field under validation must be a valid Universally Unique Lexicographically Sortable Identifier (constructULID).
func constructULID(_cfg map[string]interface{}, _args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		ulidRegex := regexp.MustCompile(`^[0-7][0-9A-HJKMNP-TV-Z]{25}$`)
		if !ulidRegex.MatchString(ctx.FieldValue) {
			return false, fmt.Errorf("the %s field must be a valid ULID", ctx.FieldName)
		}
		return true, nil
	}, nil
}

// lowercase
// The field under validation must be lowercase.
func constructLowercase(_cfg map[string]interface{}, _args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		if ctx.FieldValue != strings.ToLower(ctx.FieldValue) {
			return false, fmt.Errorf("the %s field must be lowercase", ctx.FieldName)
		}
		return true, nil
	}, nil
}

// not_in:foo,bar,...
// The field under validation must not be included in the given list of values.
func constructNotIn(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("not_in rule requires at least 1 argument")
	}
	return func(ctx *ValidationContext) (bool, error) {
		for _, disallowed := range args {
			if ctx.FieldValue == disallowed {
				return false, fmt.Errorf("the %s field must not be one of the following: %s", ctx.FieldName, strings.Join(args, ", "))
			}
		}
		return true, nil
	}, nil
}

// regex:pattern
// The field under validation must match the given regular expression.
func constructRegex(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("regex rule requires 1 argument")
	}
	pattern := strings.Join(args, ",")

	if len(pattern) >= 2 && pattern[0] == '/' {
		flags := ""
		// Find the last '/' to get the pattern without delimiters
		lastSlash := strings.LastIndex(pattern[1:], "/") + 1
		if lastSlash > 0 {
			pattern = pattern[1:lastSlash]
		}
		if len(pattern) > lastSlash+1 {
			flags = pattern[lastSlash+1:]
		}
		if strings.Contains(flags, "i") {
			pattern = "(?i)" + pattern
		}
		if strings.Contains(flags, "m") {
			pattern = "(?m)" + pattern
		}
		if strings.Contains(flags, "s") {
			pattern = "(?s)" + pattern
		}
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %s", pattern)
	}
	return func(ctx *ValidationContext) (bool, error) {
		if !re.MatchString(ctx.FieldValue) {
			return false, fmt.Errorf("the %s field format is invalid", ctx.FieldName)
		}
		return true, nil
	}, nil
}

// not_regex:pattern
// The field under validation must not match the given regular expression.
func constructNotRegex(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("not_regex rule requires 1 argument")
	}
	pattern := args[0]
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %s", pattern)
	}
	return func(ctx *ValidationContext) (bool, error) {
		if re.MatchString(ctx.FieldValue) {
			return false, fmt.Errorf("the %s field format is invalid", ctx.FieldName)
		}
		return true, nil
	}, nil
}

// same:field
// The given field must match the field under validation.
func constructSame(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("same rule requires 1 argument")
	}
	otherField := args[0]
	return func(ctx *ValidationContext) (bool, error) {
		if otherValue, ok := ctx.Raw[otherField]; !ok || otherValue != ctx.FieldValue {
			return false, fmt.Errorf("the %s field must be the same as %s", ctx.FieldName, otherField)
		}
		return true, nil
	}, nil
}

// starts_with:foo,bar,...
// The field under validation must start with one of the given values.
func constructStartsWith(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("starts_with rule requires at least 1 argument")
	}
	return func(ctx *ValidationContext) (bool, error) {
		for _, prefix := range args {
			if strings.HasPrefix(ctx.FieldValue, prefix) {
				return true, nil
			}
		}
		return false, fmt.Errorf("the %s field must start with one of the following: %s", ctx.FieldName, strings.Join(args, ", "))
	}, nil
}

// string
// The field under validation must be a string.
func constructString(_cfg map[string]interface{}, _args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		// In Go, all fields are strings, so always pass
		return true, nil
	}, nil
}

// uppercase
// The field under validation must be uppercase.
func constructUppercase(_cfg map[string]interface{}, _args ...string) (ValidationRule, error) {
	return func(ctx *ValidationContext) (bool, error) {
		if ctx.FieldValue != strings.ToUpper(ctx.FieldValue) {
			return false, fmt.Errorf("the %s field must be uppercase", ctx.FieldName)
		}
		return true, nil
	}, nil
}

// url
// The field under validation must be a valid constructURL.
func constructURL(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	schemes := []string{"http", "https"}
	if len(args) > 0 {
		schemes = args
	}
	return func(ctx *ValidationContext) (bool, error) {
		u, err := url.Parse(ctx.FieldValue)
		if err != nil {
			return false, fmt.Errorf("the %s field must be a valid URL", ctx.FieldName)
		}
		validScheme := false
		for _, scheme := range schemes {
			if u.Scheme == scheme {
				validScheme = true
				break
			}
		}
		if !validScheme {
			return false, fmt.Errorf("the %s field must be a valid URL", ctx.FieldName)
		}
		return true, nil
	}, nil
}

// doesnt_start_with:foo,bar,...
// The field under validation must not start with one of the given values.
func constructDoesntStartWith(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("doesnt_start_with rule requires at least 1 argument")
	}
	return func(ctx *ValidationContext) (bool, error) {
		for _, prefix := range args {
			if strings.HasPrefix(ctx.FieldValue, prefix) {
				return false, fmt.Errorf("the %s field must not start with one of the following: %s", ctx.FieldName, strings.Join(args, ", "))
			}
		}
		return true, nil
	}, nil
}

// doesnt_end_with:foo,bar,...
// The field under validation must not end with one of the given values.
func constructDoesntEndWith(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("doesnt_end_with rule requires at least 1 argument")
	}
	return func(ctx *ValidationContext) (bool, error) {
		for _, suffix := range args {
			if strings.HasSuffix(ctx.FieldValue, suffix) {
				return false, fmt.Errorf("the %s field must not end with one of the following: %s", ctx.FieldName, strings.Join(args, ", "))
			}
		}
		return true, nil
	}, nil
}

// uuid
// The field under validation must be a valid RFC 9562 (version 1, 3, 4, 5, 6, 7, or 8) universally unique identifier (constructUUID).
func constructUUID(_cfg map[string]interface{}, args ...string) (ValidationRule, error) {
	version := ""
	if len(args) > 0 {
		version = args[0]
	}
	return func(ctx *ValidationContext) (bool, error) {
		uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
		if !uuidRegex.MatchString(ctx.FieldValue) {
			return false, fmt.Errorf("the %s field must be a valid UUID", ctx.FieldName)
		}
		if version != "" {
			// Check version
			if len(ctx.FieldValue) >= 14 && ctx.FieldValue[14] != version[0] {
				return false, fmt.Errorf("the %s field must be a valid UUID version %s", ctx.FieldName, version)
			}
		}
		return true, nil
	}, nil
}

var embeddedStringRules = map[string]RuleConstructor{
	"alpha":             constructAlphaRule,
	"alpha_dash":        constructAlphaDashRule,
	"alpha_num":         constructAlphaNumRule,
	"ascii":             constructAsciiRule,
	"confirmed":         constructConfirmed,
	"different":         constructDifferent,
	"doesnt_end_with":   constructDoesntEndWith,
	"doesnt_start_with": constructDoesntStartWith,
	"email":             constructEmail,
	"ends_with":         constructEndsWith,
	"hex_color":         constructHexColor,
	"in":                constructIn,
	"ip":                constructIP,
	"ipv4":              constructIPv4,
	"ipv6":              constructIPv6,
	"json":              constructJSON,
	"lowercase":         constructLowercase,
	"mac_address":       constructMACAddress,
	"not_in":            constructNotIn,
	"regex":             constructRegex,
	"not_regex":         constructNotRegex,
	"same":              constructSame,
	"starts_with":       constructStartsWith,
	"string":            constructString,
	"ulid":              constructULID,
	"uppercase":         constructUppercase,
	"url":               constructURL,
	"uuid":              constructUUID,
}
