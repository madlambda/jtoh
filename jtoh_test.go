package jtoh_test

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/katcipis/jtoh"
)

// TODO:
// JSON stream that has a list inside
// JSON List that has a list inside
// Selector is non-ascii char
// Invalid selectors

func TestTransform(t *testing.T) {
	type Test struct {
		name     string
		selector string
		input    []string
		output   []string
		wantErr  error
	}

	tests := []Test{
		{
			name:     "EmptyInput",
			selector: ":field",
			input:    []string{},
			output:   []string{},
		},
		{
			name:     "SingleSelectStringField",
			selector: ":string",
			input:    []string{`{"string":"lala"}`},
			output:   []string{`lala`},
		},
		{
			name:     "SingleSelectNumberField",
			selector: ":number",
			input:    []string{`{"number":666}`},
			output:   []string{`666`},
		},
		{
			name:     "SingleSelectBoolField",
			selector: ":bool",
			input:    []string{`{"bool":true}`},
			output:   []string{`true`},
		},
		{
			name:     "SingleSelectNullField",
			selector: ":null",
			input:    []string{`{"null":null}`},
			output:   []string{`<nil>`},
		},
		{
			name:     "SingleNestedSelectStringField",
			selector: ":nested.string",
			input:    []string{`{"nested" : { "string":"lala"} }`},
			output:   []string{`lala`},
		},
		{
			name:     "SingleNestedSelectNumberField",
			selector: ":nested.number",
			input:    []string{`{"nested" : { "number":13} }`},
			output:   []string{`13`},
		},
		{
			name:     "MultipleSelectedFields",
			selector: ":string:number:bool",
			input:    []string{`{"string":"hi","number":7,"bool":false}`},
			output:   []string{`hi:7:false`},
		},
		{
			name:     "IncompletePathToField",
			selector: ":nested.number",
			input:    []string{`{"nested" : {} }`},
			output:   []string{missingFieldErrMsg("nested.number")},
		},
		{
			name:     "PathToFieldWithWrongType",
			selector: ":nested.number",
			input:    []string{`{"nested" : "notObj" }`},
			output:   []string{missingFieldErrMsg("nested.number")},
		},
		{
			name:     "UnselectedFieldIsIgnored",
			selector: ":number",
			input:    []string{`{"number":666,"ignored":"hi"}`},
			output:   []string{`666`},
		},
		{
			name:     "MissingField",
			selector: ":missing",
			input:    []string{`{"number":666,"ignored":"hi"}`},
			output:   []string{missingFieldErrMsg("missing")},
		},
		{
			name:     "IgnoreSpacesOnBeginning",
			selector: ":string",
			input:    []string{` {"string":"lala"}`},
			output:   []string{`lala`},
		},
		{
			name:     "IgnoreTabsOnBeginning",
			selector: ":string",
			input: []string{`	{"string":"lala"}`},
			output: []string{`lala`},
		},
		{
			name:     "IgnoreNewlinesOnBeginning",
			selector: ":string",
			input: []string{`
				{"string":"lala"}`,
			},
			output: []string{`lala`},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run(test.name+"WithList", func(t *testing.T) {
			input := strings.NewReader("[" + strings.Join(test.input, ",") + "]")
			testTransform(t, input, test.selector, test.output, test.wantErr)
		})

		t.Run(test.name+"WithStream", func(t *testing.T) {
			input := strings.NewReader(strings.Join(test.input, "\n"))
			testTransform(t, input, test.selector, test.output, test.wantErr)
		})
	}
}

func testTransform(
	t *testing.T,
	input io.Reader,
	selector string,
	want []string,
	wantErr error,
) {
	t.Helper()

	j, err := jtoh.New(selector)

	if wantErr != nil {
		if !errors.Is(err, wantErr) {
			t.Errorf("got err[%v] wanted[%v]", err, wantErr)
		}
		return
	}

	if err != nil {
		t.Errorf("unexpected error [%v]", err)
		return
	}

	output := bytes.Buffer{}

	j.Do(input, &output)

	gotLines := bufio.NewScanner(&output)
	lineCount := 0

	for gotLines.Scan() {
		gotLine := gotLines.Text()
		if lineCount > len(want) {
			t.Errorf("unexpected extra line: %q", gotLine)
			continue
		}
		wantLine := want[lineCount]
		if gotLine != wantLine {
			t.Errorf("line[%d]: got %q != want %q", lineCount, gotLine, wantLine)
		}
		lineCount += 1
	}

	if lineCount != len(want) {
		t.Errorf("got %d lines, want %d", lineCount, len(want))
	}

	if err := gotLines.Err(); err != nil {
		t.Errorf("unexpected error scanning output lines: %v", err)
	}
}

func readAll(t *testing.T, r io.Reader) string {
	t.Helper()

	v, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	return string(v)
}

func missingFieldErrMsg(selector string) string {
	return fmt.Sprintf("<jtoh:missing field %q>", selector)
}
