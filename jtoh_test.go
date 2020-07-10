package jtoh_test

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/katcipis/jtoh"
)

// TODO:
// Test JSON stream that has a list inside
// Test JSON List that has a list inside

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
			name:     "SingleSelectIntField",
			selector: ":int",
			input:    []string{`{"int":666}`},
			output:   []string{`666`},
		},
		{
			name:     "UnselectedFieldIsIgnored",
			selector: ":int",
			input:    []string{`{"int":666,"ignored":"hi"}`},
			output:   []string{`666`},
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
