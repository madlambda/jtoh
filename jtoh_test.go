package jtoh_test

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/katcipis/jtoh"
)

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
	}

	for i := range tests {
		test := tests[i]
		output := strings.NewReader(strings.Join(test.output, "\n"))

		t.Run(test.name+"WithList", func(t *testing.T) {
			input := strings.NewReader("[" + strings.Join(test.input, ",") + "]")
			testTransform(t, input, test.selector, output, test.wantErr)
		})

		t.Run(test.name+"WithStream", func(t *testing.T) {
			input := strings.NewReader(strings.Join(test.input, "\n"))
			testTransform(t, input, test.selector, output, test.wantErr)
		})
	}
}

func testTransform(
	t *testing.T,
	input io.Reader,
	selector string,
	want io.Reader,
	wantErr error,
) {
	t.Helper()

	got, err := jtoh.Transform(input, selector)

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

	if got == nil {
		t.Errorf("unexpected nil reader as result")
		return
	}

	gotData := readAll(t, got)
	wantData := readAll(t, want)

	if gotData != wantData {
		t.Errorf("got:\n\n%s\n\nwant:\n\n%s\n", gotData, wantData)
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
