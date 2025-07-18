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

// Error represents an error with stack trace information
type stackFrameError struct {
	message string
	frame   string
	wrapper error
}

func (e *stackFrameError) Error() string {
	return e.message
}

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
	e := &stackFrameError{}
	if errors.As(err, &e) {
		frame := fmt.Sprintf("[%d] %s\n", cnt, e.frame)
		return fetchFrames(e.wrapper, cnt+1) + frame
	}
	return ""
}

func Error(previousErr error) error {
	e := &stackFrameError{
		message: previousErr.Error(),
		frame:   currentFrame(),
		wrapper: previousErr,
	}
	return e
}

func Errorf(previousErr error, format string, args ...any) error {
	e := &stackFrameError{
		message: fmt.Sprintf(format, args...),
		frame:   currentFrame(),
		wrapper: previousErr,
	}
	return e
}

func Fatalf(format string, args ...any) {
	Fatal(Errorf(nil, format, args...))
}

func Fatal(err error) {
	if err == nil {
		panic("Fatal error: unknown")
	}
	e := &stackFrameError{}
	if errors.As(err, &e) {
		frames := fetchFrames(err, 0)
		msg := fmt.Sprintf("Error:\n%s\n\nStack:\n%s", e.message, frames)
		_, _ = fmt.Fprint(os.Stderr, msg)
		os.Exit(1)
	}
	panic(err)
}
