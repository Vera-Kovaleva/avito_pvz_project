package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"avito_pvz/internal/domain"
	"avito_pvz/internal/generated/mocks"
	"avito_pvz/internal/infra/repository"
)

func TestReceptionIntegration(t *testing.T) {
	rollback(t, func(ctx context.Context, connection domain.Connection) {

		repoReception := repository.NewReceptions()

		pvzID := uuid.New()
		receptionID := uuid.New()

		_ = fixtureCreatePVZ(ctx, t, connection, pvzID, "Казань")

		reception := fixtureCreateReceptin(ctx, t, connection, receptionID, pvzID)
		receptionFound, err := repoReception.FindActive(ctx, connection, pvzID)

		require.NoError(t, err)
		require.Equal(t, receptionFound.ID, reception.ID)
		require.Equal(t, receptionFound.PVZID, reception.PVZID)
		require.Equal(t, receptionFound.Status, domain.InProgress)

		_ = repoReception.Close(ctx, connection, receptionID)
		_, err = repoReception.FindActive(ctx, connection, pvzID)
		errors.Is(err, sql.ErrNoRows)

		receptionID2 := uuid.New()
		receptionID3 := uuid.New()

		reception2 := fixtureCreateReceptin(ctx, t, connection, receptionID2, pvzID)
		_ = repoReception.Close(ctx, connection, receptionID2)

		_ = fixtureCreateReceptin(ctx, t, connection, receptionID3, pvzID)
		receptionIDs := []domain.ReceptionID{receptionID2, receptionID3}

		receptionsFoundByIDs, err := repoReception.FindByIDs(ctx, connection, receptionIDs)

		require.NoError(t, err)
		require.Len(t, receptionsFoundByIDs, 2)
		require.Equal(t, receptionsFoundByIDs[0].PVZID, reception2.PVZID)
		require.Equal(t, receptionsFoundByIDs[0].Status, domain.Close)
		require.Equal(t, receptionsFoundByIDs[1].Status, domain.InProgress)

	})
}

func TestReceptionUnitCreate(t *testing.T) {
	connection := mocks.NewMockConnection(t)

	connection.EXPECT().ExecContext(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(0, errors.New("some error")).
		Once()

	err := repository.NewReceptions().Create(t.Context(), connection, domain.Reception{})

	require.ErrorIs(t, err, repository.ErrCreateReception)
	require.ErrorContains(t, err, "some error")
}

func TestReceptionUnitClose(t *testing.T) {
	connection := mocks.NewMockConnection(t)

	connection.EXPECT().ExecContext(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(0, errors.New("some error")).
		Once()

	var receptionID domain.ReceptionID
	err := repository.NewReceptions().Close(t.Context(), connection, receptionID)

	require.ErrorIs(t, err, repository.ErrCloseReception)
	require.ErrorContains(t, err, "some error")
}

func TestReceptionUnitFindByIDs(t *testing.T) {
	connection := mocks.NewMockConnection(t)

	connection.EXPECT().SelectContext(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("some error")).
		Once()

	var receptionIDs []domain.ReceptionID
	_, err := repository.NewReceptions().FindByIDs(t.Context(), connection, receptionIDs)

	require.ErrorIs(t, err, repository.ErrFindByIDsReception)
	require.ErrorContains(t, err, "some error")
}

func fixtureCreateReceptin(ctx context.Context, t *testing.T, connection domain.Connection, id domain.ReceptionID, pvzID domain.PVZID) domain.Reception {
	reception := domain.Reception{
		ID:     id,
		PVZID:  pvzID,
		Status: "in_progress",
	}
	require.NoError(t, repository.NewReceptions().Create(ctx, connection, reception))

	return reception
}
