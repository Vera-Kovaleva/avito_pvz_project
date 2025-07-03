package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"avito_pvz/internal/domain"
	"avito_pvz/internal/infra/repository"
)

func TestPVZsIntegration(t *testing.T) {
	rollback(t, func(ctx context.Context, connection domain.Connection) {
		repoPvz := repository.NewPVZ()

		_ = fixtureCreatePVZ(t, ctx, connection, "Москва")

		pvzs, err := repoPvz.FindAll(ctx, connection)

		require.NoError(t, err)
		require.Equal(t, 1, len(pvzs))
	})
}

func fixtureCreatePVZ(t *testing.T, ctx context.Context, connection domain.Connection, city string) domain.PVZ {
	pvz := domain.PVZ{
		ID:           domain.PVZID(uuid.New()),
		City:         domain.PVZCity(city),
		RegisteredAt: time.Now(),
	}
	require.NoError(t, repository.NewPVZ().Create(ctx, connection, pvz))

	return pvz
}
