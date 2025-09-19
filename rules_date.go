package validation

import (
	"strings"
	"time"
)

// Date validation rules

// DateRule validates that a field is a valid date
type DateRule struct{}

func (r *DateRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if str == "" {
		return false
	}
	
	// Try common date formats
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006/01/02",
		"2006/01/02 15:04:05",
		"01/02/2006",
		"01/02/2006 15:04:05",
		"02-01-2006",
		"02-01-2006 15:04:05",
		time.RFC3339,
		time.RFC822,
		time.RFC850,
		time.Kitchen,
	}
	
	for _, format := range formats {
		if _, err := time.Parse(format, str); err == nil {
			return true
		}
	}
	
	return false
}

func (r *DateRule) Message() string {
	return "The :attribute is not a valid date."
}

// DateFormatRule validates that a field matches one of the specified date formats
type DateFormatRule struct {
	Formats []string
}

func (r *DateFormatRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if str == "" {
		return false
	}
	
	for _, format := range r.Formats {
		// Convert Laravel/PHP format to Go format
		goFormat := r.convertPHPDateFormat(format)
		if _, err := time.Parse(goFormat, str); err == nil {
			return true
		}
	}
	
	return false
}

func (r *DateFormatRule) Message() string {
	return "The :attribute does not match the required format."
}

// convertPHPDateFormat converts PHP date format to Go date format
func (r *DateFormatRule) convertPHPDateFormat(phpFormat string) string {
	// Basic conversion map from PHP to Go date format
	replacements := map[string]string{
		"Y": "2006",
		"y": "06",
		"m": "01",
		"n": "1",
		"d": "02",
		"j": "2",
		"H": "15",
		"h": "03",
		"i": "04",
		"s": "05",
		"A": "PM",
		"a": "pm",
	}
	
	result := phpFormat
	for php, go_ := range replacements {
		result = strings.ReplaceAll(result, php, go_)
	}
	
	return result
}

// AfterRule validates that a field is after a specified date or field value
type AfterRule struct {
	Date string
	data map[string]interface{}
}

func (r *AfterRule) Passes(attribute string, value interface{}) bool {
	valueTime := r.parseTime(ToString(value))
	if valueTime.IsZero() {
		return false
	}
	
	var compareTime time.Time
	
	// Check if Date is a field reference
	if otherValue, exists := r.data[r.Date]; exists {
		compareTime = r.parseTime(ToString(otherValue))
	} else {
		compareTime = r.parseTime(r.Date)
	}
	
	if compareTime.IsZero() {
		return false
	}
	
	return valueTime.After(compareTime)
}

func (r *AfterRule) parseTime(str string) time.Time {
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006/01/02",
		"2006/01/02 15:04:05",
		time.RFC3339,
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, str); err == nil {
			return t
		}
	}
	
	return time.Time{}
}

func (r *AfterRule) Message() string {
	return "The :attribute must be a date after " + r.Date + "."
}

func (r *AfterRule) SetData(data map[string]interface{}) {
	r.data = data
}

// BeforeRule validates that a field is before a specified date or field value
type BeforeRule struct {
	Date string
	data map[string]interface{}
}

func (r *BeforeRule) Passes(attribute string, value interface{}) bool {
	valueTime := r.parseTime(ToString(value))
	if valueTime.IsZero() {
		return false
	}
	
	var compareTime time.Time
	
	// Check if Date is a field reference
	if otherValue, exists := r.data[r.Date]; exists {
		compareTime = r.parseTime(ToString(otherValue))
	} else {
		compareTime = r.parseTime(r.Date)
	}
	
	if compareTime.IsZero() {
		return false
	}
	
	return valueTime.Before(compareTime)
}

func (r *BeforeRule) parseTime(str string) time.Time {
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006/01/02",
		"2006/01/02 15:04:05",
		time.RFC3339,
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, str); err == nil {
			return t
		}
	}
	
	return time.Time{}
}

func (r *BeforeRule) Message() string {
	return "The :attribute must be a date before " + r.Date + "."
}

