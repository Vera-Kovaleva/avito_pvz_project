package domain_test

import (
	"avito_pvz/internal/domain"
	"avito_pvz/internal/generated/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServiceReception_Create(t *testing.T) {
	t.Parallel()

	pvzID := domain.PVZID(uuid.New())
	reception := domain.Reception{
		ID:    domain.ReceptionID(uuid.New()),
		PVZID: pvzID,
	}
	var invalidPVZID = domain.PVZID(uuid.Nil)

	tests := []struct {
		name         string
		authUser     domain.AuthenticatedUser
		pvzID        domain.PVZID
		prepareMocks func(*mocks.MockConnectionProvider, *mocks.MockReceptionsRepository)
		check        func(*testing.T, domain.Reception, error)
	}{
		{
			name:     "Success",
			authUser: nil,
			pvzID:    pvzID,
			prepareMocks: func(provider *mocks.MockConnectionProvider, repo *mocks.MockReceptionsRepository) {
				provider.EXPECT().
					ExecuteTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				provider.EXPECT().
					Execute(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				repo.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).
					Return(nil).Once()
				repo.EXPECT().FindActive(mock.Anything, mock.Anything, pvzID).
					Return(reception, nil).Once()
			},
			check: func(t *testing.T, reception domain.Reception, err error) {
				require.NoError(t, err)
				require.Equal(t, pvzID, reception.PVZID)
			},
		},
		{
			name:     "DB Error",
			authUser: nil,
			pvzID:    pvzID,
			prepareMocks: func(provider *mocks.MockConnectionProvider, repo *mocks.MockReceptionsRepository) {
				provider.EXPECT().
					ExecuteTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				repo.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("some error")).Once()
			},
			check: func(t *testing.T, reception domain.Reception, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "some error")
				require.Contains(t, err.Error(), "create failed")
			},
		},
		{
			name:     "Error find active",
			authUser: nil,
			pvzID:    pvzID,
			prepareMocks: func(provider *mocks.MockConnectionProvider, repo *mocks.MockReceptionsRepository) {
				provider.EXPECT().
					ExecuteTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				provider.EXPECT().
					Execute(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				repo.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).
					Return(nil).Once()
				repo.EXPECT().FindActive(mock.Anything, mock.Anything, pvzID).
					Return(reception, errors.New("some error")).Once()
			},
			check: func(t *testing.T, reception domain.Reception, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "some error")
				require.Contains(t, err.Error(), "find active failed")
			},
		},
		{
			name:     "Invalid ID",
			authUser: nil,
			pvzID:    invalidPVZID,
			prepareMocks: func(provider *mocks.MockConnectionProvider, repo *mocks.MockReceptionsRepository) {
			},
			check: func(t *testing.T, r domain.Reception, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "invalid pvz id")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			provider := mocks.NewMockConnectionProvider(t)

			repoReception := mocks.NewMockReceptionsRepository(t)
			repoProduct := mocks.NewMockProductsRepository(t)

			if test.prepareMocks != nil {
				test.prepareMocks(provider, repoReception)
			}

			reception, err := domain.NewReceptionService(provider, repoReception, repoProduct).Create(t.Context(), test.authUser, test.pvzID)

			test.check(t, reception, err)
		})
	}
}

