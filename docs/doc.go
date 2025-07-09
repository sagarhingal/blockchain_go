package docs

import "embed"

// FS contains project documentation and swagger specification.
//
//go:embed *.md swagger/*
var FS embed.FS
