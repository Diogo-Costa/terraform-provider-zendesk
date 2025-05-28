//go:build tools
// +build tools

package tools

// Development tool dependencies
// These are not imported by the main package but are used for development
import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/katbyte/terrafmt"
)