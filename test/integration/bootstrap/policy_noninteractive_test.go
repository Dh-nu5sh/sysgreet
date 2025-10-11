package bootstrap_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestNonInteractiveRequiresPolicy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping policy integration in short mode")
	}

	binaryPath := buildBinary(t)
	cfgDir := t.TempDir()
	cfgPath := filepath.Join(cfgDir, "config.yaml")

	cmd := exec.Command(binaryPath)
	cmd.Env = append(os.Environ(), "CI=1", "SYSGREET_CONFIG="+cfgPath)
	cmd.Stdin = bytes.NewBuffer(nil)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err == nil {
		t.Fatalf("expected error when non-interactive run lacks explicit policy")
	}
	if !strings.Contains(out.String(), "config policy") {
		t.Fatalf("expected guidance about config policy, got %s", out.String())
	}
}

func TestNonInteractiveUsesExplicitPolicy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping policy integration in short mode")
	}

	binaryPath := buildBinary(t)
	cfgDir := t.TempDir()
	cfgPath := filepath.Join(cfgDir, "config.yaml")

	cmd := exec.Command(binaryPath, "--config-policy=keep")
	cmd.Env = append(os.Environ(), "CI=1", "SYSGREET_CONFIG="+cfgPath)
	cmd.Stdin = bytes.NewBuffer(nil)
	var stderr bytes.Buffer
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("expected keep policy to succeed: %v (%s)", err, stderr.String())
	}

	if _, err := os.Stat(cfgPath); err == nil {
		t.Fatalf("expected config file to remain absent when keep policy and non-existent file")
	}

	cmd = exec.Command(binaryPath)
	cmd.Env = append(os.Environ(), "CI=1", "SYSGREET_CONFIG="+cfgPath, "SYSGREET_CONFIG_POLICY=overwrite")
	cmd.Stdin = bytes.NewBuffer(nil)
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("expected overwrite policy to succeed: %v (%s)", err, stderr.String())
	}
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("expected config file to be created under overwrite policy: %v", err)
	}
}
