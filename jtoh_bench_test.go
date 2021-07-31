package jtoh_test

import (
	"strings"
	"testing"

	"github.com/katcipis/jtoh"
	"github.com/madlambda/spells/iotest"
)

func BenchmarkNonJSONStream10Msgs(b *testing.B) {
	benchmarkNonJSONStream(b, 10)
}

func BenchmarkNonJSONStream100Msgs(b *testing.B) {
	benchmarkNonJSONStream(b, 100)
}

func BenchmarkNonJSONStream1000Msgs(b *testing.B) {
	benchmarkNonJSONStream(b, 1000)
}

func BenchmarkNonJSONStream10000Msgs(b *testing.B) {
	benchmarkNonJSONStream(b, 10000)
}

func benchmarkNonJSONStream(b *testing.B, msgCount int) {
	const testmsg = "non-json-test-msg\n"
	const selector = ":field"

	for i := 0; i < b.N; i++ {
		j, err := jtoh.New(selector)
		if err != nil {
			b.Errorf("unexpected error [%v]", err)
			return
		}

		repeater := iotest.NewRepeatReader(strings.NewReader(testmsg), msgCount)
		j.Do(repeater, NopWriter{})
	}
}

type NopWriter struct{}

func (n NopWriter) Write(a []byte) (int, error) {
	return len(a), nil
}
