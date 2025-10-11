package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBinaryExecution(t *testing.T) {
	// Build the binary for testing
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "sysgreet")
	if testing.Short() {
		t.Skip("skipping binary build in short mode")
	}

	buildCmd := exec.Command("go", "build", "-o", binaryPath, "../../cmd/sysgreet")
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build binary: %v\nOutput: %s", err, buildOutput)
	}

	tests := []struct {
		name         string
		args         []string
		env          map[string]string
		wantErr      bool
		wantContains []string
		wantEmpty    bool
	}{
		{
			name:      "disable flag produces no output",
			args:      []string{"--disable"},
			wantEmpty: true,
			wantErr:   false,
		},
		{
			name:         "default execution produces output",
			args:         []string{},
			wantErr:      false,
			wantContains: []string{},
			wantEmpty:    false,
		},
		{
			name: "monochrome mode",
			args: []string{},
			env: map[string]string{
				"SYSGREET_ASCII_MONOCHROME": "true",
			},
			wantErr:   false,
			wantEmpty: false,
		},
		{
			name: "disable specific sections",
			args: []string{},
			env: map[string]string{
				"SYSGREET_DISPLAY_UPTIME": "false",
				"SYSGREET_DISPLAY_MEMORY": "false",
			},
			wantErr:   false,
			wantEmpty: false,
		},
		{
			name: "compact layout",
			args: []string{},
			env: map[string]string{
				"SYSGREET_LAYOUT_COMPACT": "true",
			},
			wantErr:      false,
			wantContains: []string{"|"}, // Compact mode uses pipe separator
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)

			// Use temp config to avoid prompts
			testConfigPath := filepath.Join(tmpDir, "test-config.yaml")

			// Set environment variables
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, "SYSGREET_CONFIG="+testConfigPath, "CI=1", "SYSGREET_CONFIG_POLICY=overwrite")
			for k, v := range tt.env {
				cmd.Env = append(cmd.Env, k+"="+v)
			}

			// Set timeout to prevent hanging
			timeout := time.Second * 5
			timer := time.AfterFunc(timeout, func() {
				if cmd.Process != nil {
					_ = cmd.Process.Kill() // Best effort kill on timeout
				}
			})
			defer timer.Stop()

			output, err := cmd.CombinedOutput()

			if (err != nil) != tt.wantErr {
				t.Errorf("execution error = %v, wantErr %v\nOutput: %s", err, tt.wantErr, output)
				return
			}

			outputStr := string(output)

			if tt.wantEmpty {
				if len(strings.TrimSpace(outputStr)) > 0 {
					t.Errorf("expected empty output, got: %s", outputStr)
				}
				return
			}

			if !tt.wantEmpty && len(strings.TrimSpace(outputStr)) == 0 {
				t.Error("expected non-empty output, got empty")
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(outputStr, want) {
					t.Errorf("output missing %q\nGot: %s", want, outputStr)
				}
			}
		})
	}
}

func TestBinaryStartupTime(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping performance test in short mode")
	}

	// Build the binary for testing
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "sysgreet")

	buildCmd := exec.Command("go", "build", "-o", binaryPath, "../../cmd/sysgreet")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}

	// Run multiple iterations to get average startup time
	iterations := 10
	var totalDuration time.Duration

	for i := 0; i < iterations; i++ {
		start := time.Now()
		cmd := exec.Command(binaryPath, "--disable")
		if err := cmd.Run(); err != nil {
			t.Fatalf("iteration %d failed: %v", i, err)
		}
		duration := time.Since(start)
		totalDuration += duration
	}

	avgDuration := totalDuration / time.Duration(iterations)
	maxAllowed := 80 * time.Millisecond

	if avgDuration > maxAllowed {
		t.Errorf("average startup time %v exceeds maximum allowed %v", avgDuration, maxAllowed)
	}

	t.Logf("Average startup time over %d iterations: %v", iterations, avgDuration)
}

func TestBinaryWithInvalidConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	// Build the binary
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "sysgreet")

	buildCmd := exec.Command("go", "build", "-o", binaryPath, "../../cmd/sysgreet")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}

	// Create invalid config file
	configPath := filepath.Join(tmpDir, "config.yaml")
	invalidConfig := "invalid: yaml: content: [[[["
	if err := os.WriteFile(configPath, []byte(invalidConfig), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Run with invalid config and "keep" policy - should fail to load the invalid YAML
	cmd := exec.Command(binaryPath)
	cmd.Env = append(os.Environ(), "SYSGREET_CONFIG="+configPath, "CI=1", "SYSGREET_CONFIG_POLICY=keep")

	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error with invalid config, got success\nOutput: %s", output)
	}

	// Should contain error message
	if !strings.Contains(string(output), "sysgreet:") {
		t.Errorf("error output should contain 'sysgreet:' prefix\nGot: %s", output)
	}
}

func TestBinaryMemoryFootprint(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	// This test is informational and doesn't fail
	// It helps track memory usage over time
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "sysgreet")

	buildCmd := exec.Command("go", "build", "-o", binaryPath, "../../cmd/sysgreet")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}

	// Get binary size
	info, err := os.Stat(binaryPath)
	if err != nil {
		t.Fatalf("failed to stat binary: %v", err)
	}

	sizeMB := float64(info.Size()) / (1024 * 1024)
	maxSizeMB := 10.0

	if sizeMB > maxSizeMB {
		t.Errorf("binary size %.2f MB exceeds maximum %.2f MB", sizeMB, maxSizeMB)
	}

	t.Logf("Binary size: %.2f MB", sizeMB)
}
