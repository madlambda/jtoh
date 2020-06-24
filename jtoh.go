package jtoh

import (
	"io"
	"strings"
)

// Transform received a json stream reader and transforms it
// in a newline separated text with only the fields defined
// by the given selector and in the same order.
//
// The selector is on the form <separator><field selector 1><separator><field selector 2>
// For example, given ":" as a separator you can define:
//
// :fieldA:fieldB:fieldC
//
// Accessing a nested field is done with dot notation, like this:
//
// :field.nested
//
// Making "." the only character that will not be allowed to be used
// as a separator since it is already a selector for nested fields.
//
// If the jsons reader returns a non-nil non-EOF error the error
// will also be returned on the transformed reader Read call.
func Transform(jsons io.Reader, selector string) (io.Reader, error) {
	return strings.NewReader(""), nil
}
