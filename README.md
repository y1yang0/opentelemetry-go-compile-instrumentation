# OpenTelemetry Go Compile Instrumentation

<img src="./docs/assets/otel-logo.png" alt="OpenTelemetry Logo" width="500">

> [!IMPORTANT]
> This is a work in progress and not ready for production use. 🚨

This project provides a tool to automatically instrument Go applications with [OpenTelemetry](https://opentelemetry.io/) at compile time. It modifies the Go build process to inject OpenTelemetry code into the application without requiring manual changes to the source code.

## Getting Started

1. Build the otel tool
```bash
$ git clone https://github.com/open-telemetry/opentelemetry-go-compile-instrumentation.git
$ cd opentelemetry-go-compile-instrumentation
$ make build
```

2. Build the application with the tool and run it
```bash
$ make demo
```

## Contributing

See the [contributing documentation](CONTRIBUTING.md) for general contribution guidelines.

See the [developing documentation](./docs/developing.md) for tool development.

For the code of conduct, please refer to our [OpenTelemetry Community Code of Conduct](https://github.com/open-telemetry/community/blob/main/code-of-conduct.md)

## License

## License

This project is licensed under the terms of the [Apache Software License version 2.0].
See the [license file](./LICENSE) for more details.
