// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package setup

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/open-telemetry/opentelemetry-go-compile-instrumentation/tool/ex"
	"github.com/open-telemetry/opentelemetry-go-compile-instrumentation/tool/util"
)

const (
	goCacheDir = "gocache"
)

type SetupPhase struct {
	logger *slog.Logger
}

func (sp *SetupPhase) Info(msg string, args ...any)  { sp.logger.Info(msg, args...) }
func (sp *SetupPhase) Error(msg string, args ...any) { sp.logger.Error(msg, args...) }
func (sp *SetupPhase) Warn(msg string, args ...any)  { sp.logger.Warn(msg, args...) }
func (sp *SetupPhase) Debug(msg string, args ...any) { sp.logger.Debug(msg, args...) }

// keepForDebug copies the file to the build temp directory for debugging
// Error is tolerated as it's not critical.
func (sp *SetupPhase) keepForDebug(name string) {
	dstFile := filepath.Join(util.GetBuildTemp("debug"), "main", name)
	err := util.CopyFile(name, dstFile)
	if err != nil {
		sp.Warn("failed to record added file", "file", name, "error", err)
	}
}

// This function can be used to check if the setup has been completed.
func isSetup() bool {
	// TODO: Implement Task
	return false
}

// Setup prepares the environment for further instrumentation.
func Setup(ctx context.Context) error {
	logger := util.LoggerFromContext(ctx)

	if isSetup() {
		logger.InfoContext(ctx, "Setup has already been completed, skipping setup.")
		return nil
	}

	sp := &SetupPhase{
		logger: logger,
	}
	// Find all dependencies of the project being build
	deps, err := sp.findDeps(ctx, os.Args[1:])
	if err != nil {
		return err
	}
	// Match the hook code with these dependencies
	matched, err := sp.matchDeps(deps)
	if err != nil {
		return err
	}
	// Introduce additional hook code by generating otel.instrumentation.go
	err = sp.addDeps(matched)
	if err != nil {
		return err
	}
	// Extract the embedded instrumentation modules into local directory
	err = sp.extract()
	if err != nil {
		return err
	}
	// Sync new dependencies to go.mod or vendor/modules.txt
	err = sp.syncDeps(ctx, matched)
	if err != nil {
		return err
	}
	// Write the matched hook to matched.txt for further instrument phase
	err = sp.store(matched)
	if err != nil {
		return err
	}
	sp.Info("Setup completed successfully")
	return nil
}

// BuildWithToolexec builds the project with the toolexec mode
func BuildWithToolexec(ctx context.Context, args []string) error {
	logger := util.LoggerFromContext(ctx)

	// Add -toolexec=otel to the original build command
	execPath, err := os.Executable()
	if err != nil {
		return ex.Wrapf(err, "failed to get executable path")
	}
	insert := "-toolexec=" + execPath + " toolexec"
	const additionalCount = 2
	newArgs := make([]string, 0, len(args)+additionalCount) // Avoid in-place modification
	// Add "go build"
	newArgs = append(newArgs, "go")
	newArgs = append(newArgs, args[:1]...)
	// Add "-work" to give us a chance to debug instrumented code if needed
	newArgs = append(newArgs, "-work")
	// Add "-toolexec=..."
	newArgs = append(newArgs, insert)
	// TODO: We should support incremental build in the future, so we don't need
	// to force rebuild here.
	// Add "-a" to force rebuild
	newArgs = append(newArgs, "-a")
	// Add the rest
	newArgs = append(newArgs, args[1:]...)
	logger.InfoContext(ctx, "Build with toolexec", "args", newArgs)

	// Assemble the environment variables
	env := os.Environ()
	pwd := util.GetOtelWorkDir()
	util.Assert(pwd != "", "invalid working directory")
	// Add the working directory because sub-process should know where we are
	// going to build
	env = append(env, fmt.Sprintf("%s=%s", util.EnvOtelWorkDir, pwd))
	goCachePath, err := filepath.Abs(util.GetBuildTemp(goCacheDir))
	if err != nil {
		return err
	}
	// Add the temp go cache to prevent the instrumented code from using the
	// global go cache and polluting it.
	env = append(env, fmt.Sprintf("GOCACHE=%s", goCachePath))
	logger.InfoContext(ctx, "Build with tolexec", "env", env)

	// Good, run the build command with the toolexec mode
	return util.RunCmdWithEnv(ctx, env, newArgs...)
}

func cleanup(backupFiles []string) {
	_ = os.RemoveAll(OtelRuntimeFile)
	_ = util.RestoreFile(backupFiles)
	_ = os.RemoveAll(util.GetBuildTemp(goCacheDir))
}

func GoBuild(ctx context.Context, args []string) error {
	backupFiles := []string{"go.mod", "go.sum", "go.work", "go.work.sum"}
	_ = util.BackupFile(backupFiles)
	defer cleanup(backupFiles)

	err := Setup(ctx)
	if err != nil {
		return err
	}
	err = BuildWithToolexec(ctx, args)
	if err != nil {
		return err
	}
	return nil
}
