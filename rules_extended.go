package validation

import (
	"fmt"
	"net"
	"net/url"
	"reflect"
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

// Conditional required rules

// RequiredIfRule validates that a field is required when another field equals a specific value
type RequiredIfRule struct {
	Field string
	Value string
	data  map[string]interface{}
}

func (r *RequiredIfRule) Passes(attribute string, value interface{}) bool {
	// Check if the condition field equals the required value
	if fieldValue, exists := r.data[r.Field]; exists {
		if ToString(fieldValue) == r.Value {
			// Field is required, check if it's present and not empty
			if IsNil(value) {
				return false
			}
			
			switch v := value.(type) {
			case string:
				return strings.TrimSpace(v) != ""
			case []interface{}, map[string]interface{}:
				return reflect.ValueOf(v).Len() > 0
			default:
				rv := reflect.ValueOf(value)
				switch rv.Kind() {
				case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
					return rv.Len() > 0
				}
			}
			
			return true
		}
	}
	
	// Field is not required in this case
	return true
}

func (r *RequiredIfRule) Message() string {
	return "The :attribute field is required when " + r.Field + " is " + r.Value + "."
}

func (r *RequiredIfRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *RequiredIfRule) IsImplicit() bool {
	return true
}

// RequiredUnlessRule validates that a field is required unless another field equals a specific value
type RequiredUnlessRule struct {
	Field string
	Value string
	data  map[string]interface{}
}

func (r *RequiredUnlessRule) Passes(attribute string, value interface{}) bool {
	// Check if the condition field equals the exception value
	if fieldValue, exists := r.data[r.Field]; exists {
		if ToString(fieldValue) == r.Value {
			// Field is not required in this case
			return true
		}
	}
	
	// Field is required, check if it's present and not empty
	if IsNil(value) {
		return false
	}
	
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) != ""
	case []interface{}, map[string]interface{}:
		return reflect.ValueOf(v).Len() > 0
	default:
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
			return rv.Len() > 0
		}
	}
	
	return true
}

func (r *RequiredUnlessRule) Message() string {
	return "The :attribute field is required unless " + r.Field + " is " + r.Value + "."
}

func (r *RequiredUnlessRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *RequiredUnlessRule) IsImplicit() bool {
	return true
}

// RequiredWithRule validates that a field is required when any of the other specified fields are present
type RequiredWithRule struct {
	Fields []string
	data   map[string]interface{}
}

func (r *RequiredWithRule) Passes(attribute string, value interface{}) bool {
	// Check if any of the specified fields are present and not empty
	anyFieldPresent := false
	for _, field := range r.Fields {
		if fieldValue, exists := r.data[field]; exists && !IsNil(fieldValue) {
			switch v := fieldValue.(type) {
			case string:
				if strings.TrimSpace(v) != "" {
					anyFieldPresent = true
					break
				}
			case []interface{}, map[string]interface{}:
				if reflect.ValueOf(v).Len() > 0 {
					anyFieldPresent = true
					break
				}
			default:
				rv := reflect.ValueOf(fieldValue)
				switch rv.Kind() {
				case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
					if rv.Len() > 0 {
						anyFieldPresent = true
						break
					}
				default:
					anyFieldPresent = true
					break
				}
			}
		}
	}
	
	if !anyFieldPresent {
		// None of the fields are present, so this field is not required
		return true
	}
	
	// At least one field is present, so this field is required
	if IsNil(value) {
		return false
	}
	
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) != ""
	case []interface{}, map[string]interface{}:
		return reflect.ValueOf(v).Len() > 0
	default:
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
			return rv.Len() > 0
		}
	}
	
	return true
}

func (r *RequiredWithRule) Message() string {
	return "The :attribute field is required when " + strings.Join(r.Fields, ", ") + " are present."
}

func (r *RequiredWithRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *RequiredWithRule) IsImplicit() bool {
	return true
}

// RequiredWithoutRule validates that a field is required when any of the other specified fields are not present
type RequiredWithoutRule struct {
	Fields []string
	data   map[string]interface{}
}

