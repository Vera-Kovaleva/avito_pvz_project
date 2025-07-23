package domain_test

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"avito_pvz/internal/domain"
	"avito_pvz/internal/generated/mocks"
	"avito_pvz/internal/infra/pointer"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServicePVZ_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		authUser     domain.AuthenticatedUser
		pvzCity      domain.PVZCity
		prepareMocks func(*mocks.MockConnection, *mocks.MockPVZsRepository, *mocks.MockMetrics)
		check        func(*testing.T, domain.PVZ, error)
	}{
		{
			name:     "Success",
			authUser: nil,
			pvzCity:  domain.Kzn,
			prepareMocks: func(_ *mocks.MockConnection, repo *mocks.MockPVZsRepository, m *mocks.MockMetrics) {
				repo.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).
					Return(nil).
					Once()
				m.EXPECT().IncPVZs().Return().Once()
			},
			check: func(t *testing.T, pvz domain.PVZ, err error) {
				require.NoError(t, err)
				require.Equal(t, domain.Kzn, pvz.City)
			},
		},
		{
			name:     "DB Error",
			authUser: nil,
			pvzCity:  domain.Kzn,
			prepareMocks: func(_ *mocks.MockConnection, repo *mocks.MockPVZsRepository, m *mocks.MockMetrics) {
				repo.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("some error")).
					Once()
			},
			check: func(t *testing.T, _ domain.PVZ, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "some error")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			connection := mocks.NewMockConnection(t)
			provider := mocks.NewMockConnectionProvider(t)

			repoPVZ := mocks.NewMockPVZsRepository(t)
			repoProduct := mocks.NewMockProductsRepository(t)
			repoReception := mocks.NewMockReceptionsRepository(t)
			metrics := mocks.NewMockMetrics(t)

			if test.prepareMocks != nil {
				test.prepareMocks(connection, repoPVZ, metrics)
			}

			provider.EXPECT().
				ExecuteTx(mock.Anything, mock.Anything).
				RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
					return f(ctx, connection)
				}).
				Once()

			pvz, err := domain.NewPVZService(provider, repoPVZ, repoProduct, repoReception, metrics).
				Create(t.Context(), test.authUser, test.pvzCity)

			test.check(t, pvz, err)
		})
	}
}

func TestServicePVZ_FindAll(t *testing.T) {
	t.Parallel()

	pvzID := uuid.New()

	tests := []struct {
		name         string
		id           domain.PVZID
		prepareMocks func(*mocks.MockConnection, *mocks.MockPVZsRepository)
		check        func(*testing.T, []domain.PVZ, error)
	}{
		{
			name: "Success",
			prepareMocks: func(_ *mocks.MockConnection, repo *mocks.MockPVZsRepository) {
				repo.EXPECT().FindAll(mock.Anything, mock.Anything).
					Return([]domain.PVZ{{ID: pvzID, City: domain.Kzn}}, nil).
					Once()
			},
			check: func(t *testing.T, pvzs []domain.PVZ, err error) {
				require.NoError(t, err)
				require.Len(t, pvzs, 1)
				require.Equal(t, pvzID, pvzs[0].ID)
			},
		},
		{
			name: "DB error",
			prepareMocks: func(_ *mocks.MockConnection, repo *mocks.MockPVZsRepository) {
				repo.EXPECT().FindAll(mock.Anything, mock.Anything).
					Return(nil, errors.New("some error")).
					Once()
			},
			check: func(t *testing.T, _ []domain.PVZ, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "some error")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			connection := mocks.NewMockConnection(t)
			provider := mocks.NewMockConnectionProvider(t)

			repoPVZ := mocks.NewMockPVZsRepository(t)
			repoProduct := mocks.NewMockProductsRepository(t)
			repoReception := mocks.NewMockReceptionsRepository(t)
			metrics := mocks.NewMockMetrics(t)

			test.prepareMocks(connection, repoPVZ)

			provider.EXPECT().Execute(mock.Anything, mock.Anything).
				RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
					return f(ctx, connection)
				}).
				Once()

			pvzs, err := domain.NewPVZService(provider, repoPVZ, repoProduct, repoReception, metrics).
				FindAll(t.Context())
			test.check(t, pvzs, err)
		})
	}
}

