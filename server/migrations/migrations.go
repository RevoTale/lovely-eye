package migrations

import (
	"embed"
	"fmt"
	"io/fs"
	"os"

	"github.com/uptrace/bun/migrate"
)

//go:embed sqlite/*.sql postgres/*.sql
var sqlMigrations embed.FS

func NewMigrations() (*migrate.Migrations, error) {
	// Determine database driver from environment
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "sqlite"
	}

	migrations := migrate.NewMigrations()

	var migrationFS fs.FS
	var err error

	// Load migrations for the appropriate database
	switch driver {
	case "postgres":
		migrationFS, err = fs.Sub(sqlMigrations, "postgres")
	case "sqlite":
		migrationFS, err = fs.Sub(sqlMigrations, "sqlite")
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	if err != nil {
		return nil, err
	}

	if err := migrations.Discover(migrationFS); err != nil {
		return nil, err
	}

	return migrations, nil
}
