# Introduction
This proposal outlines a method for injecting runtime hooks into target functions in Go programs, enabling dynamic monitoring and modification of function behavior. The approach leverages trampoline code injection and function pointer redirection to seamlessly integrate monitoring logic without requiring significant changes to the user's codebase.

The goal is to provide a flexible and non-intrusive mechanism for instrumenting Go applications, particularly for use cases such as observability (e.g., OpenTelemetry), debugging, or performance profiling.

# Core Principles
## 1. Trampoline Code Injection
We inject trampoline code into the Target (lib-side) function, which ultimately jumps to the actual monitoring code via the function pointer Hook.

```go
func Target() {
    Trampoline()
    ....
}

func Trampoline() {
    Hook()
}

//go:linkname Hook github.com/open-telemetry/opentelemetry-go-compile-instrumentation/sdk/hook.MyHook
var Hook func()
```

## 2. Linkage via golinkname
The `Hook` function is linked to the monitoring code using the `//go:linkname` directive. This allows us to dynamically associate the hook with the target function at compile time.

# Implementation Details
## 1. Adding Dependencies
To enable the instrumentation, we first add the necessary dependencies to the user's project. This can be done by creating or modifying files. For example:

```go
// otel_importer.go
package main

import _ "github.com/open-telemetry/opentelemetry-go-compile-instrumentation/sdk"
```

The user's project successfully imports our SDK, which contains the hook code. This makes it possible to link the target function and the hook code via golinkname.

## 2. Updating `go.mod`
Next, execute `go mod tidy` to properly update the go.mod file.

## 3. Building with instrumentation
Finally, invoke the build process with a custom toolchain:

```bash
go build -toolexec=otel
```

The `-toolexec=otel` flag specifies a custom tool (e.g., otel) that intercepts compilation commands. We find the target function from compilation commands and inject trampoline code into the AST. Since the hook dependency is already imported, we can link the target function to the hook code via golinkname without any other modifications.

# Interface Design
## 1. Context
The Context interface is designed to provide a structured way to access and manipulate the parameters, return values, and other relevant data of the target function. This allows the hook code to interact with the target function's execution context seamlessly.

Example:
```go
func MyHook(ctx Context) {
	ctx.GetFuncName()
	ctx.GetParam(1)
	ctx.SetParam(1, "new value")
	ctx.GetReturnValue(1)	
	ctx.SetReturnValue(1, "new value")
}
```

## 2. Full Context API
The full context API is listed below.
```go
type Context interface {
	// Skip the original function call
	SetSkipCall(bool)
	// Check if the original function call should be skipped
	IsSkipCall() bool
	// Set the data field, can be used to pass information between OnEnter & OnExit
	SetData(interface{})
	// Get the data field, can be used to pass information between OnEnter & OnExit
	GetData() interface{}
	// Get the map data field by key
	GetKeyData(key string) interface{}
	// Set the map data field by key
	SetKeyData(key string, val interface{})
	// Has the map data field by key
	HasKeyData(key string) bool
	// Get the original function parameter at index idx
	GetParam(idx int) interface{}
	// Change the original function parameter at index idx
	SetParam(idx int, val interface{})
	// Get the original function return value at index idx
	GetReturnVal(idx int) interface{}
	// Change the original function return value at index idx
	SetReturnVal(idx int, val interface{})
	// Get the original function name
	GetFuncName() string
	// Get the package name of the original function
	GetPackageName() string
}
```

Key Features:
- Parameter Access and Modification: Allows inspection and alteration of function arguments.
- Return Value Manipulation: Enables modification of return values before they are passed back to the caller.
- Flow Control: Provides the ability to skip the original function call entirely (SetSkipCall).
- State Management: Facilitates communication between the OnEnter and OnExit phases of the hook.