func TestServicePVZ_FindPVZReceptionProducts(t *testing.T) {
	t.Parallel()

	now := time.Now()
	pvzID1, pvzID2 := uuid.New(), uuid.New()
	receptionID1, receptionID2 := uuid.New(), uuid.New()
	productID1, productID2 := uuid.New(), uuid.New()

	tests := []struct {
		name         string
		authUser     domain.AuthenticatedUser
		from, to     *time.Time
		page, limit  *int
		prepareMocks func(*mocks.MockConnection, *mocks.MockProductsRepository, *mocks.MockReceptionsRepository, *mocks.MockPVZsRepository)
		check        func(*testing.T, []domain.PVZReceptionsProducts, error)
	}{
		{
			name:     "Success one product, one reception, one pvz",
			authUser: nil,
			from:     &now,
			to:       &now,
			page:     pointer.Ref(1),
			limit:    pointer.Ref(10),
			prepareMocks: func(
				_ *mocks.MockConnection,
				productRepo *mocks.MockProductsRepository,
				receptionRepo *mocks.MockReceptionsRepository,
				pvzRepo *mocks.MockPVZsRepository,
			) {
				products := []domain.Product{
					{
						ID:          productID1,
						ReceptionID: receptionID1,
						Type:        "electronics",
						CreatedAt:   now,
					},
				}
				productRepo.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(products, nil).
					Once()

				receptions := []domain.Reception{
					{ID: receptionID1, PVZID: pvzID1, Status: domain.InProgress, CreatedAt: now},
				}
				receptionRepo.EXPECT().
					FindByIDs(mock.Anything, mock.Anything, mock.Anything).
					Return(receptions, nil).
					Once()

				pvzs := []domain.PVZ{
					{ID: pvzID1, City: domain.Kzn, RegisteredAt: now},
				}
				pvzRepo.EXPECT().
					FindByIDs(mock.Anything, mock.Anything, mock.Anything).
					Return(pvzs, nil).
					Once()
			},
			check: func(t *testing.T, result []domain.PVZReceptionsProducts, err error) {
				require.NoError(t, err)
				require.Len(t, result, 1)
				require.Equal(t, pvzID1, result[0].PVZ.ID)
				require.Len(t, result[0].Receptions, 1)
				require.Equal(t, receptionID1, result[0].Receptions[0].Reception.ID)
				require.Len(t, result[0].Receptions[0].Products, 1)
				require.Equal(t, productID1, result[0].Receptions[0].Products[0].ID)
			},
		},
		{
			name:     "Success two products, one reception, one pvz",
			authUser: nil,
			from:     &now,
			to:       &now,
			page:     pointer.Ref(1),
			limit:    pointer.Ref(10),
			prepareMocks: func(
				_ *mocks.MockConnection,
				productRepo *mocks.MockProductsRepository,
				receptionRepo *mocks.MockReceptionsRepository,
				pvzRepo *mocks.MockPVZsRepository,
			) {
				products := []domain.Product{
					{
						ID:          productID1,
						ReceptionID: receptionID1,
						Type:        "электроника",
						CreatedAt:   now,
					},
					{
						ID:          productID2,
						ReceptionID: receptionID1,
						Type:        "одежда",
						CreatedAt:   now.Add(-2 * time.Hour),
					},
				}
				productRepo.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(products, nil).
					Once()

				receptions := []domain.Reception{
					{ID: receptionID1, PVZID: pvzID1, Status: domain.InProgress, CreatedAt: now},
				}
				receptionRepo.EXPECT().
					FindByIDs(mock.Anything, mock.Anything, mock.Anything).
					Return(receptions, nil).
					Once()

				pvzs := []domain.PVZ{
					{ID: pvzID1, City: domain.Kzn, RegisteredAt: now},
				}
				pvzRepo.EXPECT().
					FindByIDs(mock.Anything, mock.Anything, mock.Anything).
					Return(pvzs, nil).
					Once()
			},
			check: func(t *testing.T, result []domain.PVZReceptionsProducts, err error) {
				require.NoError(t, err)
				require.Len(t, result, 1)
				require.Equal(t, pvzID1, result[0].PVZ.ID)
				require.Len(t, result[0].Receptions, 1)
				require.Equal(t, receptionID1, result[0].Receptions[0].Reception.ID)

				products := result[0].Receptions[0].Products
				slices.SortFunc(products, func(p1, p2 domain.Product) int {
					return p1.CreatedAt.Compare(p2.CreatedAt)
				})

				require.Len(t, products, 2)
				require.Equal(t, now.Add(-2*time.Hour), products[0].CreatedAt)
				require.Equal(t, productID1, products[1].ID)
			},
		},
		{
			name:     "Success two products, two reception, one pvz",
			authUser: nil,
			from:     &now,
			to:       &now,
			page:     pointer.Ref(1),
			limit:    pointer.Ref(10),
			prepareMocks: func(
				_ *mocks.MockConnection,
				productRepo *mocks.MockProductsRepository,
				receptionRepo *mocks.MockReceptionsRepository,
				pvzRepo *mocks.MockPVZsRepository,
			) {
				products := []domain.Product{
					{
						ID:          productID2,
						ReceptionID: receptionID1,
						Type:        "одежда",
						CreatedAt:   now.Add(-2 * time.Hour),
					},
					{
						ID:          productID1,
						ReceptionID: receptionID1,
						Type:        "электроника",
						CreatedAt:   now,
					},
				}
				productRepo.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(products, nil).
					Once()

				receptions := []domain.Reception{
					{
						ID:        receptionID1,
						PVZID:     pvzID1,
						Status:    domain.Close,
						CreatedAt: now.Add(-2 * time.Hour),
					},
					{ID: receptionID2, PVZID: pvzID1, Status: domain.InProgress, CreatedAt: now},
				}
				receptionRepo.EXPECT().
					FindByIDs(mock.Anything, mock.Anything, mock.Anything).
					Return(receptions, nil).
					Once()

				pvzs := []domain.PVZ{
					{ID: pvzID1, City: domain.Kzn, RegisteredAt: now},
				}
				pvzRepo.EXPECT().
					FindByIDs(mock.Anything, mock.Anything, mock.Anything).
					Return(pvzs, nil).
					Once()
			},
			check: func(t *testing.T, result []domain.PVZReceptionsProducts, err error) {
				require.NoError(t, err)
				require.Len(t, result, 1)
				require.Equal(t, pvzID1, result[0].PVZ.ID)

				receptions := result[0].Receptions
				slices.SortFunc(receptions, func(r1, r2 domain.ReceptionsProducts) int {
					return r1.Reception.CreatedAt.Compare(r2.Reception.CreatedAt)
				})
				require.Len(t, receptions, 2)
				require.Equal(t, domain.Close, receptions[0].Reception.Status)
				require.Equal(t, receptionID2, receptions[1].Reception.ID)

				products := result[0].Receptions[0].Products
				slices.SortFunc(products, func(p1, p2 domain.Product) int {
					return p1.CreatedAt.Compare(p2.CreatedAt)
				})
				require.Len(t, products, 2)
				require.Equal(t, now.Add(-2*time.Hour), products[0].CreatedAt)
				require.Equal(t, productID1, products[1].ID)
			},
		},
		{
			name:     "Success two products, two reception, two pvz",
			authUser: nil,
			from:     &now,
			to:       &now,
			page:     pointer.Ref(1),
			limit:    pointer.Ref(10),
			prepareMocks: func(
				_ *mocks.MockConnection,
				productRepo *mocks.MockProductsRepository,
				receptionRepo *mocks.MockReceptionsRepository,
				pvzRepo *mocks.MockPVZsRepository,
			) {
				products := []domain.Product{
					{
						ID:          productID2,
						ReceptionID: receptionID1,
						Type:        "одежда",
						CreatedAt:   now.Add(-2 * time.Hour),
					},
					{
						ID:          productID1,
						ReceptionID: receptionID1,
						Type:        "электроника",
						CreatedAt:   now,
					},
				}
				productRepo.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(products, nil).
					Once()

				receptions := []domain.Reception{
					{
						ID:        receptionID1,
						PVZID:     pvzID1,
						Status:    domain.Close,
						CreatedAt: now.Add(-2 * time.Hour),
					},
					{ID: receptionID2, PVZID: pvzID1, Status: domain.InProgress, CreatedAt: now},
				}
				receptionRepo.EXPECT().
					FindByIDs(mock.Anything, mock.Anything, mock.Anything).
					Return(receptions, nil).
					Once()

				pvzs := []domain.PVZ{
					{ID: pvzID1, City: domain.Kzn, RegisteredAt: now.Add(-2 * time.Hour)},
					{ID: pvzID2, City: domain.Msk, RegisteredAt: now},
				}
				pvzRepo.EXPECT().
					FindByIDs(mock.Anything, mock.Anything, mock.Anything).
					Return(pvzs, nil).
					Once()
			},
			check: func(t *testing.T, result []domain.PVZReceptionsProducts, err error) {
				require.NoError(t, err)

				pvzs := result
				slices.SortFunc(pvzs, func(pvz1, pvz2 domain.PVZReceptionsProducts) int {
					return pvz1.PVZ.RegisteredAt.Compare(pvz2.PVZ.RegisteredAt)
				})
				require.Len(t, result, 2)
				require.Equal(t, now.Add(-2*time.Hour), result[0].PVZ.RegisteredAt)
				require.Equal(t, pvzID2, result[1].PVZ.ID)

				receptions := result[0].Receptions
				slices.SortFunc(receptions, func(r1, r2 domain.ReceptionsProducts) int {
					return r1.Reception.CreatedAt.Compare(r2.Reception.CreatedAt)
				})
				require.Len(t, receptions, 2)
				require.Equal(t, domain.Close, receptions[0].Reception.Status)
				require.Equal(t, receptionID2, receptions[1].Reception.ID)

				products := result[0].Receptions[0].Products
				slices.SortFunc(products, func(p1, p2 domain.Product) int {
					return p1.CreatedAt.Compare(p2.CreatedAt)
				})
				require.Len(t, products, 2)
				require.Equal(t, now.Add(-2*time.Hour), products[0].CreatedAt)
				require.Equal(t, productID1, products[1].ID)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			connection := mocks.NewMockConnection(t)
			provider := mocks.NewMockConnectionProvider(t)

			repoPVZ := mocks.NewMockPVZsRepository(t)
			repoProduct := mocks.NewMockProductsRepository(t)
			repoReception := mocks.NewMockReceptionsRepository(t)
			metrics := mocks.NewMockMetrics(t)

			test.prepareMocks(connection, repoProduct, repoReception, repoPVZ)

			provider.EXPECT().Execute(mock.Anything, mock.Anything).
				RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
					return f(ctx, connection)
				}).
				Times(3)

			result, err := domain.NewPVZService(provider, repoPVZ, repoProduct, repoReception, metrics).
				FindPVZReceptionProducts(t.Context(), test.authUser, test.from, test.to, test.page, test.limit)
			test.check(t, result, err)
		})
	}
}

