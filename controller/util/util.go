package util

import (
	"reflect"
	"strings"
)

// Stringify gets the string name of any type
// todo test
func Stringify(a any) string {
	return strings.TrimPrefix(
		reflect.TypeOf(a).String(),
		reflect.ValueOf(a).String(),
	)
}