func (r *RequiredWithoutRule) Passes(attribute string, value interface{}) bool {
	// Check if any of the specified fields are missing or empty
	anyFieldMissing := false
	for _, field := range r.Fields {
		if fieldValue, exists := r.data[field]; !exists || IsNil(fieldValue) {
			anyFieldMissing = true
			break
		} else {
			switch v := fieldValue.(type) {
			case string:
				if strings.TrimSpace(v) == "" {
					anyFieldMissing = true
					break
				}
			case []interface{}, map[string]interface{}:
				if reflect.ValueOf(v).Len() == 0 {
					anyFieldMissing = true
					break
				}
			default:
				rv := reflect.ValueOf(fieldValue)
				switch rv.Kind() {
				case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
					if rv.Len() == 0 {
						anyFieldMissing = true
						break
					}
				}
			}
		}
	}
	
	if !anyFieldMissing {
		// All fields are present, so this field is not required
		return true
	}
	
	// At least one field is missing, so this field is required
	if IsNil(value) {
		return false
	}
	
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) != ""
	case []interface{}, map[string]interface{}:
		return reflect.ValueOf(v).Len() > 0
	default:
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
			return rv.Len() > 0
		}
	}
	
	return true
}

func (r *RequiredWithoutRule) Message() string {
	return "The :attribute field is required when " + strings.Join(r.Fields, ", ") + " are not present."
}

func (r *RequiredWithoutRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *RequiredWithoutRule) IsImplicit() bool {
	return true
}

// MissingIfRule validates that a field is not present when another field equals a value
type MissingIfRule struct {
	Field string
	Value string
	data  map[string]interface{}
	validator Validator
}

func (r *MissingIfRule) Passes(attribute string, value interface{}) bool {
	// Check if the condition field equals the required value
	if fieldValue, exists := r.data[r.Field]; exists {
		if ToString(fieldValue) == r.Value {
			// Field should be missing
			if r.validator != nil {
				return !r.validator.HasField(attribute)
			}
			return false // If we're here, field is present
		}
	}
	
	// Condition not met, field can be present or absent
	return true
}

func (r *MissingIfRule) Message() string {
	return "The :attribute field must not be present when " + r.Field + " is " + r.Value + "."
}

func (r *MissingIfRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *MissingIfRule) SetValidator(validator Validator) {
	r.validator = validator
}

func (r *MissingIfRule) IsImplicit() bool {
	return true
}

// MissingUnlessRule validates that a field is not present unless another field equals a value
type MissingUnlessRule struct {
	Field string
	Value string
	data  map[string]interface{}
	validator Validator
}

func (r *MissingUnlessRule) Passes(attribute string, value interface{}) bool {
	// Check if the condition field equals the exception value
	if fieldValue, exists := r.data[r.Field]; exists {
		if ToString(fieldValue) == r.Value {
			// Field can be present
			return true
		}
	}
	
	// Field should be missing
	if r.validator != nil {
		return !r.validator.HasField(attribute)
	}
	return false // If we're here, field is present
}

func (r *MissingUnlessRule) Message() string {
	return "The :attribute field must not be present unless " + r.Field + " is " + r.Value + "."
}

func (r *MissingUnlessRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *MissingUnlessRule) SetValidator(validator Validator) {
	r.validator = validator
}

func (r *MissingUnlessRule) IsImplicit() bool {
	return true
}

// MissingWithRule validates that a field is not present when other fields are present
type MissingWithRule struct {
	Fields []string
	data   map[string]interface{}
	validator Validator
}

func (r *MissingWithRule) Passes(attribute string, value interface{}) bool {
	// Check if any of the specified fields are present
	anyFieldPresent := false
	for _, field := range r.Fields {
		if _, exists := r.data[field]; exists {
			anyFieldPresent = true
			break
		}
	}
	
	if !anyFieldPresent {
		// No dependency fields present, this field can be present or absent
		return true
	}
	
	// At least one dependency field is present, this field should be missing
	if r.validator != nil {
		return !r.validator.HasField(attribute)
	}
	return false // If we're here, field is present
}

func (r *MissingWithRule) Message() string {
	return "The :attribute field must not be present when " + strings.Join(r.Fields, ", ") + " are present."
}

func (r *MissingWithRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *MissingWithRule) SetValidator(validator Validator) {
	r.validator = validator
}

