# Embedded FIGlet Fonts

The Sysgreet banner embeds multiple FIGlet fonts to ensure the CLI operates offline:

## Available Fonts

- **`ANSI Regular.flf`** (default) - Compact Unicode block characters with solid fills
- **`ANSI Shadow.flf`** - Unicode blocks with shadow effects
- **`Block.flf`** - Classic blocky style
- **`Blocks.flf`** - Variant of block style
- **`DOS Rebel.flf`** - Modern, clean DOS-style font
- **`Basic.flf`** - Simple hash-mark style
- **`standard.flf`** - Classic FIGlet standard font
- **`slant.flf`** - Tall, angular slanted font

All fonts originate from the FIGlet project (<http://www.figlet.org/>) and community repositories. They are licensed under the FIGlet Font License, which permits redistribution and embedding in binaries. License text is included at the top of each font file.

Fonts are stored in `assets/fonts/` and shipped with the binary via Go's `//go:embed` directive.

## Usage

Fonts can be configured in your config file:

```yaml
ascii:
  font: "ANSI Regular"  # Default - compact Unicode blocks
  # or "ANSI Shadow", "Block", "DOS Rebel", "Basic", "standard", "slant"
```

Or test fonts quickly:

```bash
sysgreet --text "Test" --config-policy=keep
```
