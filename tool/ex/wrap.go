// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ex

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

const numSkipFrame = 3 // skip the Errorf/Fatalf caller

// stackfulError represents an error with stack trace information
type stackfulError struct {
	message string
	frame   []string
	wrapped error
}

func (e *stackfulError) Error() string { return e.message }

func getFrames() []string {
	frameList := make([]string, 0)
	pcs := make([]uintptr, 30)
	n := runtime.Callers(numSkipFrame, pcs[:])
	if n == 0 {
		return frameList
	}
	pcs = pcs[:n]
	frames := runtime.CallersFrames(pcs)
	cnt := 0
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		const prefix = "github.com/open-telemetry/opentelemetry-go-compile-instrumentation/"
		fnName := strings.TrimPrefix(frame.Function, prefix)
		f := fmt.Sprintf("[%d]%s:%d %s", cnt, frame.File, frame.Line, fnName)
		frameList = append(frameList, f)
		cnt++
	}
	return frameList
}

// Error wraps an error with stack trace information
// If you don't want to decorate the existing error, use it.
func Error(previousErr error) error {
	e := &stackfulError{
		message: previousErr.Error(),
		frame:   getFrames(),
		wrapped: previousErr,
	}
	return e
}

// Errorf wraps an error with stack trace information and a formatted message
// If you want to decorate the existing error, use it.
func Errorf(previousErr error, format string, args ...any) error {
	e := &stackfulError{
		message: fmt.Sprintf(format, args...),
		frame:   getFrames(),
		wrapped: previousErr,
	}
	return e
}

func NewError(format string, args ...any) error {
	return Errorf(nil, format, args...)
}

func Fatalf(format string, args ...any) { Fatal(Errorf(nil, format, args...)) }

func Fatal(err error) {
	if err == nil {
		panic("Fatal error: unknown")
	}
	e := &stackfulError{}
	if errors.As(err, &e) {
		msg := fmt.Sprintf("Error:\n%s\n\nStack:\n%s", e.message, strings.Join(e.frame, "\n"))
		_, _ = fmt.Fprint(os.Stderr, msg)
		os.Exit(1)
	}
	panic(err)
}
