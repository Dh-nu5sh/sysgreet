# Changelog

## v0.9.1

### Added

- **Gradient color support** - Banner lines can cycle through color gradients (default: brightblue → blue → cyan → brightcyan → white)
- **6 new Unicode block fonts** - ANSI Regular, ANSI Shadow, Block, Blocks, DOS Rebel, Basic for compact, professional banners
- **`--demo` flag** - Display 'SYSGREET' banner with realistic fake data, perfect for screenshots and demos
- **`--text` flag** - Render custom text as ASCII art (e.g., `--text "Production DB"`)
- **Visual padding** - Added blank line before banner output for better aesthetics

### Changed

- **Default font** - Changed from `slant` to `ANSI Regular` (compact Unicode blocks)
- **Default colors** - Now uses gradient instead of single random color
- **Banner style** - Unicode block characters (█) for tighter, more readable output
- **Bootstrap message** - Updated to reflect new defaults (ANSI Regular font with gradient)

### Documentation

- Updated README.md with hero image and demo screenshot
- Updated all example configs to show gradient configuration
- Documented all 8 available fonts in docs/examples/fonts.md
- Added special modes section to README (demo, text, disable)
- Updated quickstart guide with gradient and new flag examples

## v0.1.0

- Initial cross-platform sysgreet banner implementation
- Embedded FIGlet fonts with ASCII-art hostname rendering
- System, network, and resource collectors with graceful degradation
- YAML/TOML configuration support with optional monochrome mode
- GoReleaser pipeline and GitHub Actions release workflow
