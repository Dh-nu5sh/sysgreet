# Embedded FIGlet Fonts

The Hostinfo banner embeds two FIGlet fonts to ensure the CLI operates offline:

- `standard.flf`
- `slant.flf`

Both files originate from the FIGlet project (<http://www.figlet.org/>). They are licensed under the FIGlet Font License, which permits redistribution and embedding in binaries. License text is included at the top of each font file.

Fonts are stored in `assets/fonts/` and shipped with the binary via Go's `//go:embed` directive.
