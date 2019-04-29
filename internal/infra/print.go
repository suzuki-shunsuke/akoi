package infra

import (
	"fmt"
	"io"
)

type (
	// Printer implements domain.Printer .
	Printer struct {
		DryRun bool
	}
)

// Fprintf implements domain.Printer .
func (printer Printer) Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	if printer.DryRun {
		return 0, nil
	}
	return fmt.Fprintf(w, format, a...)
}

// Fprintln implements domain.Printer .
func (printer Printer) Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	if printer.DryRun {
		return 0, nil
	}
	return fmt.Fprintln(w, a...)

}

// Printf implements domain.Printer .
func (printer Printer) Printf(format string, a ...interface{}) (n int, err error) {
	if printer.DryRun {
		return 0, nil
	}
	return fmt.Printf(format, a...)
}

// Println implements domain.Printer .
func (printer Printer) Println(a ...interface{}) (n int, err error) {
	if printer.DryRun {
		return 0, nil
	}
	return fmt.Println(a...)
}
