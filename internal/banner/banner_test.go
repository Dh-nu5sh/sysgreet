package banner

import (
	"context"
	"testing"
	"time"

	"github.com/veteranbv/hostinfo/internal/ascii"
	"github.com/veteranbv/hostinfo/internal/collectors"
	"github.com/veteranbv/hostinfo/internal/config"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		renderer    *ascii.Renderer
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid renderer",
			renderer: mustRenderer(t),
			wantErr:  false,
		},
		{
			name:        "nil renderer",
			renderer:    nil,
			wantErr:     true,
			errContains: "ascii renderer is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			providers := collectors.Providers{}
			banner, err := New(providers, tt.renderer, nil)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && err.Error() != tt.errContains {
					t.Errorf("error = %q, want %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if banner == nil {
				t.Fatal("expected banner, got nil")
			}
		})
	}
}

func TestBanner_Build(t *testing.T) {
	tests := []struct {
		name         string
		providers    collectors.Providers
		builders     []Builder
		cfg          config.Config
		wantSections int
	}{
		{
			name:         "empty providers and builders",
			providers:    collectors.Providers{},
			builders:     []Builder{},
			cfg:          config.Default(),
			wantSections: 0,
		},
		{
			name: "with system section builder",
			providers: collectors.Providers{
				System: &mockSystemCollector{
					info: collectors.SystemInfo{
						Hostname:    "testhost",
						OS:          "Linux",
						OSVersion:   "6.0",
						Arch:        "x86_64",
						Uptime:      time.Hour * 24,
						CurrentUser: "testuser",
						HomeDir:     "/home/testuser",
						Datetime:    time.Now(),
					},
				},
			},
			builders: []Builder{
				SystemSectionBuilder{},
			},
			cfg:          config.Default(),
			wantSections: 1,
		},
		{
			name: "with disabled sections in config",
			providers: collectors.Providers{
				System: &mockSystemCollector{
					info: collectors.SystemInfo{
						Hostname: "testhost",
						Uptime:   time.Hour * 24,
					},
				},
			},
			builders: []Builder{
				SystemSectionBuilder{},
			},
			cfg: config.Config{
				Display: config.DisplayConfig{
					Uptime:   false,
					User:     false,
					Datetime: false,
				},
			},
			wantSections: 0,
		},
		{
			name: "with nil builder in list",
			providers: collectors.Providers{
				System: &mockSystemCollector{
					info: collectors.SystemInfo{
						Hostname: "testhost",
						Uptime:   time.Hour * 24,
					},
				},
			},
			builders: []Builder{
				nil,
				SystemSectionBuilder{},
			},
			cfg:          config.Default(),
			wantSections: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			banner, err := New(tt.providers, mustRenderer(t), tt.builders)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			output, snap, err := banner.Build(context.Background(), tt.cfg)
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			if len(output.Sections) != tt.wantSections {
				t.Errorf("got %d sections, want %d", len(output.Sections), tt.wantSections)
			}

			// Verify snapshot is populated
			if tt.providers.System != nil {
				if snap.System.Hostname == "" {
					t.Error("expected hostname in snapshot")
				}
			}
		})
	}
}

func TestBanner_buildHeader(t *testing.T) {
	tests := []struct {
		name         string
		hostname     string
		os           string
		osVersion    string
		arch         string
		displayOS    bool
		wantHostname string
		wantLines    int
	}{
		{
			name:         "valid hostname with OS",
			hostname:     "myhost",
			os:           "Linux",
			osVersion:    "6.0",
			arch:         "x86_64",
			displayOS:    true,
			wantHostname: "myhost",
			wantLines:    1,
		},
		{
			name:         "empty hostname fallback",
			hostname:     "",
			os:           "Linux",
			displayOS:    true,
			wantHostname: "hostinfo",
			wantLines:    1,
		},
		{
			name:         "whitespace hostname fallback",
			hostname:     "   ",
			os:           "Linux",
			displayOS:    true,
			wantHostname: "hostinfo",
			wantLines:    1,
		},
		{
			name:         "OS display disabled",
			hostname:     "myhost",
			os:           "Linux",
			displayOS:    false,
			wantHostname: "myhost",
			wantLines:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			banner, err := New(collectors.Providers{}, mustRenderer(t), nil)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			snap := collectors.Snapshot{
				System: collectors.SystemInfo{
					Hostname:  tt.hostname,
					OS:        tt.os,
					OSVersion: tt.osVersion,
					Arch:      tt.arch,
				},
			}

			cfg := config.Config{
				Display: config.DisplayConfig{
					OS: tt.displayOS,
				},
				ASCII: config.ASCIIConfig{
					Font:       "standard",
					Color:      "cyan",
					Monochrome: false,
				},
			}

			header := banner.buildHeader(snap, cfg)

			if header.Art == "" {
				t.Error("expected non-empty ASCII art")
			}

			if len(header.Lines) != tt.wantLines {
				t.Errorf("got %d header lines, want %d", len(header.Lines), tt.wantLines)
			}
		})
	}
}

func TestBanner_buildSections(t *testing.T) {
	tests := []struct {
		name         string
		builders     []Builder
		cfg          config.Config
		wantSections int
	}{
		{
			name:         "no builders",
			builders:     []Builder{},
			cfg:          config.Default(),
			wantSections: 0,
		},
		{
			name: "enabled builder",
			builders: []Builder{
				SystemSectionBuilder{},
			},
			cfg:          config.Default(),
			wantSections: 1,
		},
		{
			name: "disabled builder",
			builders: []Builder{
				SystemSectionBuilder{},
			},
			cfg: config.Config{
				Display: config.DisplayConfig{
					Uptime:   false,
					User:     false,
					Datetime: false,
				},
			},
			wantSections: 0,
		},
		{
			name: "nil builder in list",
			builders: []Builder{
				nil,
				SystemSectionBuilder{},
			},
			cfg:          config.Default(),
			wantSections: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			banner, err := New(collectors.Providers{}, mustRenderer(t), tt.builders)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			snap := collectors.Snapshot{
				System: collectors.SystemInfo{
					Uptime:   time.Hour * 24,
					Datetime: time.Now(),
				},
			}

			sections := banner.buildSections(snap, tt.cfg)

			if len(sections) != tt.wantSections {
				t.Errorf("got %d sections, want %d", len(sections), tt.wantSections)
			}
		})
	}
}

// mockSystemCollector implements collectors.SystemCollector for testing.
type mockSystemCollector struct {
	info collectors.SystemInfo
	err  error
}

func (m *mockSystemCollector) CollectSystem(ctx context.Context) (collectors.SystemInfo, error) {
	return m.info, m.err
}

// mustRenderer creates a renderer or fails the test.
func mustRenderer(t *testing.T) *ascii.Renderer {
	t.Helper()
	r, err := ascii.NewRenderer()
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}
	return r
}
