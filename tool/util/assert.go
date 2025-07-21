// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package util

import "fmt"

func Assert(condition bool, message string) {
	if !condition {
		panic("Assertion failed: " + message)
	}
}

func ShouldNotReachHere() {
	panic("should not reach here")
}

func Fatal(format string, args ...any) {
	panic("Fatal error: " + fmt.Sprintf(format, args...))
}

func Unimplemented() {
	panic("Unimplemented yet")
}