func (r *MissingWithRule) IsImplicit() bool {
	return true
}

// MissingWithAllRule validates that a field is not present when all other fields are present
type MissingWithAllRule struct {
	Fields []string
	data   map[string]interface{}
	validator Validator
}

func (r *MissingWithAllRule) Passes(attribute string, value interface{}) bool {
	// Check if all specified fields are present
	allFieldsPresent := true
	for _, field := range r.Fields {
		if _, exists := r.data[field]; !exists {
			allFieldsPresent = false
			break
		}
	}
	
	if !allFieldsPresent {
		// Not all dependency fields present, this field can be present or absent
		return true
	}
	
	// All dependency fields are present, this field should be missing
	if r.validator != nil {
		return !r.validator.HasField(attribute)
	}
	return false // If we're here, field is present
}

func (r *MissingWithAllRule) Message() string {
	return "The :attribute field must not be present when all of " + strings.Join(r.Fields, ", ") + " are present."
}

func (r *MissingWithAllRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *MissingWithAllRule) SetValidator(validator Validator) {
	r.validator = validator
}

func (r *MissingWithAllRule) IsImplicit() bool {
	return true
}

// PresentIfRule validates that a field is present when another field equals a value
type PresentIfRule struct {
	Field string
	Value string
	data  map[string]interface{}
	validator Validator
}

func (r *PresentIfRule) Passes(attribute string, value interface{}) bool {
	// Check if the condition field equals the required value
	if fieldValue, exists := r.data[r.Field]; exists {
		if ToString(fieldValue) == r.Value {
			// Field must be present
			if r.validator != nil {
				return r.validator.HasField(attribute)
			}
			return true // If we're here, field is present
		}
	}
	
	// Condition not met, field can be present or absent
	return true
}

func (r *PresentIfRule) Message() string {
	return "The :attribute field must be present when " + r.Field + " is " + r.Value + "."
}

func (r *PresentIfRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *PresentIfRule) SetValidator(validator Validator) {
	r.validator = validator
}

func (r *PresentIfRule) IsImplicit() bool {
	return true
}

// PresentUnlessRule validates that a field is present unless another field equals a value
type PresentUnlessRule struct {
	Field string
	Value string
	data  map[string]interface{}
	validator Validator
}

func (r *PresentUnlessRule) Passes(attribute string, value interface{}) bool {
	// Check if the condition field equals the exception value
	if fieldValue, exists := r.data[r.Field]; exists {
		if ToString(fieldValue) == r.Value {
			// Field can be absent
			return true
		}
	}
	
	// Field must be present
	if r.validator != nil {
		return r.validator.HasField(attribute)
	}
	return true // If we're here, field is present
}

func (r *PresentUnlessRule) Message() string {
	return "The :attribute field must be present unless " + r.Field + " is " + r.Value + "."
}

func (r *PresentUnlessRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *PresentUnlessRule) SetValidator(validator Validator) {
	r.validator = validator
}

func (r *PresentUnlessRule) IsImplicit() bool {
	return true
}

// PresentWithRule validates that a field is present when other fields are present
type PresentWithRule struct {
	Fields []string
	data   map[string]interface{}
	validator Validator
}

func (r *PresentWithRule) Passes(attribute string, value interface{}) bool {
	// Check if any of the specified fields are present
	anyFieldPresent := false
	for _, field := range r.Fields {
		if _, exists := r.data[field]; exists {
			anyFieldPresent = true
			break
		}
	}
	
	if !anyFieldPresent {
		// No dependency fields present, this field can be present or absent
		return true
	}
	
	// At least one dependency field is present, this field must be present
	if r.validator != nil {
		return r.validator.HasField(attribute)
	}
	return true // If we're here, field is present
}

func (r *PresentWithRule) Message() string {
	return "The :attribute field must be present when " + strings.Join(r.Fields, ", ") + " are present."
}

func (r *PresentWithRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *PresentWithRule) SetValidator(validator Validator) {
	r.validator = validator
}

func (r *PresentWithRule) IsImplicit() bool {
	return true
}

// PresentWithAllRule validates that a field is present when all other fields are present
type PresentWithAllRule struct {
	Fields []string
	data   map[string]interface{}
	validator Validator
}

