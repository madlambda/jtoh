package jtoh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// J is a jtoh transformer, it transforms JSON into something more human
type J struct {
	separator      string
	fieldSelectors []string
}

// Err is an exported jtoh error
type Err string

// InvalidSelectorErr represents errors with the provided fields selector
const InvalidSelectorErr Err = "invalid selector"

// New creates a new jtoh transformer using the given selector.
// The selector is on the form <separator><field selector 1><separator><field selector 2>
// For example, given ":" as a separator you can define:
//
// :fieldA:fieldB:fieldC
//
// Accessing a nested field is done with dot to access nested fields, like this:
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
	if separator == "." {
		return J{}, fmt.Errorf("%w:can't use '.' as separator", InvalidSelectorErr)
	}
	return J{
		separator:      separator,
		fieldSelectors: trimSpaces(strings.Split(string(selector[1:]), separator)),
	}, nil
}

// Do receives a json stream as input and transforms it
// in lines of text (newline-delimited) which is
// then written in the provided writer.
//
// This function will block until all data is read from the input
// and written on the output.
func (j J) Do(jsonInput io.Reader, linesOutput io.Writer) {
	jsonInput, ok := isList(jsonInput)
	// Why not bufio ? what we need here is kinda like
	// buffered io, but not exactly the same (was not able to
	// come up with a better name to it).
	bufinput := bufferedReader{r: jsonInput}
	dec := json.NewDecoder(&bufinput)

	if ok {
		// WHY: To handle properly gigantic lists of JSON objs
		// Really don't need the return value, but linters can be annoying =P
		_, _ = dec.Token()
	}

	var errBuffer []byte

	for dec.More() {
		m := map[string]interface{}{}
		err := dec.Decode(&m)
		dataUsedOnDecode := bufinput.readBuffer()
		bufinput.reset()

		if err != nil {
			errBuffer = append(errBuffer, dataUsedOnDecode...)
			dec = json.NewDecoder(&bufinput)
			continue
		}

		writeErrs(linesOutput, errBuffer)
		errBuffer = nil

		fieldValues := make([]string, len(j.fieldSelectors))
		for i, fieldSelector := range j.fieldSelectors {
			fieldValues[i] = selectField(fieldSelector, m)
		}
		fmt.Fprint(linesOutput, strings.Join(fieldValues, j.separator)+"\n")
	}

	writeErrs(linesOutput, errBuffer)
}

func writeErrs(w io.Writer, errBuffer []byte) {
	if len(errBuffer) == 0 {
		return
	}

	errBuffer = append(errBuffer, '\n')
	// TODO: handle write errors
	n, err := w.Write(errBuffer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "jtoh:error writing error buffer: wrote %d bytes, details: %v\n", n, err)
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
		n, err := jsons.Read(buf)
		if err != nil {
			// FIXME: Probably would be better to fail here with a more clear error =P
			return jsons, false
		}
		if n == 0 {
			// From the docs:
			//
			// https://golang.org/pkg/io/#Reader
			//
			// Implementations of Read are discouraged from
			// returning a zero byte count with a nil error,
			// except when len(p) == 0. Callers should treat a
			// return of 0 and nil as indicating that nothing happened;
			// in particular it does not indicate EOF.
			//
			// Hope it doesn't result in some infinite loop =/
			continue
		}

		firstToken := buf[0]
		if isSpace(firstToken) {
			continue
		}

		isList := firstToken == '['
		return io.MultiReader(bytes.NewBuffer([]byte{firstToken}), jsons), isList
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

// bufferedReader is not exactly like the bufio on stdlib.
// The idea is to use it as a means to buffer read data
// until reset is called. We need this so when
// the JSON decoder finds an error in the stream we can retrieve
// exactly how much has been read between the last successful
// decode and the current error and echo it.
//
// To guarantee that we provide data byte per byte, which is
// not terribly efficient but was the only way so far to be sure
// (assuming that the json decoder does no lookahead) that when
// an error occurs on the json decoder we have the exact byte stream that
// caused the error (I would welcome with open arms a better solution x_x).
type bufferedReader struct {
	r      io.Reader
	buffer []byte
}

func (b *bufferedReader) Read(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}

	data = data[:1]
	n, err := b.r.Read(data)

	if n > 0 {
		b.buffer = append(b.buffer, data[0])
	}

	return n, err
}

func (b *bufferedReader) readBuffer() []byte {
	return b.buffer
}

func (b *bufferedReader) reset() {
	b.buffer = make([]byte, 0, 1024)
}
