---
title: Home
linkTitle: Home
description: Compile-time OpenTelemetry instrumentation for Go—zero manual code changes, third-party library support, and OpenTelemetry-aligned signals.
menu:
  main:
    weight: 10
---

{{% blocks/cover title="OpenTelemetry Go Compile-time Instrumentation" image_anchor="top" height="med" color="dark" %}}

{{% param description %}}

<div class="mt-3 mb-4">
<a class="btn btn-lg btn-info me-3 mb-2" href="docs/">Documentation</a>
<a class="btn btn-lg btn-outline-light mb-2" href="https://github.com/open-telemetry/opentelemetry-go-compile-instrumentation">GitHub</a>
</div>

{{% blocks/link-down color="info" %}}

{{% /blocks/cover %}}

{{% blocks/lead color="primary" %}}

**OpenTelemetry Go compile-time instrumentation** (`otelc`) injects hooks and telemetry during `go build` using `-toolexec`, so you get traces and metrics aligned with [OpenTelemetry](https://opentelemetry.io/) without scattering instrumentation code through your repository.

{{% /blocks/lead %}}

{{% blocks/section color="dark" %}}
<div class="home-section">
<h2>How it works</h2>
<p>
<code>otelc</code> wraps your <code>go build</code> command via <code>-toolexec</code>. In the <strong>setup</strong> phase it resolves instrumentation dependencies; in the <strong>instrument</strong> phase it injects trampolines and OpenTelemetry-aligned hooks at compile time—no manual code changes required.
</p>
<div class="command-example">
<span class="prompt">$</span> otelc go build ./...
</div>
</div>
{{% /blocks/section %}}

{{% blocks/section color="primary" %}}
<div class="home-section">
<h2>Supported Integrations</h2>
<p style="text-align:center;">Instrument popular Go libraries—even third-party dependencies you don't own.</p>
<div class="integration-grid">
<div class="integration-item">gRPC</div>
<div class="integration-item">net/http</div>
<div class="integration-item">database/sql</div>
<div class="integration-item">Redis</div>
</div>
</div>
{{% /blocks/section %}}

{{% blocks/section color="dark" %}}
<div class="home-section">
<h2>Resources</h2>
<div class="resource-links">
<a href="docs/getting-started/">Getting Started</a>
<a href="docs/implementation/">Architecture</a>
<a href="docs/rules/">Instrumentation Rules</a>
<a href="https://www.youtube.com/watch?v=xEsVOhBdlZY">Talks &amp; Videos</a>
</div>
</div>
{{% /blocks/section %}}
