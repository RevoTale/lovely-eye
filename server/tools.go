//go:build tools

//go:generate go run github.com/99designs/gqlgen generate

package tools

// This file declares tool dependencies for go generate.
// Run `go mod tidy` to ensure these are in go.mod.

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/Khan/genqlient"
)
