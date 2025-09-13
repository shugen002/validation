package validation

import (
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Field relationship validation rules

// SameRule validates that a field has the same value as another field
type SameRule struct {
	Field string
	data  map[string]interface{}
}

func (r *SameRule) Passes(attribute string, value interface{}) bool {
	otherValue, exists := r.data[r.Field]
	if !exists {
		return false
	}
	return value == otherValue
}

func (r *SameRule) Message() string {
	return "The :attribute and " + r.Field + " must match."
}

func (r *SameRule) SetData(data map[string]interface{}) {
	r.data = data
}

// DifferentRule validates that a field has a different value from other fields
type DifferentRule struct {
	Fields []string
	data   map[string]interface{}
}

func (r *DifferentRule) Passes(attribute string, value interface{}) bool {
	for _, field := range r.Fields {
		if otherValue, exists := r.data[field]; exists {
			if value == otherValue {
				return false
			}
		}
	}
	return true
}

func (r *DifferentRule) Message() string {
	return "The :attribute and " + strings.Join(r.Fields, ", ") + " must be different."
}

func (r *DifferentRule) SetData(data map[string]interface{}) {
	r.data = data
}

// ConfirmedRule validates that a field has a matching confirmation field
type ConfirmedRule struct {
	Field string
	data  map[string]interface{}
}

func (r *ConfirmedRule) Passes(attribute string, value interface{}) bool {
	confirmationField := r.Field
	if confirmationField == "" {
		confirmationField = attribute + "_confirmation"
	}
	
	otherValue, exists := r.data[confirmationField]
	if !exists {
		return false
	}
	return value == otherValue
}

func (r *ConfirmedRule) Message() string {
	return "The :attribute confirmation does not match."
}

func (r *ConfirmedRule) SetData(data map[string]interface{}) {
	r.data = data
}

// Network validation rules

// UrlRule validates that a field is a valid URL
type UrlRule struct {
	Protocols []string
}

func (r *UrlRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if str == "" {
		return false
	}
	
	parsedURL, err := url.Parse(str)
	if err != nil {
		return false
	}
	
	// Check if URL has a scheme
	if parsedURL.Scheme == "" {
		return false
	}
	
	// Check if URL has a host
	if parsedURL.Host == "" {
		return false
	}
	
	// If specific protocols are specified, check them
	if len(r.Protocols) > 0 {
		for _, protocol := range r.Protocols {
			if parsedURL.Scheme == protocol {
				return true
			}
		}
		return false
	}
	
	// Default: allow http and https
	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}

func (r *UrlRule) Message() string {
	return "The :attribute format is invalid."
}

// IpRule validates that a field is a valid IP address
type IpRule struct{}

func (r *IpRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	return net.ParseIP(str) != nil
}

func (r *IpRule) Message() string {
	return "The :attribute must be a valid IP address."
}

// Ipv4Rule validates that a field is a valid IPv4 address
type Ipv4Rule struct{}

func (r *Ipv4Rule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	ip := net.ParseIP(str)
	return ip != nil && ip.To4() != nil
}

func (r *Ipv4Rule) Message() string {
	return "The :attribute must be a valid IPv4 address."
}

// Ipv6Rule validates that a field is a valid IPv6 address
type Ipv6Rule struct{}

func (r *Ipv6Rule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	ip := net.ParseIP(str)
	return ip != nil && ip.To4() == nil
}

func (r *Ipv6Rule) Message() string {
	return "The :attribute must be a valid IPv6 address."
}

// MacAddressRule validates that a field is a valid MAC address
type MacAddressRule struct{}

func (r *MacAddressRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	_, err := net.ParseMAC(str)
	return err == nil
}

func (r *MacAddressRule) Message() string {
	return "The :attribute must be a valid MAC address."
}

// Date and time validation rules

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
		"01/02/2006",
		"02-01-2006",
		time.RFC3339,
		time.RFC822,
		time.RFC1123,
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

// DateFormatRule validates that a field matches specific date formats
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
	return "The :attribute does not match the format " + strings.Join(r.Formats, " or ") + "."
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

// AfterRule validates that a date is after another date
type AfterRule struct {
	Date string
	data map[string]interface{}
}

func (r *AfterRule) Passes(attribute string, value interface{}) bool {
	valueTime := r.parseDate(ToString(value))
	if valueTime.IsZero() {
		return false
	}
	
	var compareTime time.Time
	
	// Check if Date is a field name
	if fieldValue, exists := r.data[r.Date]; exists {
		compareTime = r.parseDate(ToString(fieldValue))
	} else {
		compareTime = r.parseDate(r.Date)
	}
	
	if compareTime.IsZero() {
		return false
	}
	
	return valueTime.After(compareTime)
}

func (r *AfterRule) Message() string {
	return "The :attribute must be a date after " + r.Date + "."
}

func (r *AfterRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *AfterRule) parseDate(dateStr string) time.Time {
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}
	
	return time.Time{}
}

// BeforeRule validates that a date is before another date
type BeforeRule struct {
	Date string
	data map[string]interface{}
}

func (r *BeforeRule) Passes(attribute string, value interface{}) bool {
	valueTime := r.parseDate(ToString(value))
	if valueTime.IsZero() {
		return false
	}
	
	var compareTime time.Time
	
	// Check if Date is a field name
	if fieldValue, exists := r.data[r.Date]; exists {
		compareTime = r.parseDate(ToString(fieldValue))
	} else {
		compareTime = r.parseDate(r.Date)
	}
	
	if compareTime.IsZero() {
		return false
	}
	
	return valueTime.Before(compareTime)
}

