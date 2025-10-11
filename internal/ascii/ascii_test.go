package ascii

import "testing"

func TestRendererRenderSpecificFont(t *testing.T) {
	r, err := NewRenderer()
	if err != nil {
		t.Fatalf("NewRenderer() error = %v", err)
	}
	art, font, _, err := r.Render("host", "standard", "reset", true)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if font != "standard" {
		t.Fatalf("expected font 'standard', got %s", font)
	}
	if len(art) == 0 {
		t.Fatalf("expected non-empty art output")
	}
	if art == "host" {
		t.Fatalf("expected ASCII art, got plain text")
	}
}

func TestRendererRandomFontSelection(t *testing.T) {
	r, err := NewRenderer()
	if err != nil {
		t.Fatalf("NewRenderer() error = %v", err)
	}
	_, font, color, err := r.Render("host", "random", "random", false)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	fonts := r.Fonts()
	found := false
	for _, candidate := range fonts {
		if font == candidate {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("random font %s not in available list", font)
	}
	if color == "reset" {
		t.Fatalf("expected random color selection")
	}
}
