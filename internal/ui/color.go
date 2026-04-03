package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

var (
	green  = color.New(color.FgGreen)
	red    = color.New(color.FgRed)
	yellow = color.New(color.FgYellow)
	cyan   = color.New(color.FgCyan)
	bold   = color.New(color.Bold)
)

func init() {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		color.NoColor = true
	}
}

func Success(w io.Writer, format string, a ...interface{}) {
	fmt.Fprint(w, green.Sprintf(format, a...))
}

func Successln(w io.Writer, a ...interface{}) {
	fmt.Fprintln(w, green.Sprint(a...))
}

func Error(w io.Writer, format string, a ...interface{}) {
	fmt.Fprint(w, red.Sprintf(format, a...))
}

func Errorln(w io.Writer, a ...interface{}) {
	fmt.Fprintln(w, red.Sprint(a...))
}

func Warn(w io.Writer, format string, a ...interface{}) {
	fmt.Fprint(w, yellow.Sprintf(format, a...))
}

func Warnln(w io.Writer, a ...interface{}) {
	fmt.Fprintln(w, yellow.Sprint(a...))
}

func Info(w io.Writer, format string, a ...interface{}) {
	fmt.Fprint(w, cyan.Sprintf(format, a...))
}

func Infoln(w io.Writer, a ...interface{}) {
	fmt.Fprintln(w, cyan.Sprint(a...))
}

func Boldln(w io.Writer, a ...interface{}) {
	fmt.Fprintln(w, bold.Sprint(a...))
}
