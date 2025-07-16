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

func TestUserIntegration(t *testing.T) {
	rollback(t, func(ctx context.Context, connection domain.Connection) {
		pvzID := uuid.New()
		receptionID := uuid.New()
		productID := uuid.New()
		userID := uuid.New()

		_ = fixtureCreatePVZ(ctx, t, connection, pvzID, "Казань")
		_ = fixtureCreateReceptin(ctx, t, connection, receptionID, pvzID)
		_ = fixtureCreateProduct(
			ctx,
			t,
			connection,
			productID,
			receptionID,
			"электроника",
			time.Now(),
		)

		users := repository.NewUsers()

		user := fixtureCreateUser(ctx, t, connection, userID, "employee")

		userRead, err := users.ReadByEmail(ctx, connection, user.Email)
		require.NoError(t, err)
		require.Equal(t, user, userRead)

		newUser := domain.User{
			ID:           userID,
			Role:         "moderator",
			Email:        "newUser@email.foo",
			PasswordHash: "some new password hash",
			Token:        "some new secret token",
		}

		err = users.Update(ctx, connection, newUser)
		require.NoError(t, err)
		userRead, err = users.ReadByEmail(ctx, connection, newUser.Email)
		require.NoError(t, err)
		require.Equal(t, newUser, userRead)

		err = users.UpdateTokenByEmail(ctx, connection, newUser.Email, "some new new token")
		require.NoError(t, err)
		userRead, err = users.ReadByEmail(ctx, connection, newUser.Email)
		require.NoError(t, err)
		require.Equal(t, userRead.Token, "some new new token")
	})
}

func TestUserUnitCreate(t *testing.T) {
	connection := mocks.NewMockConnection(t)

	connection.EXPECT().
		ExecContext(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(0, errors.New("some error")).
		Once()

	err := repository.NewUsers().Create(t.Context(), connection, domain.User{})
	require.ErrorIs(t, err, repository.ErrUsersCreate)
	require.ErrorContains(t, err, "some error")
}

func TestUserUnitReadByEmail(t *testing.T) {
	connection := mocks.NewMockConnection(t)

	connection.EXPECT().GetContext(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("some error")).
		Once()

	_, err := repository.NewUsers().ReadByEmail(t.Context(), connection, "")
	require.ErrorIs(t, err, repository.ErrUsersReadByEmail)
	require.ErrorContains(t, err, "some error")
}

func TestUserUnitUpdate(t *testing.T) {
	connection := mocks.NewMockConnection(t)

	connection.EXPECT().
		ExecContext(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(0, errors.New("some error")).
		Once()

	err := repository.NewUsers().Update(t.Context(), connection, domain.User{})
	require.ErrorIs(t, err, repository.ErrUsersUpdate)
	require.ErrorContains(t, err, "some error")
}

func TestUserUnitUpdateTokenByEmail(t *testing.T) {
	connection := mocks.NewMockConnection(t)

	connection.EXPECT().
		ExecContext(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(0, errors.New("some error")).
		Once()

	err := repository.NewUsers().UpdateTokenByEmail(t.Context(), connection, "", "")
	require.ErrorIs(t, err, repository.ErrUsersUpdateTokenByEmail)
	require.ErrorContains(t, err, "some error")
}

func fixtureCreateUser(
	ctx context.Context,
	t *testing.T,
	connection domain.Connection,
	id domain.UserID,
	role domain.UserRole,
) domain.User {
	user := domain.User{
		ID:           id,
		Role:         role,
		Email:        "user@email.foo",
		PasswordHash: "some password hash",
		Token:        "some secret token",
	}
	require.NoError(t, repository.NewUsers().Create(ctx, connection, user))

	return user
}
