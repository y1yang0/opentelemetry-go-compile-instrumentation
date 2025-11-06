// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"bufio"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func startServerAndWaitForReady(t *testing.T, serverApp *exec.Cmd, outputPipe io.ReadCloser, readyMsg string) {
	t.Helper()
	t.Cleanup(func() {
		if serverApp.Process != nil {
			serverApp.Process.Kill()
		}
	})

	readyChan := make(chan struct{})
	go func() {
		scanner := bufio.NewScanner(outputPipe)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, readyMsg) {
				close(readyChan)
				return
			}
		}
	}()

	select {
	case <-readyChan:
		t.Logf("Server is ready!")
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for server to be ready")
	}
}

func TestHttp(t *testing.T) {
	serverDir := filepath.Join("..", "demo", "http", "server")
	clientDir := filepath.Join("..", "demo", "http", "client")

	// Build the server and client applications with the instrumentation tool.
	BuildApp(t, serverDir, "go", "build", "-a")
	BuildApp(t, clientDir, "go", "build", "-a")

	// Start the server and wait for it to be ready.
	serverApp, outputPipe := StartApp(t, serverDir)
	startServerAndWaitForReady(t, serverApp, outputPipe, "server started")

	// Run the client, it will send a shutdown request to the server.
	RunApp(t, clientDir, "-shutdown")

	// Wait for the server to exit.
	require.NoError(t, serverApp.Wait())
}