func (r *PresentWithAllRule) Passes(attribute string, value interface{}) bool {
	// Check if all specified fields are present
	allFieldsPresent := true
	for _, field := range r.Fields {
		if _, exists := r.data[field]; !exists {
			allFieldsPresent = false
			break
		}
	}
	
	if !allFieldsPresent {
		// Not all dependency fields present, this field can be present or absent
		return true
	}
	
	// All dependency fields are present, this field must be present
	if r.validator != nil {
		return r.validator.HasField(attribute)
	}
	return true // If we're here, field is present
}

func (r *PresentWithAllRule) Message() string {
	return "The :attribute field must be present when all of " + strings.Join(r.Fields, ", ") + " are present."
}

func (r *PresentWithAllRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *PresentWithAllRule) SetValidator(validator Validator) {
	r.validator = validator
}

func (r *PresentWithAllRule) IsImplicit() bool {
	return true
}

// ProhibitedIfAcceptedRule validates that field is prohibited if another field is accepted
type ProhibitedIfAcceptedRule struct {
	Field string
	data  map[string]interface{}
}

func (r *ProhibitedIfAcceptedRule) Passes(attribute string, value interface{}) bool {
	if fieldValue, exists := r.data[r.Field]; exists {
		// Check if field is accepted
		accepted := false
		switch v := fieldValue.(type) {
		case bool:
			accepted = v
		case string:
			accepted = v == "yes" || v == "on" || v == "1" || v == "true"
		case int, int8, int16, int32, int64:
			accepted = reflect.ValueOf(v).Int() == 1
		case uint, uint8, uint16, uint32, uint64:
			accepted = reflect.ValueOf(v).Uint() == 1
		case float32, float64:
			accepted = reflect.ValueOf(v).Float() == 1
		}
		
		if accepted {
			// Field should be prohibited (missing or empty)
			if IsNil(value) {
				return true
			}
			
			switch v := value.(type) {
			case string:
				return strings.TrimSpace(v) == ""
			case []interface{}, map[string]interface{}:
				return reflect.ValueOf(v).Len() == 0
			default:
				rv := reflect.ValueOf(value)
				switch rv.Kind() {
				case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
					return rv.Len() == 0
				}
			}
			
			return false
		}
	}
	
	return true
}

func (r *ProhibitedIfAcceptedRule) Message() string {
	return "The :attribute field is prohibited when " + r.Field + " is accepted."
}

func (r *ProhibitedIfAcceptedRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *ProhibitedIfAcceptedRule) IsImplicit() bool {
	return true
}

// ProhibitedIfDeclinedRule validates that field is prohibited if another field is declined
type ProhibitedIfDeclinedRule struct {
	Field string
	data  map[string]interface{}
}

func (r *ProhibitedIfDeclinedRule) Passes(attribute string, value interface{}) bool {
	if fieldValue, exists := r.data[r.Field]; exists {
		// Check if field is declined
		declined := false
		switch v := fieldValue.(type) {
		case bool:
			declined = !v
		case string:
			declined = v == "no" || v == "off" || v == "0" || v == "false"
		case int, int8, int16, int32, int64:
			declined = reflect.ValueOf(v).Int() == 0
		case uint, uint8, uint16, uint32, uint64:
			declined = reflect.ValueOf(v).Uint() == 0
		case float32, float64:
			declined = reflect.ValueOf(v).Float() == 0
		}
		
		if declined {
			// Field should be prohibited (missing or empty)
			if IsNil(value) {
				return true
			}
			
			switch v := value.(type) {
			case string:
				return strings.TrimSpace(v) == ""
			case []interface{}, map[string]interface{}:
				return reflect.ValueOf(v).Len() == 0
			default:
				rv := reflect.ValueOf(value)
				switch rv.Kind() {
				case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
					return rv.Len() == 0
				}
			}
			
			return false
		}
	}
	
	return true
}

func (r *ProhibitedIfDeclinedRule) Message() string {
	return "The :attribute field is prohibited when " + r.Field + " is declined."
}

func (r *ProhibitedIfDeclinedRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *ProhibitedIfDeclinedRule) IsImplicit() bool {
	return true
}

