// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

type MyStruct struct{}

func (m *MyStruct) Example() { println("MyStruct.Example") }

// Example demonstrates how to use the instrumenter.
func Example() {
	// Output:
	// [MyHook] start to instrument hello world!
	// [MyHook] hello world is instrumented!
}

func main() {
	// Call the Example function to trigger the instrumentation
	Example()
	m := &MyStruct{}
	// Add a new field to the struct
	println(m.NewField)
	m.Example()
}
