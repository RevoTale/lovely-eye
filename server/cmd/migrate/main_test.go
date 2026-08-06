package main

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"

	"github.com/lovely-eye/server/internal/platform/config"
	"github.com/lovely-eye/server/internal/platform/database"
	"github.com/lovely-eye/server/migrations"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun/migrate"
)

func TestParseCommand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		want    command
		wantErr string
	}{
		{name: "no arguments shows help", want: commandHelp},
		{name: "help flag", args: []string{"--help"}, want: commandHelp},
		{name: "up", args: []string{"up"}, want: commandUp},
		{name: "unknown", args: []string{"unknown"}, wantErr: "unknown command"},
		{name: "extra argument", args: []string{"up", "extra"}, wantErr: "exactly one command"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseCommand(test.args)
			if test.wantErr == "" {
				require.NoError(t, err)
				require.Equal(t, test.want, got)
				return
			}
			require.ErrorContains(t, err, test.wantErr)
		})
	}
}

func TestExecuteCommandLifecycle(t *testing.T) {
	db, err := database.New(context.Background(), config.DatabaseConfig{
		Driver: config.DBDriverSQLite,
		DSN:    filepath.Join(t.TempDir(), "migrations.db"),
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(db)) })

	migrationSet, err := migrations.NewMigrations()
	require.NoError(t, err)
	migrator := migrate.NewMigrator(db, migrationSet)
	ctx := context.Background()
	var output bytes.Buffer

	require.NoError(t, executeCommand(ctx, commandInit, &output, migrator))
	require.NoError(t, executeCommand(ctx, commandUp, &output, migrator))
	require.Contains(t, output.String(), "migrated to")

	output.Reset()
	require.NoError(t, executeCommand(ctx, commandStatus, &output, migrator))
	require.Contains(t, output.String(), "unapplied migrations: empty")

	output.Reset()
	require.NoError(t, executeCommand(ctx, commandDown, &output, migrator))
	require.Contains(t, output.String(), "rolled back")
}

func TestExecuteCommandRejectsUnsupportedValue(t *testing.T) {
	t.Parallel()

	err := executeCommand(context.Background(), "invalid", &bytes.Buffer{}, nil)
	require.ErrorContains(t, err, "unsupported migration command")
}