func TestServiceReception_Close(t *testing.T) {
	t.Parallel()

	pvzID := domain.PVZID(uuid.New())
	reception := domain.Reception{
		ID:    domain.ReceptionID(uuid.New()),
		PVZID: pvzID,
	}
	var invalidPVZID = domain.PVZID(uuid.Nil)

	tests := []struct {
		name         string
		authUser     domain.AuthenticatedUser
		pvzID        domain.PVZID
		prepareMocks func(*mocks.MockConnectionProvider, *mocks.MockReceptionsRepository)
		check        func(*testing.T, domain.Reception, error)
	}{
		{
			name:     "Success",
			authUser: nil,
			pvzID:    pvzID,
			prepareMocks: func(provider *mocks.MockConnectionProvider, repo *mocks.MockReceptionsRepository) {
				provider.EXPECT().
					Execute(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				provider.EXPECT().
					ExecuteTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				repo.EXPECT().FindActive(mock.Anything, mock.Anything, mock.Anything).
					Return(reception, nil).Once()
				repo.EXPECT().Close(mock.Anything, mock.Anything, reception.ID).
					Return(nil).Once()
			},
			check: func(t *testing.T, reception domain.Reception, err error) {
				require.NoError(t, err)
				require.Equal(t, pvzID, reception.PVZID)
			},
		},
		{
			name:     "Find active Error",
			authUser: nil,
			pvzID:    pvzID,
			prepareMocks: func(provider *mocks.MockConnectionProvider, repo *mocks.MockReceptionsRepository) {
				provider.EXPECT().
					Execute(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				repo.EXPECT().FindActive(mock.Anything, mock.Anything, mock.Anything).
					Return(reception, errors.New("some error")).Once()
			},
			check: func(t *testing.T, reception domain.Reception, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "some error")
				require.Contains(t, err.Error(), "find active failed")
			},
		},
		{
			name:     "Close error",
			authUser: nil,
			pvzID:    pvzID,
			prepareMocks: func(provider *mocks.MockConnectionProvider, repo *mocks.MockReceptionsRepository) {
				provider.EXPECT().
					Execute(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				provider.EXPECT().
					ExecuteTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				repo.EXPECT().FindActive(mock.Anything, mock.Anything, mock.Anything).
					Return(reception, nil).Once()
				repo.EXPECT().Close(mock.Anything, mock.Anything, reception.ID).
					Return(errors.New("some error")).Once()
			},
			check: func(t *testing.T, reception domain.Reception, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "some error")
				require.Contains(t, err.Error(), "close failed")
			},
		},
		{
			name:     "Invalid ID",
			authUser: nil,
			pvzID:    invalidPVZID,
			prepareMocks: func(provider *mocks.MockConnectionProvider, repo *mocks.MockReceptionsRepository) {
			},
			check: func(t *testing.T, r domain.Reception, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "invalid pvz id")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			provider := mocks.NewMockConnectionProvider(t)

			repoReception := mocks.NewMockReceptionsRepository(t)
			repoProduct := mocks.NewMockProductsRepository(t)

			if test.prepareMocks != nil {
				test.prepareMocks(provider, repoReception)
			}

			reception, err := domain.NewReceptionService(provider, repoReception, repoProduct).Close(t.Context(), test.authUser, test.pvzID)

			test.check(t, reception, err)
		})
	}
}

func TestServiceReception_CreateProduct(t *testing.T) {
	t.Parallel()

	pvzID := domain.PVZID(uuid.New())
	reception := domain.Reception{
		ID:     domain.ReceptionID(uuid.New()),
		PVZID:  pvzID,
		Status: domain.InProgress,
	}

	var invalidPVZID = domain.PVZID(uuid.Nil)

	tests := []struct {
		name         string
		authUser     domain.AuthenticatedUser
		pvzID        domain.PVZID
		productType  domain.ProductType
		prepareMocks func(*mocks.MockConnectionProvider, *mocks.MockReceptionsRepository, *mocks.MockProductsRepository)
		check        func(*testing.T, domain.Product, error)
	}{
		{
			name:     "Success",
			authUser: nil,
			pvzID:    pvzID,
			prepareMocks: func(provider *mocks.MockConnectionProvider, repoReception *mocks.MockReceptionsRepository, repoProduct *mocks.MockProductsRepository) {
				provider.EXPECT().
					Execute(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				provider.EXPECT().
					ExecuteTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				repoReception.EXPECT().FindActive(mock.Anything, mock.Anything, mock.Anything).
					Return(reception, nil).Once()
				repoProduct.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).
					Return(nil).Once()
			},
			check: func(t *testing.T, product domain.Product, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:     "Find active Error",
			authUser: nil,
			pvzID:    pvzID,
			prepareMocks: func(provider *mocks.MockConnectionProvider, repoReception *mocks.MockReceptionsRepository, repoProduct *mocks.MockProductsRepository) {
				provider.EXPECT().
					Execute(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				repoReception.EXPECT().FindActive(mock.Anything, mock.Anything, mock.Anything).
					Return(reception, errors.New("some error")).Once()
			},
			check: func(t *testing.T, product domain.Product, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "some error")
				require.Contains(t, err.Error(), "find active failed")
			},
		},
		{
			name:     "Create Product error",
			authUser: nil,
			pvzID:    pvzID,
			prepareMocks: func(provider *mocks.MockConnectionProvider, repoReception *mocks.MockReceptionsRepository, repoProduct *mocks.MockProductsRepository) {
				provider.EXPECT().
					Execute(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				provider.EXPECT().
					ExecuteTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				repoReception.EXPECT().FindActive(mock.Anything, mock.Anything, mock.Anything).
					Return(reception, nil).Once()
				repoProduct.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("some error")).Once()
			},
			check: func(t *testing.T, product domain.Product, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "some error")
				require.Contains(t, err.Error(), "create product failed")
			},
		},
		{
			name:     "Invalid ID",
			authUser: nil,
			pvzID:    invalidPVZID,
			prepareMocks: func(provider *mocks.MockConnectionProvider, repoReception *mocks.MockReceptionsRepository, repoProduct *mocks.MockProductsRepository) {
			},
			check: func(t *testing.T, product domain.Product, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "invalid pvz id")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			provider := mocks.NewMockConnectionProvider(t)

			repoReception := mocks.NewMockReceptionsRepository(t)
			repoProduct := mocks.NewMockProductsRepository(t)

			if test.prepareMocks != nil {
				test.prepareMocks(provider, repoReception, repoProduct)
			}

			product, err := domain.NewReceptionService(provider, repoReception, repoProduct).CreateProduct(t.Context(), test.authUser, test.pvzID, test.productType)

			test.check(t, product, err)
		})
	}
}

