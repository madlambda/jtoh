package jtoh_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/madlambda/jtoh"
	"github.com/madlambda/spells/iotest"
)

func BenchmarkNonJSONStreams(b *testing.B) {
	msgCounts := []int{1000, 10000, 100000, 1000000}
	testmsg := "non-json-test-msg\n"
	selector := ":field"

	benchmarkStream(b, selector, testmsg, msgCounts)
}

func BenchmarkJSONStreams(b *testing.B) {
	msgCounts := []int{1000, 10000, 100000, 1000000}
	testmsg := "{ \"field\" : \"bench\" }\n"
	selector := ":field"

	benchmarkStream(b, selector, testmsg, msgCounts)
}

func benchmarkStream(b *testing.B, selector, msg string, msgCounts []int) {
	for _, msgCount := range msgCounts {
		b.Run(fmt.Sprintf("%d Messages", msgCount), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				j, err := jtoh.New(selector)
				if err != nil {
					b.Errorf("unexpected error [%v]", err)
					return
				}

				repeater := iotest.NewRepeatReader(strings.NewReader(msg), msgCount)
				j.Do(repeater, NopWriter{})
			}

		})
	}
}

type NopWriter struct{}

func (n NopWriter) Write(a []byte) (int, error) {
	return len(a), nil
}
