package iostreams

import (
	"fmt"
	"io"
)

// IOStreams provides structured access to standard input/output/error streams
// with convenience methods for formatted output.
type IOStreams struct {
	// In is the input stream (stdin)
	In io.Reader
	// Out is the output stream (stdout)
	Out io.Writer
	// ErrOut is the error output stream (stderr)
	ErrOut io.Writer
}

// Fprintf writes formatted output to Out with automatic newline.
// If args are provided, the format string is processed with fmt.Sprintf.
// If no args are provided, the format string is written directly.
// Per constitution: automatically appends newline, conditionally uses fmt.Sprintf.
func (s *IOStreams) Fprintf(format string, args ...any) {
	if s.Out == nil {
		// Gracefully handle nil writer - silently ignore
		return
	}

	var message string
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	} else {
		message = format
	}

	_, _ = fmt.Fprintln(s.Out, message)
}

// Fprintln writes output to Out with automatic newline.
// This is a direct pass-through to fmt.Fprintln.
func (s *IOStreams) Fprintln(args ...any) {
	if s.Out == nil {
		// Gracefully handle nil writer - silently ignore
		return
	}

	_, _ = fmt.Fprintln(s.Out, args...)
}

// Errorf writes formatted error output to ErrOut with automatic newline.
// If args are provided, the format string is processed with fmt.Sprintf.
// If no args are provided, the format string is written directly.
// Per constitution: automatically appends newline, conditionally uses fmt.Sprintf.
func (s *IOStreams) Errorf(format string, args ...any) {
	if s.ErrOut == nil {
		// Gracefully handle nil writer - silently ignore
		return
	}

	var message string
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	} else {
		message = format
	}

	_, _ = fmt.Fprintln(s.ErrOut, message)
}

// Errorln writes error output to ErrOut with automatic newline.
// This is a direct pass-through to fmt.Fprintln on the error stream.
func (s *IOStreams) Errorln(args ...any) {
	if s.ErrOut == nil {
		// Gracefully handle nil writer - silently ignore
		return
	}

	_, _ = fmt.Fprintln(s.ErrOut, args...)
}
