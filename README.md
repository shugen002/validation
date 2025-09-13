# Laravel-style Validation Library for Go

An implementation of Laravel's validation framework in Go, providing a comprehensive set of validation rules and features similar to Laravel's validator (excluding database and file validation features).

## Features

This library implements all major Laravel validation rules except database and file/image related ones:

### Basic Type Validation
- `required` - Field must be present and not empty
- `string` - Field must be a string  
- `integer` - Field must be an integer (supports strict mode)
- `numeric` - Field must be numeric (supports strict mode)
- `boolean` - Field must be boolean (supports strict mode)
- `array` - Field must be an array
- `json` - Field must be valid JSON

### String Validation
- `email` - Valid email address
- `alpha` - Only alphabetic characters (supports ASCII mode)
- `alpha_num` - Only alphanumeric characters (supports ASCII mode)
- `alpha_dash` - Only alphanumeric, dash, and underscore characters (supports ASCII mode)
- `regex` - Matches regular expression pattern
- `not_regex` - Does not match regular expression pattern
- `starts_with` - Starts with specified values
- `ends_with` - Ends with specified values
- `doesnt_start_with` - Does not start with specified values
- `doesnt_end_with` - Does not end with specified values
- `uppercase` - All uppercase
- `lowercase` - All lowercase
- `ascii` - Only ASCII characters

### Numeric Validation
- `min` - Minimum value/length
- `max` - Maximum value/length
- `between` - Between two values/lengths
- `size` - Exact size/length
- `gt` - Greater than
- `gte` - Greater than or equal
- `lt` - Less than
- `lte` - Less than or equal

### Array/List Validation
- `in` - Value is in specified list
- `not_in` - Value is not in specified list

### Field Relationship Validation
- `same` - Same as another field
- `different` - Different from other fields
- `confirmed` - Has matching confirmation field

### Date/Time Validation
- `date` - Valid date
- `date_format` - Matches specific date format
- `after` - Date after another date/field
- `before` - Date before another date/field
- `after_or_equal` - Date after or equal to another date/field
- `before_or_equal` - Date before or equal to another date/field
- `date_equals` - Date equals another date/field
- `timezone` - Valid timezone

### Network Validation
- `url` - Valid URL
- `ip` - Valid IP address (IPv4 or IPv6)
- `ipv4` - Valid IPv4 address
- `ipv6` - Valid IPv6 address
- `mac_address` - Valid MAC address

### Other Validation
- `uuid` - Valid UUID (supports version specification)
- `ulid` - Valid ULID
- `hex_color` - Valid hexadecimal color

### Special Rules
- `nullable` - Field can be null
- `sometimes` - Conditional validation
- `bail` - Stop validation on first failure

## Installation

```bash
go get github.com/shugen002/validation
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/shugen002/validation"
)

func main() {
    factory := validation.NewFactory()
    
    data := map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
        "age":   25,
    }
    
    rules := map[string]interface{}{
        "name":  "required|string|min:2",
        "email": "required|email",
        "age":   "integer|min:18|max:100",
    }
    
    validator := factory.Make(data, rules)
    
    if validator.Passes() {
        fmt.Println("Validation passed!")
        fmt.Printf("Valid data: %+v\n", validator.Valid())
    } else {
        fmt.Println("Validation failed!")
        for field, errors := range validator.Errors().All() {
            fmt.Printf("%s: %v\n", field, errors)
        }
    }
}
```

## Usage Examples

### Basic Validation

```go
factory := validation.NewFactory()

data := map[string]interface{}{
    "username": "john_doe",
    "email":    "john@example.com",
    "age":      25,
}

rules := map[string]interface{}{
    "username": "required|alpha_dash|min:3|max:20",
    "email":    "required|email",
    "age":      "required|integer|min:18",
}

validator := factory.Make(data, rules)

if validator.Passes() {
    // Validation passed
    validData := validator.Valid()
    fmt.Printf("Valid data: %+v\n", validData)
} else {
    // Validation failed
    errors := validator.Errors()
    fmt.Printf("Errors: %+v\n", errors.All())
}
```

