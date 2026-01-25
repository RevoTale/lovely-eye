//go:build tools
// +build tools

package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-bun/bunschema"
	_ "ariga.io/atlas/sdk/recordriver"

	"github.com/lovely-eye/server/internal/models"
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
		&models.User{},
		&models.Site{},
		&models.SiteDomain{},
		&models.SiteBlockedIP{},
		&models.SiteBlockedCountry{},
		&models.Client{},
		&models.Session{},
		&models.Event{},
		&models.EventDefinition{},
		&models.EventDefinitionField{},
		&models.EventData{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
