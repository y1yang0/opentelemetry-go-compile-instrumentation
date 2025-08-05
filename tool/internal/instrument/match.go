// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package instrument

import "github.com/open-telemetry/opentelemetry-go-compile-instrumentation/tool/internal/rule"

func (*InstrumentPhase) match(args []string, rules []*rule.InstRule) bool {
	// TODO: Implement task
	_ = args
	_ = rules
	return false
}