// ProhibitedUnlessRule validates that field is prohibited unless another field equals a value
type ProhibitedUnlessRule struct {
	Field string
	Value string
	data  map[string]interface{}
}

func (r *ProhibitedUnlessRule) Passes(attribute string, value interface{}) bool {
	// Check if the condition field equals the exception value
	if fieldValue, exists := r.data[r.Field]; exists {
		if ToString(fieldValue) == r.Value {
			// Field is not prohibited in this case
			return true
		}
	}
	
	// Field should be prohibited (missing or empty)
	if IsNil(value) {
		return true
	}
	
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case []interface{}, map[string]interface{}:
		return reflect.ValueOf(v).Len() == 0
	default:
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
			return rv.Len() == 0
		}
	}
	
	return false
}

func (r *ProhibitedUnlessRule) Message() string {
	return "The :attribute field is prohibited unless " + r.Field + " is " + r.Value + "."
}

func (r *ProhibitedUnlessRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *ProhibitedUnlessRule) IsImplicit() bool {
	return true
}

// ProhibitsRule validates that if field is present, other fields must be prohibited
type ProhibitsRule struct {
	Fields []string
	data   map[string]interface{}
}

func (r *ProhibitsRule) Passes(attribute string, value interface{}) bool {
	// Check if this field is missing or empty
	if IsNil(value) {
		return true
	}
	
	isEmpty := false
	switch v := value.(type) {
	case string:
		isEmpty = strings.TrimSpace(v) == ""
	case []interface{}, map[string]interface{}:
		isEmpty = reflect.ValueOf(v).Len() == 0
	default:
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
			isEmpty = rv.Len() == 0
		}
	}
	
	if isEmpty {
		return true
	}
	
	// This field is present and not empty, check that all prohibited fields are missing or empty
	for _, field := range r.Fields {
		if fieldValue, exists := r.data[field]; exists && !IsNil(fieldValue) {
			// Check if the field is empty
			fieldEmpty := false
			switch v := fieldValue.(type) {
			case string:
				fieldEmpty = strings.TrimSpace(v) == ""
			case []interface{}, map[string]interface{}:
				fieldEmpty = reflect.ValueOf(v).Len() == 0
			default:
				rv := reflect.ValueOf(fieldValue)
				switch rv.Kind() {
				case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
					fieldEmpty = rv.Len() == 0
				}
			}
			
			if !fieldEmpty {
				return false
			}
		}
	}
	
	return true
}

func (r *ProhibitsRule) Message() string {
	return "The :attribute field prohibits " + strings.Join(r.Fields, ", ") + " from being present."
}

func (r *ProhibitsRule) SetData(data map[string]interface{}) {
	r.data = data
}

// RequiredIfAcceptedRule validates that field is required if another field is accepted
type RequiredIfAcceptedRule struct {
	Field string
	data  map[string]interface{}
}

func (r *RequiredIfAcceptedRule) Passes(attribute string, value interface{}) bool {
	if fieldValue, exists := r.data[r.Field]; exists {
		// Check if field is accepted
		accepted := false
		switch v := fieldValue.(type) {
		case bool:
			accepted = v
		case string:
			accepted = v == "yes" || v == "on" || v == "1" || v == "true"
		case int, int8, int16, int32, int64:
			accepted = reflect.ValueOf(v).Int() == 1
		case uint, uint8, uint16, uint32, uint64:
			accepted = reflect.ValueOf(v).Uint() == 1
		case float32, float64:
			accepted = reflect.ValueOf(v).Float() == 1
		}
		
		if accepted {
			// Field is required, check if it's present and not empty
			if IsNil(value) {
				return false
			}
			
			switch v := value.(type) {
			case string:
				return strings.TrimSpace(v) != ""
			case []interface{}, map[string]interface{}:
				return reflect.ValueOf(v).Len() > 0
			default:
				rv := reflect.ValueOf(value)
				switch rv.Kind() {
				case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
					return rv.Len() > 0
				}
			}
			
			return true
		}
	}
	
	return true
}

func (r *RequiredIfAcceptedRule) Message() string {
	return "The :attribute field is required when " + r.Field + " is accepted."
}

func (r *RequiredIfAcceptedRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *RequiredIfAcceptedRule) IsImplicit() bool {
	return true
}

