package database_test

import (
	"context"
	"os"
	"testing"
	"time"

	"avito_pvz/internal/domain"
	"avito_pvz/internal/infra/database"
	"avito_pvz/internal/infra/noerr"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestIntegrationPostgresProvider(t *testing.T) {
	t.Chdir("../../..")
	require.NoError(t, godotenv.Load())

	provider := database.NewPostgresProvider(
		noerr.Must(pgxpool.New(t.Context(), os.Getenv("DB_CONNECTION"))),
	)
	defer provider.Close()

	type nowRow struct {
		Now time.Time
	}

	err := provider.Execute(
		t.Context(),
		func(ctx context.Context, connection domain.Connection) error {
			var row nowRow
			return connection.GetContext(ctx, &row, "select now() as now")
		},
	)
	require.NoError(t, err)

	err = provider.Execute(
		t.Context(),
		func(ctx context.Context, connection domain.Connection) error {
			var row []nowRow
			return connection.SelectContext(ctx, &row, "select now() as now")
		},
	)
	require.NoError(t, err)

	err = provider.Execute(
		t.Context(),
		func(ctx context.Context, connection domain.Connection) error {
			_, err = connection.ExecContext(ctx, "select now() as now")

			return err
		},
	)
	require.NoError(t, err)

	err = provider.ExecuteTx(
		t.Context(),
		func(ctx context.Context, connection domain.Connection) error {
			var row nowRow
			return connection.GetContext(ctx, &row, "select now() as now")
		},
	)
	require.NoError(t, err)

	err = provider.ExecuteTx(
		t.Context(),
		func(ctx context.Context, connection domain.Connection) error {
			var row []nowRow
			return connection.SelectContext(ctx, &row, "select now() as now")
		},
	)
	require.NoError(t, err)

	err = provider.ExecuteTx(
		t.Context(),
		func(ctx context.Context, connection domain.Connection) error {
			_, err = connection.ExecContext(ctx, "select now() as now")

			return err
		},
	)
	require.NoError(t, err)
}