func (r *BeforeRule) Message() string {
	return "The :attribute must be a date before " + r.Date + "."
}

func (r *BeforeRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *BeforeRule) parseDate(dateStr string) time.Time {
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}
	
	return time.Time{}
}

// AfterOrEqualRule validates that a date is after or equal to another date
type AfterOrEqualRule struct {
	Date string
	data map[string]interface{}
}

func (r *AfterOrEqualRule) Passes(attribute string, value interface{}) bool {
	valueTime := r.parseDate(ToString(value))
	if valueTime.IsZero() {
		return false
	}
	
	var compareTime time.Time
	
	// Check if Date is a field name
	if fieldValue, exists := r.data[r.Date]; exists {
		compareTime = r.parseDate(ToString(fieldValue))
	} else {
		compareTime = r.parseDate(r.Date)
	}
	
	if compareTime.IsZero() {
		return false
	}
	
	return valueTime.After(compareTime) || valueTime.Equal(compareTime)
}

func (r *AfterOrEqualRule) Message() string {
	return "The :attribute must be a date after or equal to " + r.Date + "."
}

func (r *AfterOrEqualRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *AfterOrEqualRule) parseDate(dateStr string) time.Time {
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}
	
	return time.Time{}
}

// BeforeOrEqualRule validates that a date is before or equal to another date
type BeforeOrEqualRule struct {
	Date string
	data map[string]interface{}
}

func (r *BeforeOrEqualRule) Passes(attribute string, value interface{}) bool {
	valueTime := r.parseDate(ToString(value))
	if valueTime.IsZero() {
		return false
	}
	
	var compareTime time.Time
	
	// Check if Date is a field name
	if fieldValue, exists := r.data[r.Date]; exists {
		compareTime = r.parseDate(ToString(fieldValue))
	} else {
		compareTime = r.parseDate(r.Date)
	}
	
	if compareTime.IsZero() {
		return false
	}
	
	return valueTime.Before(compareTime) || valueTime.Equal(compareTime)
}

func (r *BeforeOrEqualRule) Message() string {
	return "The :attribute must be a date before or equal to " + r.Date + "."
}

func (r *BeforeOrEqualRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *BeforeOrEqualRule) parseDate(dateStr string) time.Time {
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}
	
	return time.Time{}
}

// DateEqualsRule validates that a date equals another date
type DateEqualsRule struct {
	Date string
	data map[string]interface{}
}

func (r *DateEqualsRule) Passes(attribute string, value interface{}) bool {
	valueTime := r.parseDate(ToString(value))
	if valueTime.IsZero() {
		return false
	}
	
	var compareTime time.Time
	
	// Check if Date is a field name
	if fieldValue, exists := r.data[r.Date]; exists {
		compareTime = r.parseDate(ToString(fieldValue))
	} else {
		compareTime = r.parseDate(r.Date)
	}
	
	if compareTime.IsZero() {
		return false
	}
	
	return valueTime.Equal(compareTime)
}

func (r *DateEqualsRule) Message() string {
	return "The :attribute must be a date equal to " + r.Date + "."
}

func (r *DateEqualsRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *DateEqualsRule) parseDate(dateStr string) time.Time {
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}
	
	return time.Time{}
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
	
	_, err := time.LoadLocation(str)
	return err == nil
}

func (r *TimezoneRule) Message() string {
	return "The :attribute must be a valid timezone."
}

// Other validation rules

// UuidRule validates that a field is a valid UUID
type UuidRule struct {
	Version interface{} // can be int (1-8) or "max"
}

func (r *UuidRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if str == "" {
		return false
	}
	
	parsedUUID, err := uuid.Parse(str)
	if err != nil {
		return false
	}
	
	// Check version if specified
	if r.Version != nil {
		if version, ok := r.Version.(int); ok {
			if version >= 1 && version <= 8 {
				return int(parsedUUID.Version()) == version
			}
		}
		// "max" means any version is acceptable
	}
	
	return true
}

func (r *UuidRule) Message() string {
	return "The :attribute must be a valid UUID."
}

// UlidRule validates that a field is a valid ULID
type UlidRule struct{}

func (r *UlidRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if str == "" {
		return false
	}
	
	// Basic ULID validation - 26 characters, Crockford Base32
	if len(str) != 26 {
		return false
	}
	
	// ULID uses Crockford Base32 alphabet
	validChars := "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	for _, char := range strings.ToUpper(str) {
		if !strings.ContainsRune(validChars, char) {
			return false
		}
	}
	
	return true
}

func (r *UlidRule) Message() string {
	return "The :attribute must be a valid ULID."
}

// HexColorRule validates that a field is a valid hex color
type HexColorRule struct{}

func (r *HexColorRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if str == "" {
		return false
	}
	
	// Hex color patterns: #RGB, #RRGGBB, #ARGB, #AARRGGBB
	pattern := `^#([0-9A-Fa-f]{3}|[0-9A-Fa-f]{6}|[0-9A-Fa-f]{4}|[0-9A-Fa-f]{8})$`
	matched, err := regexp.MatchString(pattern, str)
	return err == nil && matched
}

func (r *HexColorRule) Message() string {
	return "The :attribute must be a valid hex color."
}