func TestServicePVZ_FindPVZReceptionProducts_ErrProducts(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name         string
		authUser     domain.AuthenticatedUser
		from, to     *time.Time
		page, limit  *int
		prepareMocks func(*mocks.MockConnection, *mocks.MockProductsRepository, *mocks.MockReceptionsRepository, *mocks.MockPVZsRepository)
		check        func(*testing.T, []domain.PVZReceptionsProducts, error)
	}{
		{
			name:     "DB Products Error",
			authUser: nil,
			from:     &now,
			to:       &now,
			page:     pointer.Ref(1),
			limit:    pointer.Ref(10),
			prepareMocks: func(
				_ *mocks.MockConnection,
				productRepo *mocks.MockProductsRepository,
				_ *mocks.MockReceptionsRepository,
				_ *mocks.MockPVZsRepository,
			) {
				products := []domain.Product{}
				productRepo.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(products, errors.New("some error")).
					Once()
			},
			check: func(t *testing.T, _ []domain.PVZReceptionsProducts, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "some error")
				require.Contains(t, err.Error(), "search products failed")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			connection := mocks.NewMockConnection(t)
			provider := mocks.NewMockConnectionProvider(t)

			repoPVZ := mocks.NewMockPVZsRepository(t)
			repoProduct := mocks.NewMockProductsRepository(t)
			repoReception := mocks.NewMockReceptionsRepository(t)
			metrics := mocks.NewMockMetrics(t)

			test.prepareMocks(connection, repoProduct, repoReception, repoPVZ)

			provider.EXPECT().Execute(mock.Anything, mock.Anything).
				RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
					return f(ctx, connection)
				}).
				Once()

			result, err := domain.NewPVZService(provider, repoPVZ, repoProduct, repoReception, metrics).
				FindPVZReceptionProducts(t.Context(), test.authUser, test.from, test.to, test.page, test.limit)
			test.check(t, result, err)
		})
	}
}

