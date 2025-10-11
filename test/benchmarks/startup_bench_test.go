package benchmarks

import (
	"context"
	"testing"

	"github.com/veteranbv/hostinfo/internal/ascii"
	"github.com/veteranbv/hostinfo/internal/banner"
	"github.com/veteranbv/hostinfo/internal/collectors"
	"github.com/veteranbv/hostinfo/internal/config"
)

// BenchmarkStartup is a high-level benchmark that builds a banner using stub collectors.
func BenchmarkStartup(b *testing.B) {
	renderer, err := ascii.NewRenderer()
	if err != nil {
		b.Fatalf("failed to create renderer: %v", err)
	}
	providers := collectors.Providers{}
	nBanner, err := banner.New(providers, renderer, nil)
	if err != nil {
		b.Fatalf("failed to create banner: %v", err)
	}
	cfg := config.Default()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, _, err := nBanner.Build(ctx, cfg); err != nil {
			b.Fatalf("banner build failed: %v", err)
		}
	}
}
