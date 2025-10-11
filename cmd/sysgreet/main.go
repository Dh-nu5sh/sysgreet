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
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), "\nEnvironment variables:")
		fmt.Fprintln(flag.CommandLine.Output(), "  SYSGREET_CONFIG_POLICY   Config bootstrap policy (prompt|keep|overwrite)")
		fmt.Fprintln(flag.CommandLine.Output(), "  SYSGREET_ASSUME_TTY      Force interactive prompts (testing/support)")
		fmt.Fprintln(flag.CommandLine.Output(), "  CI                      When set, disables interactive prompts by default")
		fmt.Fprintln(flag.CommandLine.Output(), "\nBootstrap:")
		fmt.Fprintln(flag.CommandLine.Output(), "  First run writes curated defaults (slant font, metadata).")
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

	cfgPath := config.DefaultWritePath()
	if cfgPath != "" {
		if _, err := bootstrap.Bootstrap(ctx, cfgPath, bootstrap.IO{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}, bootstrap.Options{FlagPolicy: *policyFlag, EnvPolicy: policyEnv, Interactive: interactive}); err != nil {
			if errors.Is(err, bootstrap.ErrUserCancelled) {
				return
			}
			fmt.Fprintf(os.Stderr, "sysgreet: %v\n", err)
			os.Exit(1)
		}
	}

	cfg, _, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "sysgreet: %v\n", err)
		os.Exit(1)
	}

	renderer, err := ascii.NewRenderer()
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
