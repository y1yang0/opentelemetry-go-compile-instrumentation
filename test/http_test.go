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

func startServerAndWaitForReady(t *testing.T, serverApp *exec.Cmd, outputPipe io.ReadCloser, readyMsg string) func() string {
	t.Helper()

	readyChan := make(chan struct{})
	doneChan := make(chan struct{})
	output := strings.Builder{}
	go func() {
		// Scan will return false when the application exits.
		defer close(doneChan)
		scanner := bufio.NewScanner(outputPipe)
		for scanner.Scan() {
			line := scanner.Text()
			output.WriteString(line + "\n")
			if strings.Contains(line, readyMsg) {
				close(readyChan)
			}
		}
	}()

	select {
	case <-readyChan:
		t.Logf("Server is ready!")
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for server to be ready")
	}

	return func() string {
		// Wait for the server to exit
		serverApp.Wait()
		// Wait for the output goroutine to finish
		<-doneChan
		// Return the complete output
		return output.String()
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
	waitServerOutput := startServerAndWaitForReady(t, serverApp, outputPipe, "server started")

	// Run the client, it will send a shutdown request to the server.
	RunApp(t, clientDir, "-shutdown")

	// Wait for the server to exit and return the output.
	output := waitServerOutput()

	// Verify that the server hook was called.
	require.Contains(t, output, "BeforeServeHTTP")
}