### Custom Error Messages

```go
customMessages := map[string]string{
    "username.required": "Please provide a username",
    "email.email":       "Please provide a valid email address",
    "age.min":          "You must be at least 18 years old",
}

validator := factory.Make(data, rules, customMessages)
```

### Field Relationships

```go
data := map[string]interface{}{
    "password":              "secret123",
    "password_confirmation": "secret123",
}

rules := map[string]interface{}{
    "password":              "required|min:8",
    "password_confirmation": "required|same:password",
}

validator := factory.Make(data, rules)
```

### Conditional Validation

```go
validator := factory.Make(data, basicRules)

// Add conditional validation
validator.Sometimes("phone", []validation.Rule{
    &validation.RegexRule{Pattern: `^\d{10}$`},
}, func(data map[string]interface{}) bool {
    return data["contact_method"] == "phone"
})
```

### Struct Validation with Tags

```go
type User struct {
    Name     string `json:"name" validate:"required|string|min:2"`
    Email    string `json:"email" validate:"required|email"`
    Age      int    `json:"age" validate:"integer|min:18|max:100"`
    Website  string `json:"website" validate:"url"`
}

user := User{
    Name:    "Jane Doe",
    Email:   "jane@example.com",
    Age:     28,
    Website: "https://janedoe.com",
}

validator, err := factory.ValidateStruct(user)
if err != nil {
    log.Fatal(err)
}

if validator.Passes() {
    fmt.Println("User is valid!")
}
```

### Stop on First Failure

```go
validator := factory.Make(data, rules).StopOnFirstFailure()

if validator.Fails() {
    // Only the first error will be present
    firstError := validator.Errors().First("field_name")
    fmt.Println("First error:", firstError)
}
```

## Rule Reference

### Rule Syntax

Rules can be specified in several ways:

```go
// String with pipe separator
rules := map[string]interface{}{
    "field": "required|string|min:3|max:20",
}

// Array of strings
rules := map[string]interface{}{
    "field": []string{"required", "string", "min:3", "max:20"},
}

// Array of Rule objects
rules := map[string]interface{}{
    "field": []validation.Rule{
        &validation.RequiredRule{},
        &validation.StringRule{},
        &validation.MinRule{Min: 3},
        &validation.MaxRule{Max: 20},
    },
}
```

### Available Rules

#### Basic Type Rules
- `required` - Field must be present and not empty
- `string` - Must be a string
- `integer[:strict]` - Must be an integer (strict mode only accepts actual int types)
- `numeric[:strict]` - Must be numeric (strict mode rejects string numbers)
- `boolean[:strict]` - Must be boolean (strict mode only accepts actual bool types)
- `array[:key1,key2]` - Must be an array (optionally restrict allowed keys)
- `json` - Must be valid JSON string

#### String Rules
- `email` - Valid email address
- `alpha[:ascii]` - Alphabetic characters only (ASCII mode for ASCII-only)
- `alpha_num[:ascii]` - Alphanumeric characters only
- `alpha_dash[:ascii]` - Alphanumeric, dash, and underscore only
- `regex:pattern` - Must match regular expression
- `not_regex:pattern` - Must not match regular expression
- `starts_with:value1,value2` - Must start with one of the values
- `ends_with:value1,value2` - Must end with one of the values
- `doesnt_start_with:value1,value2` - Must not start with any of the values
- `doesnt_end_with:value1,value2` - Must not end with any of the values
- `uppercase` - Must be all uppercase
- `lowercase` - Must be all lowercase
- `ascii` - Must contain only ASCII characters

#### Numeric Rules
- `min:value` - Minimum value/length
- `max:value` - Maximum value/length
- `between:min,max` - Between minimum and maximum values/lengths
- `size:value` - Exact size/length/value
- `gt:field` - Greater than another field
- `gte:field` - Greater than or equal to another field
- `lt:field` - Less than another field
- `lte:field` - Less than or equal to another field

