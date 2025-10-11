package bootstrap

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"testing"
)

func BenchmarkBootstrapFirstRun(b *testing.B) {
	ctx := context.Background()
	dir := b.TempDir()
	for i := 0; i < b.N; i++ {
		cfgPath := filepath.Join(dir, fmt.Sprintf("config-%d.yaml", i))
		b.StartTimer()
		_, err := Bootstrap(ctx, cfgPath, IO{Stderr: io.Discard}, Options{Interactive: true})
		b.StopTimer()
		if err != nil {
			b.Fatalf("bootstrap error: %v", err)
		}
	}
}
