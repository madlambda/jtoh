package jtoh

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// J is a jtoh transformer, it transforms JSON into something more human
type J struct {
	selector string
}

// New creates a new jtoh transformer using the given selector.
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
// If the selector is invalid it returns an error.
func New(selector string) (J, error) {
	// TODO:
	// - selector validation
	return J{
		selector: selector[1:],
	}, nil
}

// Do receives a json stream as input and transforms it
// in simple lines of text (newline-delimited) which is
// then written in the provided writer.
//
// This function will block until all data is read from the input
// and written on the output.
func (j J) Do(jsonInput io.Reader, textOutput io.Writer) {
	jsonInput, _ = isList(jsonInput)
	dec := json.NewDecoder(jsonInput)

	for dec.More() {
		m := map[string]interface{}{}
		err := dec.Decode(&m)
		if err != nil {
			// TODO: handle non disruptive parse errors
			// Ideally we want the original non-JSON data
			// Will need some form of extended reader that remembers
			// part of the read data (not all, don't want O(N) spatial
			// complexity).
			fmt.Printf("TODO:HANDLERR:%v\n", err)
			return
		}

		fmt.Fprint(textOutput, selectField(j.selector, m))
	}
}

func selectField(selector string, doc map[string]interface{}) string {
	v, ok := doc[selector]
	if !ok {
		return ""
	}
	return fmt.Sprint(v)

}

func isList(jsons io.Reader) (io.Reader, bool) {
	buf := make([]byte, 1)

	for {
		_, err := jsons.Read(buf)
		if err != nil {
			return jsons, false
		}

		firstToken := buf[0]
		if isSpace(firstToken) {
			continue
		}

		if firstToken == '[' {
			return jsons, true
		}

		// Got a JSON stream, need to prepend the { back again
		return io.MultiReader(strings.NewReader("{"), jsons), false
	}
}

func isSpace(c byte) bool {
	// TODO: test all this space handling
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}