#### List Rules
- `in:value1,value2,value3` - Must be one of the specified values
- `not_in:value1,value2,value3` - Must not be one of the specified values

#### Field Relationship Rules
- `same:field` - Must be the same as another field
- `different:field1,field2` - Must be different from specified fields
- `confirmed[:field]` - Must have matching confirmation field (defaults to fieldname_confirmation)

#### Date Rules
- `date` - Must be a valid date
- `date_format:format1,format2` - Must match one of the specified date formats
- `after:date_or_field` - Must be after the specified date or field value
- `before:date_or_field` - Must be before the specified date or field value
- `after_or_equal:date_or_field` - Must be after or equal to the specified date or field value
- `before_or_equal:date_or_field` - Must be before or equal to the specified date or field value
- `date_equals:date_or_field` - Must equal the specified date or field value
- `timezone[:group[:country]]` - Must be a valid timezone

#### Network Rules
- `url[:protocol1,protocol2]` - Must be a valid URL (optionally restrict protocols)
- `ip` - Must be a valid IP address (IPv4 or IPv6)
- `ipv4` - Must be a valid IPv4 address
- `ipv6` - Must be a valid IPv6 address
- `mac_address` - Must be a valid MAC address

#### Other Rules
- `uuid[:version]` - Must be a valid UUID (optionally specify version 1-8 or "max")
- `ulid` - Must be a valid ULID
- `hex_color` - Must be a valid hexadecimal color

#### Special Rules
- `nullable` - Field can be null/nil
- `sometimes` - Only validate if field is present
- `bail` - Stop validation on first failure

## Error Handling

### Error Bag Methods

```go
errors := validator.Errors()

// Check if field has errors
if errors.Has("email") {
    // Handle email errors
}

// Get all errors for a field
emailErrors := errors.Get("email")

// Get first error for a field
firstError := errors.First("email")

// Get all errors
allErrors := errors.All()

// Check if empty
if errors.IsEmpty() {
    // No errors
}

// Get total error count
count := errors.Count()
```

### Validator Methods

```go
// Check if validation passed
if validator.Passes() {
    // Success
}

// Check if validation failed
if validator.Fails() {
    // Failure
}

// Get valid data only
validData := validator.Valid()

// Get invalid data only
invalidData := validator.Invalid()

// Validate and throw exception on failure
err := validator.Validate()
if err != nil {
    validationErr := err.(*validation.ValidationException)
    fmt.Printf("Validation failed: %s\n", validationErr.Message)
    fmt.Printf("Errors: %+v\n", validationErr.Errors.All())
}
```

## Advanced Features

### Custom Validation Rules

You can create custom validation rules by implementing the `Rule` interface:

```go
type CustomRule struct {
    Parameter string
}

func (r *CustomRule) Passes(attribute string, value interface{}) bool {
    // Your validation logic here
    return true
}

func (r *CustomRule) Message() string {
    return "The :attribute field is invalid."
}

// Use the custom rule
validator.AddRule("field", &CustomRule{Parameter: "value"})
```

### Factory Configuration

```go
factory := validation.NewFactory()

// Set global custom messages
factory.SetCustomMessages(map[string]string{
    "required": "This field is mandatory",
    "email":    "Please provide a valid email address",
})

// Set global custom attributes
factory.SetCustomAttributes(map[string]string{
    "email": "Email Address",
    "phone": "Phone Number",
})
```

## Differences from Laravel

1. **No Database Rules**: `exists`, `unique` and other database-dependent rules are not implemented
2. **No File Rules**: `file`, `image`, `mimes`, `dimensions` and other file-related rules are not implemented
3. **Go Types**: Rules work with Go types instead of PHP types
4. **Struct Tags**: Additional support for validating structs using struct tags

## License

This project is open source and available under the [MIT License](LICENSE.md).
