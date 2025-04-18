package static

import "embed"

//go:embed public/*
var Dist embed.FS
