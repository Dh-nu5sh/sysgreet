package ascii

import (
	"bytes"
	"fmt"
	"io/fs"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/veteranbv/sysgreet/assets"
)

var colorCodes = map[string]string{
	"reset":  "\033[0m",
	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"blue":   "\033[34m",
	"purple": "\033[35m",
	"cyan":   "\033[36m",
	"gray":   "\033[37m",
	"white":  "\033[97m",
}

var supportedColors = []string{"red", "green", "yellow", "blue", "purple", "cyan", "gray", "white"}

// Renderer produces ASCII art banners from embedded FIGlet fonts.
type Renderer struct {
	fonts map[string][]byte
	order []string
	rnd   *rand.Rand
}

// NewRenderer loads embedded fonts into memory.
func NewRenderer() (*Renderer, error) {
	paths, err := fs.Glob(assets.FontsFS, "fonts/*.flf")
	if err != nil {
		return nil, err
	}
	fonts := make(map[string][]byte)
	var order []string
	for _, p := range paths {
		data, err := assets.FontsFS.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("load font %s: %w", p, err)
		}
		name := strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))
		fonts[name] = data
		order = append(order, name)
	}
	if len(order) == 0 {
		return nil, fmt.Errorf("no fonts embedded")
	}
	return &Renderer{
		fonts: fonts,
		order: order,
		rnd:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

// Fonts returns the list of available fonts.
func (r *Renderer) Fonts() []string {
	return append([]string{}, r.order...)
}

// Render creates ASCII art for the input string using the specified font and color.
// If fontName or colorName equal "random", selections are randomized from the embedded sets.
func (r *Renderer) Render(text, fontName, colorName string, monochrome bool) (string, string, string, error) {
	if fontName == "" || fontName == "random" {
		fontName = r.randomFont()
	} else if _, ok := r.fonts[fontName]; !ok {
		fontName = r.randomFont()
	}

	fontData := r.fonts[fontName]
	fig := figure.NewFigureWithFont(text, bytes.NewReader(fontData), true)
	rows := fig.Slicify()
	asciiArt := strings.Join(rows, "\n")

	color := "reset"
	if !monochrome {
		color = r.pickColor(colorName)
		if code, ok := colorCodes[color]; ok && color != "reset" {
			asciiArt = fmt.Sprintf("%s%s%s", code, asciiArt, colorCodes["reset"])
		}
	}

	return asciiArt, fontName, color, nil
}

func (r *Renderer) randomFont() string {
	return r.order[r.rnd.Intn(len(r.order))]
}

func (r *Renderer) pickColor(name string) string {
	if name == "" || name == "random" {
		return supportedColors[r.rnd.Intn(len(supportedColors))]
	}
	name = strings.ToLower(name)
	if _, ok := colorCodes[name]; ok {
		return name
	}
	return "reset"
}
