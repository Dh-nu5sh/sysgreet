package config

import (
	os "os"
	testing "testing"
)

func TestLoad_DefaultsWhenNoFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cfg, path, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if path != "" {
		t.Fatalf("expected no config path, got %q", path)
	}
	def := Default()
	if cfg.Layout.Compact != def.Layout.Compact || len(cfg.Layout.Sections) != len(def.Layout.Sections) {
		t.Fatalf("expected defaults, got %+v", cfg)
	}
}

func TestLoad_YAMLOverrides(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	cfgDir := dir + "/.config/hostinfo"
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	content := []byte(`display:
  hostname: false
ascii:
  font: "slant"
network:
  max_interfaces: 1
`)
	if err := os.WriteFile(cfgDir+"/config.yaml", content, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, used, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if used == "" {
		t.Fatalf("expected config path to be returned")
	}
	if cfg.Display.Hostname {
		t.Fatalf("expected hostname disabled")
	}
	if cfg.ASCII.Font != "slant" {
		t.Fatalf("expected font override, got %s", cfg.ASCII.Font)
	}
	if cfg.Network.MaxInterfaces != 1 {
		t.Fatalf("expected max interfaces 1, got %d", cfg.Network.MaxInterfaces)
	}
}

func TestLoad_EnvOverrides(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("HOSTINFO_DISPLAY_REMOTE_IP", "false")
	t.Setenv("HOSTINFO_LAYOUT_SECTIONS", "header,resources")
	t.Setenv("HOSTINFO_NETWORK_MAX_INTERFACES", "5")

	cfg, _, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Display.RemoteIP {
		t.Fatalf("env override for remote ip failed")
	}
	if len(cfg.Layout.Sections) != 2 || cfg.Layout.Sections[1] != "resources" {
		t.Fatalf("layout sections override failed: %+v", cfg.Layout.Sections)
	}
	if cfg.Network.MaxInterfaces != 5 {
		t.Fatalf("expected max interfaces 5, got %d", cfg.Network.MaxInterfaces)
	}
}

func TestDefaultSectionsOrder(t *testing.T) {
	cfg := Default()
	if len(cfg.Layout.Sections) == 0 {
		t.Fatalf("expected default sections")
	}
	if cfg.Layout.Sections[0] != "header" {
		t.Fatalf("expected header first, got %s", cfg.Layout.Sections[0])
	}
}
