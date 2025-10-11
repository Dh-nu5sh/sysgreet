package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/veteranbv/sysgreet/internal/ascii"
	"github.com/veteranbv/sysgreet/internal/banner"
	"github.com/veteranbv/sysgreet/internal/bootstrap"
	"github.com/veteranbv/sysgreet/internal/collectors"
	"github.com/veteranbv/sysgreet/internal/config"
	"github.com/veteranbv/sysgreet/internal/render"
)

func main() {
	ctx := context.Background()

	settings := parseFlags()
	policyEnv := os.Getenv("SYSGREET_CONFIG_POLICY")
	interactive := resolveInteractivity()
	if settings.Disable {
		return
	}

	renderer, err := ascii.NewRenderer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "sysgreet: %v\n", err)
		os.Exit(1)
	}

	if handled, err := runTextMode(renderer, settings.Text); err != nil {
		fmt.Fprintf(os.Stderr, "sysgreet: %v\n", err)
		os.Exit(1)
	} else if handled {
		return
	}

	if handled, err := runDemoMode(renderer, settings.Demo); err != nil {
		fmt.Fprintf(os.Stderr, "sysgreet: %v\n", err)
		os.Exit(1)
	} else if handled {
		return
	}

	// Normal mode: bootstrap config and collect real data
	cfgPath := config.DefaultWritePath()
	if err := maybeBootstrap(ctx, cfgPath, settings.PolicyFlag, policyEnv, interactive, bootstrap.IO{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}); err != nil {
		if errors.Is(err, bootstrap.ErrUserCanceled) {
			return
		}
		fmt.Fprintf(os.Stderr, "sysgreet: %v\n", err)
		os.Exit(1)
	}

	cfg, _, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "sysgreet: %v\n", err)
		os.Exit(1)
	}

	providers := collectors.Providers{
		System:    collectors.NewSystemCollector(),
		Network:   collectors.NewNetworkCollector(cfg.Network.MaxInterfaces),
		Resources: collectors.NewResourceCollector(),
		Session:   collectors.NewSessionCollector(),
		LastLogin: collectors.NewLastLoginCollector(),
	}

	hostBanner, err := banner.New(providers, renderer, banner.BuildersForConfig(cfg))
	if err != nil {
		fmt.Fprintf(os.Stderr, "sysgreet: %v\n", err)
		os.Exit(1)
	}

	output, _, err := hostBanner.Build(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sysgreet: %v\n", err)
		os.Exit(1)
	}

	layout := render.NewRenderer(cfg.ASCII.Monochrome)
	fmt.Println(layout.Render(output, cfg))
}

type runSettings struct {
	PolicyFlag string
	Disable    bool
	Demo       bool
	Text       string
}

func parseFlags() runSettings {
	policyFlag := flag.String("config-policy", "", "Config bootstrap policy: prompt, keep, or overwrite")
	disable := flag.Bool("disable", false, "Disable sysgreet output")
	demo := flag.Bool("demo", false, "Demo mode with 'SYSGREET' banner and fake data")
	text := flag.String("text", "", "Render custom text as ASCII art (e.g., --text \"Tea Pot\")")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), "\nEnvironment variables:")
		fmt.Fprintln(flag.CommandLine.Output(), "  SYSGREET_CONFIG_POLICY   Config bootstrap policy (prompt|keep|overwrite)")
		fmt.Fprintln(flag.CommandLine.Output(), "  SYSGREET_ASSUME_TTY      Force interactive prompts (testing/support)")
		fmt.Fprintln(flag.CommandLine.Output(), "  CI                      When set, disables interactive prompts by default")
		fmt.Fprintln(flag.CommandLine.Output(), "\nBootstrap:")
		fmt.Fprintln(flag.CommandLine.Output(), "  First run writes curated defaults (ANSI Regular font with gradient, metadata).")
		fmt.Fprintln(flag.CommandLine.Output(), "  Existing configs stay untouched unless you opt in via config policy.")
	}
	flag.Parse()
	return runSettings{
		PolicyFlag: *policyFlag,
		Disable:    *disable,
		Demo:       *demo,
		Text:       *text,
	}
}

func resolveInteractivity() bool {
	interactive := isInteractive()
	if os.Getenv("CI") != "" {
		interactive = false
	}
	if os.Getenv("SYSGREET_ASSUME_TTY") != "" {
		interactive = true
	}
	return interactive
}

func runTextMode(renderer *ascii.Renderer, text string) (bool, error) {
	if text == "" {
		return false, nil
	}
	cfg := config.Default()
	art, _, _, err := renderer.RenderWithGradient(text, cfg.ASCII.Font, cfg.ASCII.Color, cfg.ASCII.Gradient, cfg.ASCII.Monochrome)
	if err != nil {
		return false, err
	}
	fmt.Printf("\n%s\n\n", art)
	return true, nil
}

func runDemoMode(renderer *ascii.Renderer, enabled bool) (bool, error) {
	if !enabled {
		return false, nil
	}
	cfg := config.Default()
	demoSnap := collectors.DemoSnapshot()
	providers := collectors.Providers{} // Empty providers for demo
	hostBanner, err := banner.New(providers, renderer, banner.BuildersForConfig(cfg))
	if err != nil {
		return false, err
	}
	output := hostBanner.BuildWithSnapshot(demoSnap, cfg)
	layout := render.NewRenderer(cfg.ASCII.Monochrome)
	fmt.Println(layout.Render(output, cfg))
	return true, nil
}

func maybeBootstrap(ctx context.Context, cfgPath, policyFlag, policyEnv string, interactive bool, io bootstrap.IO) error {
	if cfgPath == "" {
		return nil
	}
	info, statErr := os.Stat(cfgPath)
	policyProvided := policyFlag != "" || policyEnv != ""
	configMissing := errors.Is(statErr, os.ErrNotExist)
	configIsDir := statErr == nil && info.IsDir()
	if statErr != nil && !configMissing {
		return fmt.Errorf("stat config: %w", statErr)
	}
	if !policyProvided && !configMissing && !configIsDir {
		return nil
	}
	_, err := bootstrap.Bootstrap(ctx, cfgPath, io, bootstrap.Options{FlagPolicy: policyFlag, EnvPolicy: policyEnv, Interactive: interactive})
	return err
}

func isInteractive() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
