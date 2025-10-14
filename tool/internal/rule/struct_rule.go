// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package rule

// InstStructRule represents a rule that guides hook struct field injection into
// appropriate target struct locations. For example, if we want to inject
// custom Foo field at the target struct Bar, we can define a rule:
//
//	rule:
//		name: "rule"
//		target: "main"
//		struct: "Bar"
//		field_name: "Foo"
//		field_type: "int"
//
// The rule will be matched against the target struct and the hook field will be
// injected at the appropriate location.
type InstStructRule struct {
	InstBaseRule
	Struct    string `json:"struct"     yaml:"struct"`     // The type name of the struct to be instrumented
	FieldName string `json:"field_name" yaml:"field_name"` // The name of the field to be added
	FieldType string `json:"field_type" yaml:"field_type"` // The type of the field to be added
}
