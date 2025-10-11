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
		fmt.Fprintln(flag.CommandLine.Output(), "  Existing configs prompt to keep or overwrite unless a policy is supplied.")
	}
	flag.Parse()
	policyEnv := os.Getenv("SYSGREET_CONFIG_POLICY")
	interactive := isInteractive()
	if os.Getenv("CI") != "" {
		interactive = false
	}
	if os.Getenv("SYSGREET_ASSUME_TTY") != "" {
		interactive = true
	}
	if *disable {
		return
	}

	renderer, err := ascii.NewRenderer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "sysgreet: %v\n", err)
		os.Exit(1)
	}

	// Text mode: render custom ASCII art only
	if *text != "" {
		cfg := config.Default()
		art, _, _, err := renderer.RenderWithGradient(*text, cfg.ASCII.Font, cfg.ASCII.Color, cfg.ASCII.Gradient, cfg.ASCII.Monochrome)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sysgreet: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\n%s\n\n", art)
		return
	}

	// Demo mode: use fake data and default config
	if *demo {
		cfg := config.Default()
		demoSnap := collectors.DemoSnapshot()
		providers := collectors.Providers{} // Empty providers for demo
		hostBanner, err := banner.New(providers, renderer, banner.BuildersForConfig(cfg))
		if err != nil {
			fmt.Fprintf(os.Stderr, "sysgreet: %v\n", err)
			os.Exit(1)
		}
		output := hostBanner.BuildWithSnapshot(demoSnap, cfg)
		layout := render.NewRenderer(cfg.ASCII.Monochrome)
		fmt.Println(layout.Render(output, cfg))
		return
	}

	// Normal mode: bootstrap config and collect real data
	cfgPath := config.DefaultWritePath()
	if cfgPath != "" {
		info, statErr := os.Stat(cfgPath)
		policyProvided := *policyFlag != "" || policyEnv != ""
		configMissing := errors.Is(statErr, os.ErrNotExist)
		configIsDir := statErr == nil && info.IsDir()
		if statErr != nil && !configMissing {
			fmt.Fprintf(os.Stderr, "sysgreet: stat config: %v\n", statErr)
			os.Exit(1)
		}
		if policyProvided || configMissing || configIsDir {
			if _, err := bootstrap.Bootstrap(ctx, cfgPath, bootstrap.IO{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}, bootstrap.Options{FlagPolicy: *policyFlag, EnvPolicy: policyEnv, Interactive: interactive}); err != nil {
				if errors.Is(err, bootstrap.ErrUserCanceled) {
					return
				}
				fmt.Fprintf(os.Stderr, "sysgreet: %v\n", err)
				os.Exit(1)
			}
		}
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

func isInteractive() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
