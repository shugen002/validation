package validation

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Numeric validation rules

// MinRule validates minimum value/length
type MinRule struct {
	Min float64
}

func (r *MinRule) Passes(attribute string, value interface{}) bool {
	size, ok := GetSize(value)
	if !ok {
		return false
	}
	return size >= r.Min
}

func (r *MinRule) Message() string {
	return fmt.Sprintf("The :attribute must be at least %.0f.", r.Min)
}

// MaxRule validates maximum value/length
type MaxRule struct {
	Max float64
}

func (r *MaxRule) Passes(attribute string, value interface{}) bool {
	size, ok := GetSize(value)
	if !ok {
		return false
	}
	return size <= r.Max
}

func (r *MaxRule) Message() string {
	return fmt.Sprintf("The :attribute may not be greater than %.0f.", r.Max)
}

// BetweenRule validates that a value/length is between min and max
type BetweenRule struct {
	Min float64
	Max float64
}

func (r *BetweenRule) Passes(attribute string, value interface{}) bool {
	size, ok := GetSize(value)
	if !ok {
		return false
	}
	return size >= r.Min && size <= r.Max
}

func (r *BetweenRule) Message() string {
	return fmt.Sprintf("The :attribute must be between %.0f and %.0f.", r.Min, r.Max)
}

// SizeRule validates exact size/length/value
type SizeRule struct {
	Size float64
}

func (r *SizeRule) Passes(attribute string, value interface{}) bool {
	size, ok := GetSize(value)
	if !ok {
		return false
	}
	return size == r.Size
}

func (r *SizeRule) Message() string {
	return fmt.Sprintf("The :attribute must be %.0f.", r.Size)
}

// DigitsRule validates that a field is numeric and has an exact number of digits
type DigitsRule struct {
	Length int
}

func (r *DigitsRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	
	// Remove any minus sign for negative numbers
	if strings.HasPrefix(str, "-") {
		str = str[1:]
	}
	
	// Check if all characters are digits
	for _, char := range str {
		if char < '0' || char > '9' {
			return false
		}
	}
	
	return len(str) == r.Length
}

func (r *DigitsRule) Message() string {
	return fmt.Sprintf("The :attribute must be %d digits.", r.Length)
}

// DigitsBetweenRule validates that a field is numeric and has digit count between min and max
type DigitsBetweenRule struct {
	Min int
	Max int
}

func (r *DigitsBetweenRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	
	// Remove any minus sign for negative numbers
	if strings.HasPrefix(str, "-") {
		str = str[1:]
	}
	
	// Check if all characters are digits
	for _, char := range str {
		if char < '0' || char > '9' {
			return false
		}
	}
	
	length := len(str)
	return length >= r.Min && length <= r.Max
}

func (r *DigitsBetweenRule) Message() string {
	return fmt.Sprintf("The :attribute must be between %d and %d digits.", r.Min, r.Max)
}

// MinDigitsRule validates minimum number of digits
type MinDigitsRule struct {
	Min int
}

func (r *MinDigitsRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	
	// Remove any minus sign for negative numbers
	if strings.HasPrefix(str, "-") {
		str = str[1:]
	}
	
	// Count digits only
	digitCount := 0
	for _, char := range str {
		if char >= '0' && char <= '9' {
			digitCount++
		}
	}
	
	return digitCount >= r.Min
}

func (r *MinDigitsRule) Message() string {
	return fmt.Sprintf("The :attribute must have at least %d digits.", r.Min)
}

// DecimalRule validates that a field is numeric and has specified decimal places
type DecimalRule struct {
	Min int
	Max int
}

func (r *DecimalRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	
	// Check if it's a valid number
	if _, err := strconv.ParseFloat(str, 64); err != nil {
		return false
	}
	
	// Find decimal point
	parts := strings.Split(str, ".")
	if len(parts) == 1 {
		// No decimal point - this is an integer
		// For range validation, 0 decimal places should be valid if Min <= 0
		return r.Min <= 0 && (r.Max == 0 || r.Max >= 0)
	}
	
	if len(parts) != 2 {
		return false // Multiple decimal points
	}
	
	decimalPlaces := len(parts[1])
	
	if r.Max == 0 {
		// Exact number of decimal places
		return decimalPlaces == r.Min
	}
	
	// Range of decimal places
	return decimalPlaces >= r.Min && decimalPlaces <= r.Max
}

func (r *DecimalRule) Message() string {
	if r.Max == 0 {
		return fmt.Sprintf("The :attribute must have exactly %d decimal places.", r.Min)
	}
	return fmt.Sprintf("The :attribute must have between %d and %d decimal places.", r.Min, r.Max)
}

// MultipleOfRule validates that a field is a multiple of a given value
type MultipleOfRule struct {
	Value float64
}

func (r *MultipleOfRule) Passes(attribute string, value interface{}) bool {
	val, err := strconv.ParseFloat(ToString(value), 64)
	if err != nil {
		return false
	}
	
	if r.Value == 0 {
		return false
	}
	
	remainder := math.Mod(val, r.Value)
	// Use a small epsilon for floating point comparison
	return math.Abs(remainder) < 1e-10 || math.Abs(remainder-r.Value) < 1e-10
}

func (r *MultipleOfRule) Message() string {
	return fmt.Sprintf("The :attribute must be a multiple of %.2f.", r.Value)
}