func TestServicePVZ_FindPVZReceptionProducts_ErrReceptions(t *testing.T) {
	t.Parallel()

	now := time.Now()
	pvzID := uuid.New()
	receptionID := uuid.New()
	productID1 := uuid.New()

	tests := []struct {
		name         string
		authUser     domain.AuthenticatedUser
		from, to     *time.Time
		page, limit  *int
		prepareMocks func(*mocks.MockConnection, *mocks.MockProductsRepository, *mocks.MockReceptionsRepository, *mocks.MockPVZsRepository)
		check        func(*testing.T, []domain.PVZReceptionsProducts, error)
	}{
		{
			prepareMocks: func(_ *mocks.MockConnection,
				productRepo *mocks.MockProductsRepository,
				receptionRepo *mocks.MockReceptionsRepository,
				_ *mocks.MockPVZsRepository,
			) {
				products := []domain.Product{
					{ID: productID1, ReceptionID: receptionID, Type: "electronics", CreatedAt: now},
				}
				productRepo.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(products, nil).
					Once()

				receptions := []domain.Reception{
					{ID: receptionID, PVZID: pvzID, Status: domain.InProgress, CreatedAt: now},
				}
				receptionRepo.EXPECT().
					FindByIDs(mock.Anything, mock.Anything, mock.Anything).
					Return(receptions, errors.New("some error")).
					Once()
			},
			check: func(t *testing.T, _ []domain.PVZReceptionsProducts, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "some error")
				require.Contains(t, err.Error(), "search reseptions failed")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			connection := mocks.NewMockConnection(t)
			provider := mocks.NewMockConnectionProvider(t)

			repoPVZ := mocks.NewMockPVZsRepository(t)
			repoProduct := mocks.NewMockProductsRepository(t)
			repoReception := mocks.NewMockReceptionsRepository(t)
			metrics := mocks.NewMockMetrics(t)

			test.prepareMocks(connection, repoProduct, repoReception, repoPVZ)

			provider.EXPECT().Execute(mock.Anything, mock.Anything).
				RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
					return f(ctx, connection)
				}).
				Times(2)

			result, err := domain.NewPVZService(provider, repoPVZ, repoProduct, repoReception, metrics).
				FindPVZReceptionProducts(t.Context(), test.authUser, test.from, test.to, test.page, test.limit)
			test.check(t, result, err)
		})
	}
}

