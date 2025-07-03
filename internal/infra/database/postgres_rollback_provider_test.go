package database_test

import (
	"context"
	"errors"
	"testing"

	"avito_pvz/internal/domain"
	"avito_pvz/internal/generated/mocks"
	"avito_pvz/internal/infra/database"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPosgresRollbackProviderAlwaysRollsback(t *testing.T) {
	t.Parallel()

	underlyingProvider := mocks.NewMockConnectionProvider(t)
	underlyingConnection := mocks.NewMockConnection(t)

	rollbackProvider := database.NewPostgresRollbackProvider(underlyingProvider)

	underlyingProvider.EXPECT().
		ExecuteTx(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, receiver func(context.Context, domain.Connection) error) error {
			err := receiver(ctx, underlyingConnection)
			require.Error(t, err)
			require.ErrorContains(t, err, "this provider always rollbacks transactions")

			return err
		}).
		Once()
	rollbackProvider.Execute(
		t.Context(),
		func(_ context.Context, _ domain.Connection) error { return nil },
	)

	underlyingProvider.EXPECT().Close().Return(nil).Once()
	require.NoError(t, rollbackProvider.Close())
}

func TestPosgresRollbackProviderTransparentlyProxiesErrors(t *testing.T) {
	t.Parallel()

	underlyingProvider := mocks.NewMockConnectionProvider(t)
	underlyingConnection := mocks.NewMockConnection(t)
	someErr := errors.New("some err")

	rollbackProvider := database.NewPostgresRollbackProvider(underlyingProvider)

	underlyingProvider.EXPECT().
		ExecuteTx(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, receiver func(context.Context, domain.Connection) error) error {
			err := receiver(ctx, underlyingConnection)
			require.Error(t, err)
			require.ErrorIs(t, err, someErr)
			require.NotContains(t, err.Error(), "this provider always rollbacks transactions")

			return err
		}).
		Once()
	rollbackProvider.Execute(
		t.Context(),
		func(_ context.Context, _ domain.Connection) error { return someErr },
	)

	underlyingProvider.EXPECT().Close().Return(nil).Once()
	require.NoError(t, rollbackProvider.Close())
}
