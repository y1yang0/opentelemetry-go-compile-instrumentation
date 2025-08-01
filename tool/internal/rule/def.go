// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package rule

import (
	"strings"
)

type Advice struct {
	Before string `json:"before" yaml:"before"`
	After  string `json:"after"  yaml:"after"`
}

type InstRule struct {
	Name     string   `json:"name,omitempty" yaml:"name,omitempty"`
	Path     string   `json:"path"           yaml:"path"`
	Pointcut string   `json:"pointcut"       yaml:"pointcut"`
	Advice   []Advice `json:"advice"         yaml:"advice"`
}

func (r *InstRule) String() string {
	return r.Name
}

func (r *InstRule) GetFuncName() string {
	return strings.Split(r.Pointcut, ".")[1]
}

func (r *InstRule) GetFuncImportPath() string {
	return strings.Split(r.Pointcut, ".")[0]
}
