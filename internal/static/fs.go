package static

import "embed"

//go:embed dist/*
var Dist embed.FS
