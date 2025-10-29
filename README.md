<div align="center">
  <img src="./docs/assets/otel-logo.png" alt="OpenTelemetry Logo" width="500" />
  <br />
  <img src="https://img.shields.io/badge/Go-1.21%2B-4A90E2?style=flat&logo=go" alt="Go" />
  <img src="https://img.shields.io/badge/License-Apache%202.0-4A90E2?style=flat&logo=apache" alt="License" />
  <img src="https://img.shields.io/badge/Status-Development-FF6B35?style=flat&logo=github" alt="Status" />
  <img src="https://img.shields.io/badge/Slack-CNCF-FF6B35?style=flat&logo=slack" alt="Slack" />
</div>

> [!IMPORTANT]
> This is a work in progress and not ready for production use. 🚨

## Overview

This project provides a tool to automatically instrument Go applications with [OpenTelemetry](https://opentelemetry.io/) at compile-time.
It modifies the Go build process to inject OpenTelemetry code into the application **without requiring manual changes to the source code**.

Highlights:

- **🚀 Zero Runtime Overhead** - Instrumentation is baked into your binary at compile time
- **🔧 Zero Code Changes** - Automatically instrument entire applications and dependencies
- **📚 Third-Party Library Support** - Instrument libraries you don't control
- **🎯 Complete Decoupling** - Keep your codebase free from instrumentation concerns
- **⚡ Flexible Deployment** - Integrate at development time or in your CI/CD pipeline

## Quick Start

### 1. Build the Tool

```bash
git clone https://github.com/open-telemetry/opentelemetry-go-compile-instrumentation.git
cd opentelemetry-go-compile-instrumentation
make build
```

### 2. Try the Demo

```bash
make test
```

## Community

### 📚 Documentation

- [Getting Started Guide](./docs/getting-started.md) - Setup and usage
- [UX Design](./docs/ux-design.md) - Configuration options
- [Implementation Details](./docs/implementation.md) - Technical architecture
- [API Design](./docs/api-design-and-project-structure.md) - API structure

### 🎥 Video Talks

- [Project Overview](https://www.youtube.com/watch?v=xEsVOhBdlZY)
- [Deep Dive Details](https://www.youtube.com/watch?v=8Rw-fVEjihw&list=PLDWZ5uzn69ewrYyHTNrXlrWVDjLiOX0Yb&index=19)

### 💬 Get Help

- [GitHub Discussions](https://github.com/open-telemetry/opentelemetry-go-compile-instrumentation/discussions) - Ask questions
- [GitHub Issues](https://github.com/open-telemetry/opentelemetry-go-compile-instrumentation/issues) - Report bugs
- [Slack Channel](https://cloud-native.slack.com/archives/C088D8GSSSF) - Real-time chat

### 🤝 Contributing

We welcome contributions! See our [contributing guide](CONTRIBUTING.md) and [development docs](./docs/developing.md).

This project follows the [OpenTelemetry Code of Conduct](https://github.com/open-telemetry/community/blob/main/code-of-conduct.md).
