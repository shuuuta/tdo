package log

import (
	"bytes"
	"flag"
	"testing"
)

func TestLogf(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
	Init()
	verbose = true

	buf := &bytes.Buffer{}
	out = buf

	Logf("%s %s", "sample", "test")
	got := buf.String()
	if got != "sample test" {
		t.Fatalf("expect 'sample test', got %q", got)
	}
}

func TestLogfSilent(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
	Init()
	verbose = false

	buf := &bytes.Buffer{}
	out = buf

	Logf("%s %s", "sample", "test")
	if buf.Len() != 0 {
		t.Fatalf("expect no output, got %q", buf.String())
	}
}
