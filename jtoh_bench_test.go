package jtoh_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/katcipis/jtoh"
	"github.com/madlambda/spells/iotest"
)

func BenchmarkNonJSONStreams(b *testing.B) {
	msgCounts := []int{10, 100, 1000, 10000}

	for _, msgCount := range msgCounts {
		b.Run(fmt.Sprintf("%d Messages", msgCount), func(b *testing.B) {
			const testmsg = "non-json-test-msg\n"
			const selector = ":field"

			benchmarkStream(b, selector, testmsg, msgCount)
		})
	}
}

func benchmarkStream(b *testing.B, selector, msg string, msgCount int) {
	for i := 0; i < b.N; i++ {
		j, err := jtoh.New(selector)
		if err != nil {
			b.Errorf("unexpected error [%v]", err)
			return
		}

		repeater := iotest.NewRepeatReader(strings.NewReader(msg), msgCount)
		j.Do(repeater, NopWriter{})
	}
}

type NopWriter struct{}

func (n NopWriter) Write(a []byte) (int, error) {
	return len(a), nil
}
