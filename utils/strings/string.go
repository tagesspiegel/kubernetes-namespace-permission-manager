package strings

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidKeyValueString = errors.New("invalid key value string")
)

// KeyValue parses a string of key=value into key and value
//
// Example:
//
//	key, value, err := KeyValue("foo=bar")
//	if err != nil {
//		// handle error
//	}
//	fmt.Println(key, value) // foo bar
func KeyValue(s string) (key, value string, err error) {
	parts := strings.Split(s, "=")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("%w: %q", ErrInvalidKeyValueString, s)
	}
	return parts[0], parts[1], nil
}

// Properties parses a string of key=value;key=value;key=value into a slice of strings
//
// Example:
//
//	properties := Properties("foo=bar;bar=baz")
//	fmt.Println(properties) // [foo=bar bar=baz]
func Properties(s string) []string {
	return strings.Split(s, ";")
}

// Array parses a string of value,value,value into a slice of strings
//
// Example:
//
//	array := Array("foo,bar,baz")
//	fmt.Println(array) // [foo bar baz]
func Array(s string) []string {
	return strings.Split(s, ",")
}

// RemoveEmpty removes empty strings from a slice of strings
//
// Example:
//
//	array := RemoveEmpty([]string{"foo", "", "bar", "", "baz"})
//	fmt.Println(array) // [foo bar baz]
func RemoveEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// ArrayC parses a string of value::value::value into a slice of strings
//
// Example:
//
//	array := ArrayC("foo::bar::baz")
//	fmt.Println(array) // [foo bar baz]
func ArrayC(s string) []string {
	return strings.Split(s, "::")
}
