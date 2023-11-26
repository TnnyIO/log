package log_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/go-logfmt/logfmt"
	"github.com/tnnyio/log"
)

func TestLogfmtLogger(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	logger := log.NewLogfmtLogger(buf)

	if err := logger.Log("hello", "world"); err != nil {
		t.Fatal(err)
	}
	if want, have := "hello=world\n", buf.String(); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}

	buf.Reset()
	if err := logger.Log("a", 1, "err", errors.New("error")); err != nil {
		t.Fatal(err)
	}
	if want, have := "a=1 err=error\n", buf.String(); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}

	buf.Reset()
	if err := logger.Log("std_map", map[int]int{1: 2}, "my_map", mymap{0: 0}); err != nil {
		t.Fatal(err)
	}
	if want, have := "std_map=\""+logfmt.ErrUnsupportedValueType.Error()+"\" my_map=special_behavior\n", buf.String(); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}
}

func BenchmarkLogfmtLoggerSimple(b *testing.B) {
	benchmarkRunner(b, log.NewLogfmtLogger(io.Discard), baseMessage)
}

func BenchmarkLogfmtLoggerContextual(b *testing.B) {
	benchmarkRunner(b, log.NewLogfmtLogger(io.Discard), withMessage)
}

func TestLogfmtLoggerConcurrency(t *testing.T) {
	t.Parallel()
	testConcurrency(t, log.NewLogfmtLogger(io.Discard), 10000)
}

type mymap map[int]int

func (m mymap) String() string { return "special_behavior" }
