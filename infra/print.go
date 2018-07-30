package infra

import (
	"fmt"
	"io"

	"github.com/suzuki-shunsuke/akoi/domain"
)

// NewFprintf returns a domain.Fprintf which prints if flag is true.
func NewFprintf(flag bool) domain.Fprintf {
	return func(w io.Writer, format string, a ...interface{}) (n int, err error) {
		if flag {
			return fmt.Fprintf(w, format, a...)
		}
		return 0, nil
	}
}

// NewFprintln returns a domain.Fprintln which prints if flag is true.
func NewFprintln(flag bool) domain.Fprintln {
	return func(w io.Writer, a ...interface{}) (n int, err error) {
		if flag {
			return fmt.Fprintln(w, a...)
		}
		return 0, nil
	}
}

// NewPrintf returns a domain.Printf which prints if flag is true.
func NewPrintf(flag bool) domain.Printf {
	return func(format string, a ...interface{}) (n int, err error) {
		if flag {
			return fmt.Printf(format, a...)
		}
		return 0, nil
	}
}

// NewPrintln returns a domain.Println which prints if flag is true.
func NewPrintln(flag bool) domain.Println {
	return func(a ...interface{}) (n int, err error) {
		if flag {
			return fmt.Println(a...)
		}
		return 0, nil
	}
}
