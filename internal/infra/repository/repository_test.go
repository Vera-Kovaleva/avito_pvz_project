package repository_test

import (
	"context"
	"os"
	"testing"

	"avito_pvz/internal/domain"
	"avito_pvz/internal/infra/database"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func rollback(t *testing.T, receiver func(context.Context, domain.Connection)) {
	require.NoError(t, godotenv.Load("../../../.env"))

	pool, err := pgxpool.New(t.Context(), os.Getenv("DB_CONNECTION"))
	require.NoError(t, err)

	provider := database.NewPostgresRollbackProvider(database.NewPostgresProvider(pool))

	clearTable := func(t *testing.T, connection domain.Connection, table string) {
		_, errTruncate := connection.ExecContext(t.Context(), "delete from "+table+" cascade")
		require.NoError(t, errTruncate)
	}

	require.NoError(
		t,
		provider.ExecuteTx(
			t.Context(),
			func(ctx context.Context, connection domain.Connection) error {
				clearTable(t, connection, "products")
				clearTable(t, connection, "receptions")
				clearTable(t, connection, "pvz")
				clearTable(t, connection, "users")

				receiver(ctx, connection)

				return nil
			},
		),
	)
}
