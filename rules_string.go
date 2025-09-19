package validation

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode"
)

// String validation rules

// EmailRule validates that a field is a valid email address
type EmailRule struct {
	Validations []string
}

func (r *EmailRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if str == "" {
		return false
	}
	
	_, err := mail.ParseAddress(str)
	return err == nil
}

func (r *EmailRule) Message() string {
	return "The :attribute must be a valid email address."
}

// AlphaRule validates that a field contains only alphabetic characters
type AlphaRule struct {
	ASCII bool // if true, only allow ASCII alphabetic characters
}

func (r *AlphaRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if str == "" {
		return false
	}
	
	for _, char := range str {
		if r.ASCII {
			if char > 127 || !unicode.IsLetter(char) {
				return false
			}
		} else {
			if !unicode.IsLetter(char) {
				return false
			}
		}
	}
	return true
}

func (r *AlphaRule) Message() string {
	if r.ASCII {
		return "The :attribute may only contain ASCII letters."
	}
	return "The :attribute may only contain letters."
}

// AlphaNumRule validates that a field contains only alphanumeric characters
type AlphaNumRule struct {
	ASCII bool // if true, only allow ASCII alphanumeric characters
}

func (r *AlphaNumRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if str == "" {
		return false
	}
	
	for _, char := range str {
		if r.ASCII {
			if char > 127 || (!unicode.IsLetter(char) && !unicode.IsDigit(char)) {
				return false
			}
		} else {
			if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
				return false
			}
		}
	}
	return true
}

func (r *AlphaNumRule) Message() string {
	if r.ASCII {
		return "The :attribute may only contain ASCII letters and numbers."
	}
	return "The :attribute may only contain letters and numbers."
}

// AlphaDashRule validates that a field contains only alphanumeric characters, dashes, and underscores
type AlphaDashRule struct {
	ASCII bool // if true, only allow ASCII characters
}

func (r *AlphaDashRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if str == "" {
		return false
	}
	
	for _, char := range str {
		if r.ASCII {
			if char > 127 || (!unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '-' && char != '_') {
				return false
			}
		} else {
			if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '-' && char != '_' {
				return false
			}
		}
	}
	return true
}

func (r *AlphaDashRule) Message() string {
	if r.ASCII {
		return "The :attribute may only contain ASCII letters, numbers, dashes, and underscores."
	}
	return "The :attribute may only contain letters, numbers, dashes, and underscores."
}

// RegexRule validates that a field matches a regular expression
type RegexRule struct {
	Pattern string
}

func (r *RegexRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	regex, err := regexp.Compile(r.Pattern)
	if err != nil {
		return false
	}
	return regex.MatchString(str)
}

func (r *RegexRule) Message() string {
	return "The :attribute format is invalid."
}

// NotRegexRule validates that a field does not match a regular expression
type NotRegexRule struct {
	Pattern string
}

func (r *NotRegexRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	regex, err := regexp.Compile(r.Pattern)
	if err != nil {
		return true // If pattern is invalid, consider it as passed
	}
	return !regex.MatchString(str)
}

func (r *NotRegexRule) Message() string {
	return "The :attribute format is invalid."
}

// StartsWithRule validates that a field starts with one of the given values
type StartsWithRule struct {
	Prefixes []string
}

func (r *StartsWithRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	for _, prefix := range r.Prefixes {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}

func (r *StartsWithRule) Message() string {
	return "The :attribute must start with one of the following: " + strings.Join(r.Prefixes, ", ") + "."
}

// EndsWithRule validates that a field ends with one of the given values
type EndsWithRule struct {
	Suffixes []string
}

func (r *EndsWithRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	for _, suffix := range r.Suffixes {
		if strings.HasSuffix(str, suffix) {
			return true
		}
	}
	return false
}

func (r *EndsWithRule) Message() string {
	return "The :attribute must end with one of the following: " + strings.Join(r.Suffixes, ", ") + "."
}

// DoesntStartWithRule validates that a field does not start with any of the given values
type DoesntStartWithRule struct {
	Prefixes []string
}

func (r *DoesntStartWithRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	for _, prefix := range r.Prefixes {
		if strings.HasPrefix(str, prefix) {
			return false
		}
	}
	return true
}

func (r *DoesntStartWithRule) Message() string {
	return "The :attribute may not start with one of the following: " + strings.Join(r.Prefixes, ", ") + "."
}

// DoesntEndWithRule validates that a field does not end with any of the given values
type DoesntEndWithRule struct {
	Suffixes []string
}

func (r *DoesntEndWithRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	for _, suffix := range r.Suffixes {
		if strings.HasSuffix(str, suffix) {
			return false
		}
	}
	return true
}

func (r *DoesntEndWithRule) Message() string {
	return "The :attribute may not end with one of the following: " + strings.Join(r.Suffixes, ", ") + "."
}

// UppercaseRule validates that a field is all uppercase
type UppercaseRule struct{}

func (r *UppercaseRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	return strings.ToUpper(str) == str
}

func (r *UppercaseRule) Message() string {
	return "The :attribute must be uppercase."
}

// LowercaseRule validates that a field is all lowercase
type LowercaseRule struct{}

func (r *LowercaseRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	return strings.ToLower(str) == str
}

func (r *LowercaseRule) Message() string {
	return "The :attribute must be lowercase."
}

// AsciiRule validates that a field contains only ASCII characters
type AsciiRule struct{}

func (r *AsciiRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	for _, char := range str {
		if char > 127 {
			return false
		}
	}
	return true
}

func (r *AsciiRule) Message() string {
	return "The :attribute must only contain ASCII characters."
}