// RequiredIfDeclinedRule validates that field is required if another field is declined
type RequiredIfDeclinedRule struct {
	Field string
	data  map[string]interface{}
}

func (r *RequiredIfDeclinedRule) Passes(attribute string, value interface{}) bool {
	if fieldValue, exists := r.data[r.Field]; exists {
		// Check if field is declined
		declined := false
		switch v := fieldValue.(type) {
		case bool:
			declined = !v
		case string:
			declined = v == "no" || v == "off" || v == "0" || v == "false"
		case int, int8, int16, int32, int64:
			declined = reflect.ValueOf(v).Int() == 0
		case uint, uint8, uint16, uint32, uint64:
			declined = reflect.ValueOf(v).Uint() == 0
		case float32, float64:
			declined = reflect.ValueOf(v).Float() == 0
		}
		
		if declined {
			// Field is required, check if it's present and not empty
			if IsNil(value) {
				return false
			}
			
			switch v := value.(type) {
			case string:
				return strings.TrimSpace(v) != ""
			case []interface{}, map[string]interface{}:
				return reflect.ValueOf(v).Len() > 0
			default:
				rv := reflect.ValueOf(value)
				switch rv.Kind() {
				case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
					return rv.Len() > 0
				}
			}
			
			return true
		}
	}
	
	return true
}

func (r *RequiredIfDeclinedRule) Message() string {
	return "The :attribute field is required when " + r.Field + " is declined."
}

func (r *RequiredIfDeclinedRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *RequiredIfDeclinedRule) IsImplicit() bool {
	return true
}

// RequiredWithAllRule validates that field is required when all other fields are present
type RequiredWithAllRule struct {
	Fields []string
	data   map[string]interface{}
}

func (r *RequiredWithAllRule) Passes(attribute string, value interface{}) bool {
	// Check if all specified fields are present and not empty
	allFieldsPresent := true
	for _, field := range r.Fields {
		if fieldValue, exists := r.data[field]; !exists || IsNil(fieldValue) {
			allFieldsPresent = false
			break
		} else {
			switch v := fieldValue.(type) {
			case string:
				if strings.TrimSpace(v) == "" {
					allFieldsPresent = false
					break
				}
			case []interface{}, map[string]interface{}:
				if reflect.ValueOf(v).Len() == 0 {
					allFieldsPresent = false
					break
				}
			default:
				rv := reflect.ValueOf(fieldValue)
				switch rv.Kind() {
				case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
					if rv.Len() == 0 {
						allFieldsPresent = false
						break
					}
				}
			}
		}
	}
	
	if !allFieldsPresent {
		// Not all fields are present, so this field is not required
		return true
	}
	
	// All fields are present, so this field is required
	if IsNil(value) {
		return false
	}
	
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) != ""
	case []interface{}, map[string]interface{}:
		return reflect.ValueOf(v).Len() > 0
	default:
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
			return rv.Len() > 0
		}
	}
	
	return true
}

func (r *RequiredWithAllRule) Message() string {
	return "The :attribute field is required when all of " + strings.Join(r.Fields, ", ") + " are present."
}

func (r *RequiredWithAllRule) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *RequiredWithAllRule) IsImplicit() bool {
	return true
}

// RequiredArrayKeysRule validates that array contains required keys
type RequiredArrayKeysRule struct {
	Keys []string
}

func (r *RequiredArrayKeysRule) Passes(attribute string, value interface{}) bool {
	// Check if value is a map/object
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Map {
		return false
	}
	
	// Convert to map[string]interface{} if possible
	dataMap, ok := value.(map[string]interface{})
	if !ok {
		// Try to convert map to string keys
		dataMap = make(map[string]interface{})
		for _, key := range rv.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			dataMap[keyStr] = rv.MapIndex(key).Interface()
		}
	}
	
	// Check all required keys exist
	for _, requiredKey := range r.Keys {
		if _, exists := dataMap[requiredKey]; !exists {
			return false
		}
	}
	
	return true
}

func (r *RequiredArrayKeysRule) Message() string {
	return "The :attribute field must contain the keys: " + strings.Join(r.Keys, ", ") + "."
}