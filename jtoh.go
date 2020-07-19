package jtoh

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// J is a jtoh transformer, it transforms JSON into something more human
type J struct {
	separator      string
	fieldSelectors []string
}

// Err is a jtoh error
type Err string

const InvalidSelectorErr Err = "invalid selector"

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
func New(s string) (J, error) {
	selector := []rune(s)
	if len(selector) <= 1 {
		return J{}, fmt.Errorf("%w:%s", InvalidSelectorErr, s)
	}
	separator := string(selector[0])
	return J{
		separator:      separator,
		fieldSelectors: trimSpaces(strings.Split(string(selector[1:]), separator)),
	}, nil
}

// Do receives a json stream as input and transforms it
// in simple lines of text (newline-delimited) which is
// then written in the provided writer.
//
// This function will block until all data is read from the input
// and written on the output.
func (j J) Do(jsonInput io.Reader, textOutput io.Writer) {
	jsonInput, ok := isList(jsonInput)
	dec := json.NewDecoder(jsonInput)

	if ok {
		// WHY: To handle properly gigantic lists of JSON objs
		// Really don't need the return value, but linters can be annoying =P
		_, _ = dec.Token()
	}

	for dec.More() {
		m := map[string]interface{}{}
		err := dec.Decode(&m)
		if err != nil {
			// TODO: handle non disruptive parse errors
			// Ideally we want the original non-JSON data
			// Will need some form of extended reader that remembers
			// part of the read data (not all, don't want O(N) spatial
			// complexity).
			fmt.Fprintf(textOutput, "TODO:HANDLERR:%v\n", err)
			return
		}

		fieldValues := make([]string, len(j.fieldSelectors))
		for i, fieldSelector := range j.fieldSelectors {
			fieldValues[i] = selectField(fieldSelector, m)
		}
		fmt.Fprint(textOutput, strings.Join(fieldValues, j.separator)+"\n")
	}
}

func selectField(selector string, obj map[string]interface{}) string {
	const accessOp = "."

	fields := strings.Split(selector, accessOp)
	pathFields := fields[0 : len(fields)-1]
	finalField := fields[len(fields)-1]

	for _, pathField := range pathFields {
		v, ok := obj[pathField]
		if !ok {
			return missingFieldErrMsg(selector)
		}
		obj, ok = v.(map[string]interface{})
		if !ok {
			return missingFieldErrMsg(selector)
		}
	}

	v, ok := obj[finalField]
	if !ok {
		return missingFieldErrMsg(selector)
	}

	return strings.Trim(fmt.Sprint(v), "\n ")
}

func missingFieldErrMsg(selector string) string {
	return fmt.Sprintf("<jtoh:missing field %q>", selector)
}

func isList(jsons io.Reader) (io.Reader, bool) {
	buf := make([]byte, 1)

	// WHY: was unable to find something like peek on json Decoder
	for {
		_, err := jsons.Read(buf)
		if err != nil {
			// FIXME: Probably would be better to fail here with a more clear error =P
			return jsons, false
		}

		firstToken := buf[0]
		if isSpace(firstToken) {
			continue
		}

		if firstToken == '[' {
			return io.MultiReader(strings.NewReader("["), jsons), true
		}

		if firstToken == '{' {
			return io.MultiReader(strings.NewReader("{"), jsons), false
		}

		// FIXME: Probably would be better to fail here with a more clear error =P
		return jsons, false
	}
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

func (e Err) Error() string {
	return string(e)
}

func trimSpaces(s []string) []string {
	trimmed := make([]string, len(s))
	for i, v := range s {
		trimmed[i] = strings.TrimSpace(v)
	}
	return trimmed
}
