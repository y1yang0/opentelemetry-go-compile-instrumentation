// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import _ "github.com/open-telemetry/opentelemetry-go-compile-instrumentation/instrumentation/helloworld"

// Example demonstrates how to use the instrumenter.
func Example() {
	// Output:
	// [MyHook] start to instrument hello world!
	// [MyHook] hello world is instrumented!
}

func main() {
	// Call the Example function to trigger the instrumentation
	Example()
}
