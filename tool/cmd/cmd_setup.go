// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"

	"github.com/open-telemetry/opentelemetry-go-compile-instrumentation/tool/internal/setup"
	"github.com/urfave/cli/v3"
)

//nolint:gochecknoglobals // Implementation of a CLI command
var commandSetup = cli.Command{
	Name:        "setup",
	Description: "Set up the environment for instrumentation",
	Before:      addLoggerPhaseAttribute,
	Action: func(ctx context.Context, _ *cli.Command) error {
		return setup.Setup(ctx)
	},
}
