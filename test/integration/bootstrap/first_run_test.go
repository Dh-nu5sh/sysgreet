package bootstrap_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func buildBinary(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "sysgreet")
	cmd := exec.Command("go", "build", "-o", binaryPath, "../../../cmd/sysgreet")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build binary: %v\n%s", err, output)
	}
	return binaryPath
}

type configDoc struct {
	ASCII struct {
		Font string `yaml:"font"`
	} `yaml:"ascii"`
	Version   string `yaml:"version"`
	CreatedAt string `yaml:"created_at"`
}

func TestFirstRunCreatesConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping bootstrap integration in short mode")
	}

	binaryPath := buildBinary(t)
	cfgDir := t.TempDir()
	cfgPath := filepath.Join(cfgDir, "config.yaml")

	cmd := exec.Command(binaryPath)
	cmd.Env = append(os.Environ(), "SYSGREET_CONFIG="+cfgPath, "SYSGREET_ASSUME_TTY=1")
	cmd.Env = append(cmd.Env, "SYSGREET_ASCII_MONOCHROME=true")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	timer := time.AfterFunc(10*time.Second, func() {
		cmd.Process.Kill()
	})
	defer timer.Stop()

	if err := cmd.Run(); err != nil {
		t.Fatalf("sysgreet run failed: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	raw, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("expected config file to exist: %v", err)
	}

	var parsed configDoc
	if err := yaml.Unmarshal(raw, &parsed); err != nil {
		t.Fatalf("yaml parse: %v", err)
	}
	if parsed.ASCII.Font != "ANSI Regular" {
		t.Fatalf("expected ascii font ANSI Regular, got %q", parsed.ASCII.Font)
	}
	if parsed.Version == "" {
		t.Fatalf("expected version to be set")
	}
	if parsed.CreatedAt == "" {
		t.Fatalf("expected created_at to be set")
	}
	if !strings.Contains(stderr.String(), "created default config") {
		t.Fatalf("expected stderr to indicate config creation, got: %s", stderr.String())
	}
}
