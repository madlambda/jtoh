package jtoh

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Transform received a json stream reader and transforms it
// in a newline separated text stream containing only the fields defined
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
//
// The second reader returned is to be handled as an message error stream.
// If there is a parsing error in the middle of the transformation process
// it will go on writing the parse errors on the returned error stream.
//
// If the returned error is non-nil the stream can be safely ignored,
// otherwise it is the responsibility of the caller to read it
// until EOF (or other error) is reached, failing on do so will leak resources.
func Transform(selector string, jsons io.Reader) (io.Reader, error) {
	// TODO:
	// - selector validation
	// - jsons read errors
	reader, writer := io.Pipe()

	go transform(selector, jsons, writer)

	return reader, nil
}

func transform(selector string, jsons io.Reader, writer *io.PipeWriter) {
	defer writer.Close()

	// WHY: json.Decoder has no peek method :-(
	// We need to detect if this is a json stream or just an gigantic json list.
	jsons, _ = isList(jsons)
	dec := json.NewDecoder(jsons)

	for dec.More() {
		m := map[string]interface{}{}
		err := dec.Decode(&m)
		if err != nil {
			// TODO: handle non disruptive parse errors
			fmt.Fprintf(writer, "jtoh:error:%v", err)
			return
		}
		// TODO: this is obviously wrong
		for _, v := range m {
			fmt.Fprint(writer, v)
		}
	}
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
