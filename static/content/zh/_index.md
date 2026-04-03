---
title: 首页
linkTitle: 首页
description: 面向 Go 的 OpenTelemetry 编译期插桩：无需改业务代码、支持第三方库、信号与 OpenTelemetry 一致。
menu:
  main:
    weight: 10
---

{{% blocks/cover title="OpenTelemetry Go 编译期插桩" image_anchor="top" height="med" color="dark" %}}

{{% param description %}}

<div class="mt-3 mb-4">
<a class="btn btn-lg btn-info me-3 mb-2" href="docs/">文档</a>
<a class="btn btn-lg btn-outline-light mb-2" href="https://github.com/open-telemetry/opentelemetry-go-compile-instrumentation">GitHub</a>
</div>

{{% blocks/link-down color="info" %}}

{{% /blocks/cover %}}

{{% blocks/lead color="primary" %}}

**OpenTelemetry Go 编译期插桩**（`otelc`）在 `go build` 过程中通过 `-toolexec` 注入钩子与遥测数据，与 [OpenTelemetry](https://opentelemetry.io/) 对齐，同时避免在业务仓库中散落插桩代码。

{{% /blocks/lead %}}

{{% blocks/section color="dark" %}}
<div class="home-section">
<h2>工作原理</h2>
<p>
<code>otelc</code> 通过 <code>-toolexec</code> 包装 <code>go build</code>。在 <strong>Setup</strong> 阶段解析插桩依赖；在 <strong>Instrument</strong> 阶段于编译期注入 trampoline 与 OpenTelemetry 对齐的钩子——无需修改任何业务代码。
</p>
<div class="command-example">
<span class="prompt">$</span> otelc go build ./...
</div>
</div>
{{% /blocks/section %}}

{{% blocks/section color="primary" %}}
<div class="home-section">
<h2>支持的集成</h2>
<p style="text-align:center;">对常用 Go 库进行插桩——即使是你不拥有源码的第三方依赖。</p>
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
<h2>资源</h2>
<div class="resource-links">
<a href="docs/getting-started/">快速开始</a>
<a href="/en/docs/implementation/">架构与实现</a>
<a href="docs/rules/">插桩规则</a>
<a href="https://www.youtube.com/watch?v=xEsVOhBdlZY">演讲与视频</a>
</div>
</div>
{{% /blocks/section %}}