func TestServicePVZ_FindPVZReceptionProducts_ErrPVZs(t *testing.T) {
	t.Parallel()

	now := time.Now()
	pvzID := uuid.New()
	receptionID := uuid.New()
	productID1 := uuid.New()

	tests := []struct {
		name         string
		authUser     domain.AuthenticatedUser
		from, to     *time.Time
		page, limit  *int
		prepareMocks func(*mocks.MockConnection, *mocks.MockProductsRepository, *mocks.MockReceptionsRepository, *mocks.MockPVZsRepository)
		check        func(*testing.T, []domain.PVZReceptionsProducts, error)
	}{
		{
			prepareMocks: func(_ *mocks.MockConnection,
				productRepo *mocks.MockProductsRepository,
				receptionRepo *mocks.MockReceptionsRepository,
				pvzRepo *mocks.MockPVZsRepository,
			) {
				products := []domain.Product{
					{ID: productID1, ReceptionID: receptionID, Type: "electronics", CreatedAt: now},
				}
				productRepo.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(products, nil).
					Once()

				receptions := []domain.Reception{
					{ID: receptionID, PVZID: pvzID, Status: domain.InProgress, CreatedAt: now},
				}
				receptionRepo.EXPECT().
					FindByIDs(mock.Anything, mock.Anything, mock.Anything).
					Return(receptions, nil).
					Once()

				pvzs := []domain.PVZ{
					{ID: pvzID, City: domain.Kzn, RegisteredAt: now},
				}
				pvzRepo.EXPECT().
					FindByIDs(mock.Anything, mock.Anything, mock.Anything).
					Return(pvzs, errors.New("some error")).
					Once()
			},
			check: func(t *testing.T, _ []domain.PVZReceptionsProducts, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "some error")
				require.Contains(t, err.Error(), "search pvzs failed")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			connection := mocks.NewMockConnection(t)
			provider := mocks.NewMockConnectionProvider(t)

			repoPVZ := mocks.NewMockPVZsRepository(t)
			repoProduct := mocks.NewMockProductsRepository(t)
			repoReception := mocks.NewMockReceptionsRepository(t)
			metrics := mocks.NewMockMetrics(t)

			test.prepareMocks(connection, repoProduct, repoReception, repoPVZ)

			provider.EXPECT().Execute(mock.Anything, mock.Anything).
				RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
					return f(ctx, connection)
				}).
				Times(3)

			result, err := domain.NewPVZService(provider, repoPVZ, repoProduct, repoReception, metrics).
				FindPVZReceptionProducts(t.Context(), test.authUser, test.from, test.to, test.page, test.limit)
			test.check(t, result, err)
		})
	}
}
