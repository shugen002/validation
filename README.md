# Validation

A Go validation library inspired by Laravel's validation rules. This library provides a flexible and extensible way to validate data structures using a variety of built-in rules.

## Features

- **Extensible Rule System**: Easily add custom validation rules
- **Laravel-Inspired Rules**: Supports many rules similar to Laravel's validation
- **Type-Safe**: Written in Go with strong typing
- **Memory Context**: Rules can share context through memory for complex validations

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

    // Define validation rules
    rules := map[string]string{
        "email":    "required|email",
        "password": "required|min:8",
        "age":      "numeric|min:18",
    }

    // Create validator
    validator, err := factory.Parse(rules)
    if err != nil {
        panic(err)
    }

    // Data to validate
    data := map[string]string{
        "email":    "user@example.com",
        "password": "securepassword",
        "age":      "25",
    }

    // Validate
    err = validator.Validate(data)
    if err != nil {
        fmt.Printf("Validation failed: %v\n", err)
    } else {
        fmt.Println("Validation passed!")
    }
}
```

## Supported Rules

### String Rules

- `alpha` - Field must be entirely alphabetic characters
- `alpha_dash` - Field must be alpha-numeric with dashes and underscores
- `alpha_num` - Field must be alpha-numeric
- `ascii` - Field must be ASCII characters
- `confirmed` - Field must have a matching `{field}_confirmation`
- `different:field` - Field must differ from another field
- `email` - Field must be a valid email address
- `ends_with:foo,bar` - Field must end with one of the values
- `hex_color` - Field must be a valid hex color
- `in:foo,bar` - Field must be in the given list
- `ip` - Field must be a valid IP address
- `ipv4` - Field must be a valid IPv4 address
- `ipv6` - Field must be a valid IPv6 address
- `json` - Field must be valid JSON
- `lowercase` - Field must be lowercase
- `mac_address` - Field must be a valid MAC address
- `not_in:foo,bar` - Field must not be in the given list
- `regex:pattern` - Field must match regex pattern
- `same:field` - Field must match another field
- `starts_with:foo,bar` - Field must start with one of the values
- `string` - Field must be a string
- `ulid` - Field must be a valid ULID
- `uppercase` - Field must be uppercase
- `url` - Field must be a valid URL
- `uuid` - Field must be a valid UUID

### Number Rules

- `numeric` - Field must be numeric
- `integer` - Field must be an integer
- `decimal:min,max` - Field must have specified decimal places
- `digits:value` - Field must be exactly N digits
- `digits_between:min,max` - Field must be between min and max digits
- `min_digits:value` - Field must have at least N digits
- `max_digits:value` - Field must have at most N digits

### Size Rules

- `min:value` - Field must be at least value
- `max:value` - Field must be at most value
- `size:value` - Field must be exactly value
- `between:min,max` - Field must be between min and max
- `gt:field_or_value` - Field must be greater than another field or value
- `gte:field_or_value` - Field must be greater than or equal to another field or value
- `lt:field_or_value` - Field must be less than another field or value
- `lte:field_or_value` - Field must be less than or equal to another field or value

### Boolean Rules

- `accepted` - Field must be "yes", "on", 1, "1", true, or "true"
- `accepted_if:anotherfield,value,...` - Field must be accepted if another field equals specified value
- `boolean` - Field must be a boolean (true/false, 1/0, "1"/"0")
- `boolean:strict` - Field must be strictly true or false
- `declined` - Field must be "no", "off", 0, "0", false, or "false"
- `declined_if:anotherfield,value,...` - Field must be declined if another field equals specified value

### Utility Rules

- `required` - Field must be present and not empty
- `nullable` - Field may be null
- `sometimes` - Field may be present

## Custom Rules

You can register custom validation rules:

```go
factory := validation.NewFactory()

factory.RegisterRule("custom_rule", func(cfg map[string]interface{}, args ...string) (validation.ValidationRule, error) {
    return func(ctx *validation.ValidationContext) (bool, error) {
        // Your validation logic here
        if ctx.FieldValue != "expected" {
            return false, fmt.Errorf("field %s must be 'expected'", ctx.FieldName)
        }
        return true, nil
    }, nil
})
```

## Configuration

You can set global configuration:

```go
factory.SetConfig("strict", true)
```

## Testing

Run tests with:

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