func TestServiceReception_DeleteLastProduct(t *testing.T) {
	t.Parallel()

	pvzID := domain.PVZID(uuid.New())
	reception := domain.Reception{
		ID:     domain.ReceptionID(uuid.New()),
		PVZID:  pvzID,
		Status: domain.InProgress,
	}

	var invalidPVZID = domain.PVZID(uuid.Nil)

	tests := []struct {
		name         string
		authUser     domain.AuthenticatedUser
		pvzID        domain.PVZID
		productType  domain.ProductType
		prepareMocks func(*mocks.MockConnectionProvider, *mocks.MockReceptionsRepository, *mocks.MockProductsRepository)
		check        func(*testing.T, error)
	}{
		{
			name:     "Success",
			authUser: nil,
			pvzID:    pvzID,
			prepareMocks: func(provider *mocks.MockConnectionProvider, repoReception *mocks.MockReceptionsRepository, repoProduct *mocks.MockProductsRepository) {
				provider.EXPECT().
					Execute(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				provider.EXPECT().
					ExecuteTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				repoReception.EXPECT().FindActive(mock.Anything, mock.Anything, mock.Anything).
					Return(reception, nil).Once()
				repoProduct.EXPECT().DeleteLast(mock.Anything, mock.Anything, mock.Anything).
					Return(nil).Once()
			},
			check: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:     "Find active Error",
			authUser: nil,
			pvzID:    pvzID,
			prepareMocks: func(provider *mocks.MockConnectionProvider, repoReception *mocks.MockReceptionsRepository, repoProduct *mocks.MockProductsRepository) {
				provider.EXPECT().
					Execute(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				repoReception.EXPECT().FindActive(mock.Anything, mock.Anything, mock.Anything).
					Return(reception, errors.New("some error")).Once()
			},
			check: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "some error")
				require.Contains(t, err.Error(), "find active failed")
			},
		},
		{
			name:     "Delete last error",
			authUser: nil,
			pvzID:    pvzID,
			prepareMocks: func(provider *mocks.MockConnectionProvider, repoReception *mocks.MockReceptionsRepository, repoProduct *mocks.MockProductsRepository) {
				provider.EXPECT().
					Execute(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				provider.EXPECT().
					ExecuteTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, &mocks.MockConnection{})
					}).Once()
				repoReception.EXPECT().FindActive(mock.Anything, mock.Anything, mock.Anything).
					Return(reception, nil).Once()
				repoProduct.EXPECT().DeleteLast(mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("some error")).Once()
			},
			check: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "some error")
				require.Contains(t, err.Error(), "delete product failed")
			},
		},
		{
			name:     "Invalid ID",
			authUser: nil,
			pvzID:    invalidPVZID,
			prepareMocks: func(provider *mocks.MockConnectionProvider, repoReception *mocks.MockReceptionsRepository, repoProduct *mocks.MockProductsRepository) {
			},
			check: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "invalid pvz id")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			provider := mocks.NewMockConnectionProvider(t)

			repoReception := mocks.NewMockReceptionsRepository(t)
			repoProduct := mocks.NewMockProductsRepository(t)

			if test.prepareMocks != nil {
				test.prepareMocks(provider, repoReception, repoProduct)
			}

			err := domain.NewReceptionService(provider, repoReception, repoProduct).DeleteLastProduct(t.Context(), test.authUser, test.pvzID)

			test.check(t, err)
		})
	}
}
