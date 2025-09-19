package validation

import (
	"fmt"
	"reflect"
	"strings"
)

// List validation rules

// InRule validates that a field is one of the specified values
type InRule struct {
	Values []string
}

func (r *InRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	for _, v := range r.Values {
		if str == v {
			return true
		}
	}
	return false
}

func (r *InRule) Message() string {
	return "The selected :attribute is invalid."
}

// NotInRule validates that a field is not one of the specified values
type NotInRule struct {
	Values []string
}

func (r *NotInRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	for _, v := range r.Values {
		if str == v {
			return false
		}
	}
	return true
}

func (r *NotInRule) Message() string {
	return "The selected :attribute is invalid."
}

// DistinctRule validates that array values are unique
type DistinctRule struct {
	Strict     bool
	IgnoreCase bool
}

func (r *DistinctRule) Passes(attribute string, value interface{}) bool {
	rv := reflect.ValueOf(value)
	
	// Must be an array or slice
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return false
	}
	
	seen := make(map[string]bool)
	
	for i := 0; i < rv.Len(); i++ {
		item := rv.Index(i).Interface()
		key := r.getKey(item)
		
		if seen[key] {
			return false
		}
		seen[key] = true
	}
	
	return true
}

func (r *DistinctRule) getKey(value interface{}) string {
	var str string
	
	if r.Strict {
		// Include type information for strict comparison
		str = fmt.Sprintf("%T:%v", value, value)
	} else {
		str = fmt.Sprintf("%v", value)
	}
	
	if r.IgnoreCase {
		str = strings.ToLower(str)
	}
	
	return str
}

func (r *DistinctRule) Message() string {
	return "The :attribute field has duplicate values."
}