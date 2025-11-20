package ui

import (
	"embed"
)

// special comment directive
//
//go:embed "html" "static"
var Files embed.FS
