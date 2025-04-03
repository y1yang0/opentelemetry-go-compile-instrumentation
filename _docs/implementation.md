# Principles
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

# Implementation Details
First, we add new dependencies to the user's project by creating or modifying files. For example:

```go
// otel_importer.go
package main

import _ "github.com/open-telemetry/opentelemetry-go-compile-instrumentation/sdk"
```

The user's project successfully imports our SDK, which contains the hook code. This makes it possible to link the target function and the hook code via golinkname.

Next, execute `go mod tidy` to properly update the go.mod file.

Finally, we invoke `go build -toolexec=otel` to perform the actual build. Our tool intercepts the compilation commands of interest, locates the files containing the target code, and modifies them by injecting trampoline code into the AST as described above. The Hook variable is then linked to the hook code via golinkname.

# Interface Design
The design of the MyHook function is an important topic. Some basic requirements include the ability to retrieve the parameters, return values, and function name of the target function (e.g., Target in the example above) within the MyHook function. An example interface might look like this:

```go
func MyHook(ctx Context) {
	ctx.GetFuncName()
	ctx.GetParam(1)
	ctx.SetParam(1, "new value")
	ctx.GetReturnValue(1)	
	ctx.SetReturnValue(1, "new value")
}
```
All our hook code accepts a Context, through which we can access and modify the data of the target function.

The full API included in Context is as follows:

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