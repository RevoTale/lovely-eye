//go:build tools
// +build tools

package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-bun/bunschema"
	_ "ariga.io/atlas/sdk/recordriver"

	analyticspersistence "github.com/lovely-eye/server/internal/analytics/persistence"
	authpersistence "github.com/lovely-eye/server/internal/auth/persistence"
	countrypersistence "github.com/lovely-eye/server/internal/country/persistence"
	eventpersistence "github.com/lovely-eye/server/internal/event/persistence"
	sitepersistence "github.com/lovely-eye/server/internal/site/persistence"
)

func main() {
	dialect := os.Getenv("ATLAS_DIALECT")
	if dialect == "" {
		dialect = "sqlite"
	}

	var d bunschema.Dialect
	switch dialect {
	case "postgres":
		d = bunschema.DialectPostgres
	case "sqlite":
		d = bunschema.DialectSQLite
	default:
		fmt.Fprintf(os.Stderr, "unsupported dialect: %s\n", dialect)
		os.Exit(1)
	}

	stmts, err := bunschema.New(d).Load(
		&authpersistence.User{},
		&sitepersistence.Site{},
		&sitepersistence.Domain{},
		&sitepersistence.BlockedIP{},
		&sitepersistence.BlockedCountry{},
		&countrypersistence.Country{},
		&analyticspersistence.Client{},
		&analyticspersistence.Session{},
		&analyticspersistence.Event{},
		&eventpersistence.Definition{},
		&eventpersistence.Field{},
		&analyticspersistence.EventData{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
