// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"path/filepath"
)

const (
	BuildTempDir = ".otel-build"
	OtelRoot     = "github.com/open-telemetry/opentelemetry-go-compile-instrumentation"
)

// GetBuildTemp returns the path to the build temp directory $BUILD_TEMP/name
func GetBuildTemp(name string) string {
	return filepath.Join(BuildTempDir, name)
}

// BackupFile backups the source file to $BUILD_TEMP/backup/name, error is
// tolerated as it's not critical.
func copyFilesHelper(names []string, srcRoot, dstRoot string) {
	for _, name := range names {
		src := filepath.Join(srcRoot, name)
		dst := filepath.Join(dstRoot, name)
		_ = CopyFile(src, dst) // Error is tolerated
	}
}

func BackupFile(names []string) {
	backupDir := GetBuildTemp("backup")
	copyFilesHelper(names, "", backupDir)
}

func RestoreFile(names []string) {
	backupDir := GetBuildTemp("backup")
	copyFilesHelper(names, backupDir, "")
}
	for _, name := range names {
		src := name
		dst := filepath.Join(GetBuildTemp("backup"), name)
		_ = CopyFile(src, dst)
	}
}

// RestoreFile restores the source file from $BUILD_TEMP/backup/name, error is
// tolerated as it's not critical.
func RestoreFile(names []string) {
	for _, name := range names {
		src := filepath.Join(GetBuildTemp("backup"), name)
		dst := name
		_ = CopyFile(src, dst)
	}
}
