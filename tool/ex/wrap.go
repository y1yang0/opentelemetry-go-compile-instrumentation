// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ex

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// stackfulError represents an error with stack trace information
type stackfulError struct {
	message string
	frame   string
	wrapped error
}

func (e *stackfulError) Error() string { return e.message }

// currentFrame returns the "current frame" whose caller is the function that
// called Errorf.
func currentFrame() string {
	pc := make([]uintptr, 1)
	const skip = 3 // skip the Errorf caller
	n := runtime.Callers(skip, pc)
	if n == 0 {
		return ""
	}
	pc = pc[:n]
	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()
	shortFunc := frame.Function
	const prefix = "github.com/open-telemetry/opentelemetry-go-compile-instrumentation/"
	shortFunc = strings.TrimPrefix(shortFunc, prefix)
	return frame.File + ":" + strconv.Itoa(frame.Line) + " " + shortFunc
}

func fetchFrames(err error, cnt int) string {
	e := &stackfulError{}
	if errors.As(err, &e) {
		frame := fmt.Sprintf("[%d] %s\n", cnt, e.frame)
		return fetchFrames(e.wrapped, cnt+1) + frame
	}
	return ""
}

// Error wraps an error with stack trace information
// If you don't want to decorate the existing error, use it.
func Error(previousErr error) error {
	e := &stackfulError{
		message: previousErr.Error(),
		frame:   currentFrame(),
		wrapped: previousErr,
	}
	return e
}

// Errorf wraps an error with stack trace information and a formatted message
// If you want to decorate the existing error, use it.
func Errorf(previousErr error, format string, args ...any) error {
	e := &stackfulError{
		message: fmt.Sprintf(format, args...),
		frame:   currentFrame(),
		wrapped: previousErr,
	}
	return e
}

func Fatalf(format string, args ...any) { Fatal(Errorf(nil, format, args...)) }

func Fatal(err error) {
	if err == nil {
		panic("Fatal error: unknown")
	}
	e := &stackfulError{}
	if errors.As(err, &e) {
		frames := fetchFrames(err, 0)
		msg := fmt.Sprintf("Error:\n%s\n\nStack:\n%s", e.message, frames)
		_, _ = fmt.Fprint(os.Stderr, msg)
		os.Exit(1)
	}
	panic(err)
}
