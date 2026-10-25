// Package assets embeds the static files served under /assets/, so the binary
// is the only thing an air-gapped cluster needs to run.
package assets

import "embed"

//go:embed all:css all:fonts all:img all:js
var Assets embed.FS
