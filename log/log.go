package log

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	verbose bool
	out     io.Writer = os.Stderr
)

func Init() {
	flag.BoolVar(&verbose, "v", false, "enable verbose logging")
}

func Logf(format string, v ...interface{}) {
	if verbose {
		fmt.Fprintf(out, format, v...)
	}
}
