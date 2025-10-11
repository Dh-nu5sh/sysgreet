package assets

import "embed"

// FontsFS provides access to embedded FIGlet fonts.
//
//go:embed fonts/*.flf
var FontsFS embed.FS
