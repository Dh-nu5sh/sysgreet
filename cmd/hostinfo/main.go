package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/veteranbv/hostinfo/internal/ascii"
	"github.com/veteranbv/hostinfo/internal/banner"
	"github.com/veteranbv/hostinfo/internal/collectors"
	"github.com/veteranbv/hostinfo/internal/config"
	"github.com/veteranbv/hostinfo/internal/render"
)

func main() {
	ctx := context.Background()

	disable := flag.Bool("disable", false, "Disable hostinfo output")
	flag.Parse()
	if *disable {
		return
	}

	cfg, _, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "hostinfo: %v\n", err)
		os.Exit(1)
	}

	renderer, err := ascii.NewRenderer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "hostinfo: %v\n", err)
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
		fmt.Fprintf(os.Stderr, "hostinfo: %v\n", err)
		os.Exit(1)
	}

	output, _, err := hostBanner.Build(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "hostinfo: %v\n", err)
		os.Exit(1)
	}

	layout := render.NewRenderer(cfg.ASCII.Monochrome)
	fmt.Println(layout.Render(output, cfg))
}