func (r *BeforeRule) SetData(data map[string]interface{}) {
	r.data = data
}

// AfterOrEqualRule validates that a field is after or equal to a specified date or field value
type AfterOrEqualRule struct {
	Date string
	data map[string]interface{}
}

func (r *AfterOrEqualRule) Passes(attribute string, value interface{}) bool {
	valueTime := r.parseTime(ToString(value))
	if valueTime.IsZero() {
		return false
	}
	
	var compareTime time.Time
	
	// Check if Date is a field reference
	if otherValue, exists := r.data[r.Date]; exists {
		compareTime = r.parseTime(ToString(otherValue))
	} else {
		compareTime = r.parseTime(r.Date)
	}
	
	if compareTime.IsZero() {
		return false
	}
	
	return valueTime.After(compareTime) || valueTime.Equal(compareTime)
}

func (r *AfterOrEqualRule) parseTime(str string) time.Time {
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006/01/02",
		"2006/01/02 15:04:05",
		time.RFC3339,
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, str); err == nil {
			return t
		}
	}
	
	return time.Time{}
}

func (r *AfterOrEqualRule) Message() string {
	return "The :attribute must be a date after or equal to " + r.Date + "."
}

func (r *AfterOrEqualRule) SetData(data map[string]interface{}) {
	r.data = data
}

// BeforeOrEqualRule validates that a field is before or equal to a specified date or field value
type BeforeOrEqualRule struct {
	Date string
	data map[string]interface{}
}

func (r *BeforeOrEqualRule) Passes(attribute string, value interface{}) bool {
	valueTime := r.parseTime(ToString(value))
	if valueTime.IsZero() {
		return false
	}
	
	var compareTime time.Time
	
	// Check if Date is a field reference
	if otherValue, exists := r.data[r.Date]; exists {
		compareTime = r.parseTime(ToString(otherValue))
	} else {
		compareTime = r.parseTime(r.Date)
	}
	
	if compareTime.IsZero() {
		return false
	}
	
	return valueTime.Before(compareTime) || valueTime.Equal(compareTime)
}

func (r *BeforeOrEqualRule) parseTime(str string) time.Time {
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006/01/02",
		"2006/01/02 15:04:05",
		time.RFC3339,
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, str); err == nil {
			return t
		}
	}
	
	return time.Time{}
}

func (r *BeforeOrEqualRule) Message() string {
	return "The :attribute must be a date before or equal to " + r.Date + "."
}

func (r *BeforeOrEqualRule) SetData(data map[string]interface{}) {
	r.data = data
}

// DateEqualsRule validates that a field equals a specified date or field value
type DateEqualsRule struct {
	Date string
	data map[string]interface{}
}

func (r *DateEqualsRule) Passes(attribute string, value interface{}) bool {
	valueTime := r.parseTime(ToString(value))
	if valueTime.IsZero() {
		return false
	}
	
	var compareTime time.Time
	
	// Check if Date is a field reference
	if otherValue, exists := r.data[r.Date]; exists {
		compareTime = r.parseTime(ToString(otherValue))
	} else {
		compareTime = r.parseTime(r.Date)
	}
	
	if compareTime.IsZero() {
		return false
	}
	
	return valueTime.Equal(compareTime)
}

func (r *DateEqualsRule) parseTime(str string) time.Time {
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006/01/02",
		"2006/01/02 15:04:05",
		time.RFC3339,
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, str); err == nil {
			return t
		}
	}
	
	return time.Time{}
}

func (r *DateEqualsRule) Message() string {
	return "The :attribute must be a date equal to " + r.Date + "."
}

func (r *DateEqualsRule) SetData(data map[string]interface{}) {
	r.data = data
}

// TimezoneRule validates that a field is a valid timezone
type TimezoneRule struct {
	Group   string
	Country string
}

func (r *TimezoneRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if str == "" {
		return false
	}
	
	// Try to load the timezone
	_, err := time.LoadLocation(str)
	return err == nil
}

func (r *TimezoneRule) Message() string {
	return "The :attribute must be a valid timezone."
}