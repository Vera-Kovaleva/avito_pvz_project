package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"avito_pvz/internal/domain"
	"avito_pvz/internal/generated/mocks"
	"avito_pvz/internal/infra/repository"
)

func TestPVZsIntegration(t *testing.T) {
	rollback(t, func(ctx context.Context, connection domain.Connection) {
		repoPvz := repository.NewPVZ()
		uuid1 := uuid.New()
		uuid2 := uuid.New()

		pvz1 := fixtureCreatePVZ(ctx, t, connection, uuid1, "Москва")
		pvz2 := fixtureCreatePVZ(ctx, t, connection, uuid2, "Казань")

		pvzAll, err := repoPvz.FindAll(ctx, connection)

		require.NoError(t, err)
		require.Len(t, pvzAll, 2)

		var pvzIDs []domain.PVZID
		pvzIDs = append(pvzIDs, uuid1, uuid2)

		pvzFound, err := repoPvz.FindByIDs(ctx, connection, pvzIDs)

		var pvzDefault []domain.PVZ
		pvzDefault = append(pvzDefault, pvz1, pvz2)

		require.NoError(t, err)
		require.Equal(t, pvzFound[0].ID, pvzDefault[0].ID)
		require.Equal(t, pvzFound[0].City, pvzDefault[0].City)
		require.Equal(t, pvzFound[1].ID, pvzDefault[1].ID)
		require.Equal(t, pvzFound[1].City, pvzDefault[1].City)
	})
}

func TestPVZUnitCreate(t *testing.T) {
	pvz := repository.NewPVZ()
	connection := mocks.NewMockConnection(t)

	connection.EXPECT().
		ExecContext(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(0, errors.New("some error")).
		Once()

	err := pvz.Create(t.Context(), connection, domain.PVZ{})

	require.ErrorIs(t, err, repository.ErrPVZCreate)
	require.ErrorContains(t, err, "some error")
}

func TestPVZUnitFindAll(t *testing.T) {
	connection := mocks.NewMockConnection(t)

	connection.EXPECT().SelectContext(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("some error")).
		Once()

	_, err := repository.NewPVZ().FindAll(t.Context(), connection)

	require.ErrorIs(t, err, repository.ErrPVZFindAll)
	require.ErrorContains(t, err, "some error")
}

func TestPVZUnitFindByIDsl(t *testing.T) {
	connection := mocks.NewMockConnection(t)

	connection.EXPECT().SelectContext(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("some error")).
		Once()

	var pvzIDs []domain.PVZID
	_, err := repository.NewPVZ().FindByIDs(t.Context(), connection, pvzIDs)

	require.ErrorIs(t, err, repository.ErrPVZFindByIDs)
	require.ErrorContains(t, err, "some error")
}

func fixtureCreatePVZ(
	ctx context.Context,
	t *testing.T,
	connection domain.Connection,
	id domain.PVZID,
	city string,
) domain.PVZ {
	pvz := domain.PVZ{
		ID:           id,
		City:         domain.PVZCity(city),
		RegisteredAt: time.Now(),
	}
	require.NoError(t, repository.NewPVZ().Create(ctx, connection, pvz))

	return pvz
}
