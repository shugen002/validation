# Validation Rules Organization

The validation rules have been reorganized by type into logical categories for better maintainability and discoverability.

## Rule File Organization

### rules_basic_types.go
- **RequiredRule** - Field must be present and not empty
- **StringRule** - Must be a string
- **IntegerRule** - Must be an integer (with strict mode option)
- **NumericRule** - Must be numeric (with strict mode option)
- **BooleanRule** - Must be boolean (with strict mode option)
- **AcceptedRule** - Must be accepted ("yes", "on", 1, "1", true, "true")
- **AcceptedIfRule** - Must be accepted if another field equals specified value
- **DeclinedRule** - Must be declined ("no", "off", 0, "0", false, "false")
- **DeclinedIfRule** - Must be declined if another field equals specified value
- **ArrayRule** - Must be an array/slice (with optional key restrictions)
- **JsonRule** - Must be valid JSON string
- **FilledRule** - Must have a value when present (can be absent)
- **PresentRule** - Must be present in input (can be empty)
- **ProhibitedRule** - Must not be present or must be empty

### rules_string.go
- **EmailRule** - Valid email address
- **AlphaRule** - Alphabetic characters only (with ASCII mode option)
- **AlphaNumRule** - Alphanumeric characters only (with ASCII mode option)
- **AlphaDashRule** - Alphanumeric, dash, and underscore only (with ASCII mode option)
- **RegexRule** - Must match regular expression
- **NotRegexRule** - Must not match regular expression
- **StartsWithRule** - Must start with one of the specified values
- **EndsWithRule** - Must end with one of the specified values
- **DoesntStartWithRule** - Must not start with any of the specified values
- **DoesntEndWithRule** - Must not end with any of the specified values
- **UppercaseRule** - Must be all uppercase
- **LowercaseRule** - Must be all lowercase
- **AsciiRule** - Must contain only ASCII characters

### rules_numeric.go
- **MinRule** - Minimum value/length
- **MaxRule** - Maximum value/length
- **BetweenRule** - Between minimum and maximum values/lengths
- **SizeRule** - Exact size/length/value
- **DigitsRule** - Must be numeric and have exact number of digits
- **DigitsBetweenRule** - Must be numeric and have digit count between min and max
- **MinDigitsRule** - Must have minimum number of digits
- **DecimalRule** - Must have specified decimal places (exact or range)
- **MultipleOfRule** - Must be a multiple of specified value

### rules_list.go
- **InRule** - Must be one of the specified values
- **NotInRule** - Must not be one of the specified values
- **DistinctRule** - Array values must be unique (with strict and ignore_case options)

### rules_relationship.go
- **SameRule** - Must be the same as another field
- **DifferentRule** - Must be different from specified fields
- **ConfirmedRule** - Must have matching confirmation field
- **RequiredIfRule** - Required when another field equals specified value
- **RequiredUnlessRule** - Required unless another field equals specified value
- **RequiredWithRule** - Required when any of other fields are present
- **RequiredWithoutRule** - Required when any of other fields are not present
- **RequiredWithAllRule** - Required when all of other fields are present
- **RequiredIfAcceptedRule** - Required when another field is accepted
- **RequiredIfDeclinedRule** - Required when another field is declined
- **RequiredArrayKeysRule** - Array must contain specified keys

### rules_date.go
- **DateRule** - Must be a valid date
- **DateFormatRule** - Must match specified date formats
- **AfterRule** - Must be after specified date or field value
- **BeforeRule** - Must be before specified date or field value
- **AfterOrEqualRule** - Must be after or equal to specified date or field value
- **BeforeOrEqualRule** - Must be before or equal to specified date or field value
- **DateEqualsRule** - Must equal specified date or field value
- **TimezoneRule** - Must be a valid timezone

### rules_network.go
- **UrlRule** - Must be a valid URL (with optional protocol restrictions)
- **IpRule** - Must be a valid IP address (IPv4 or IPv6)
- **Ipv4Rule** - Must be a valid IPv4 address
- **Ipv6Rule** - Must be a valid IPv6 address
- **MacAddressRule** - Must be a valid MAC address
- **HexColorRule** - Must be a valid hexadecimal color

### rules_special.go
- **UuidRule** - Must be a valid UUID (with optional version specification)
- **UlidRule** - Must be a valid ULID
- **NullableRule** - Indicates field can be null (marker rule)
- **SometimesRule** - Indicates conditional validation (marker rule)
- **BailRule** - Indicates to stop validation on first failure (marker rule)
- **MissingRule** - Must not be present
- **MissingIfRule** - Must not be present when another field equals specified value
- **MissingUnlessRule** - Must not be present unless another field equals specified value
- **MissingWithRule** - Must not be present when any other fields are present
- **MissingWithAllRule** - Must not be present when all other fields are present
- **PresentIfRule** - Must be present when another field equals specified value
- **PresentUnlessRule** - Must be present unless another field equals specified value
- **PresentWithRule** - Must be present when any other fields are present
- **PresentWithAllRule** - Must be present when all other fields are present
- **ProhibitedIfAcceptedRule** - Prohibited when another field is accepted
- **ProhibitedIfDeclinedRule** - Prohibited when another field is declined
- **ProhibitedUnlessRule** - Prohibited unless another field equals specified value
- **ProhibitsRule** - If present, prohibits other fields from being present
- **GreaterThanRule** - Must be greater than another field
- **GreaterThanOrEqualRule** - Must be greater than or equal to another field
- **LessThanRule** - Must be less than another field
- **LessThanOrEqualRule** - Must be less than or equal to another field

## Benefits of This Organization

1. **Better Maintainability** - Related rules are grouped together, making it easier to find and modify specific types of validation logic.

2. **Improved Discoverability** - Developers can quickly locate the rules they need based on the validation type they want to implement.

3. **Logical Separation** - Clear separation between different types of validation concerns (basic types, strings, numbers, relationships, etc.).

4. **Easier Testing** - Rules can be tested by category, making it easier to ensure comprehensive coverage.

5. **Better Documentation** - The file structure itself serves as documentation of the available rule types.

## Migration Notes

- All existing rules have been moved from `rules_basic.go` and `rules_extended.go` to the new categorized files.
- No functionality has been changed - only the organization of the code.
- All tests continue to pass, ensuring backward compatibility.
- The factory.go file continues to work with all rules as before.
- Added missing utility functions (`IsInteger`, `IsNumeric`, `IsJSON`) to `types.go`.