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

func TestProductIntegration(t *testing.T) {
	rollback(t, func(ctx context.Context, connection domain.Connection) {
		pvzID := uuid.New()
		receptionID := uuid.New()

		_ = fixtureCreatePVZ(ctx, t, connection, pvzID, "Казань")
		_ = fixtureCreateReceptin(ctx, t, connection, receptionID, pvzID)

		products := repository.NewProduct()
		productID1 := uuid.New()
		productID2 := uuid.New()
		productID3 := uuid.New()
		now := time.Now()
		_ = fixtureCreateProduct(
			ctx,
			t,
			connection,
			productID1,
			receptionID,
			"электроника",
			now.Add(-2*time.Hour),
		)
		product2 := fixtureCreateProduct(
			ctx,
			t,
			connection,
			productID2,
			receptionID,
			"одежда",
			now.Add(-1*time.Hour),
		)
		_ = fixtureCreateProduct(ctx, t, connection, productID3, receptionID, "обувь", now)

		err := products.DeleteLast(ctx, connection, receptionID)
		require.NoError(t, err)

		productsFound, err := products.Search(ctx, connection, nil, nil, nil, nil)
		require.NoError(t, err)
		require.Len(t, productsFound, 2)

		limit := 1
		productsFound, err = products.Search(ctx, connection, nil, nil, nil, &limit)
		require.NoError(t, err)
		require.Len(t, productsFound, limit)
		require.Equal(t, product2.ID, productsFound[1].ID)
	})
}

func TestProductIntegrationSearch(t *testing.T) {
	rollback(t, func(ctx context.Context, connection domain.Connection) {
		pvzID := uuid.New()
		receptionID := uuid.New()

		_ = fixtureCreatePVZ(ctx, t, connection, pvzID, "Казань")
		_ = fixtureCreateReceptin(ctx, t, connection, receptionID, pvzID)

		products := repository.NewProduct()
		productID1 := uuid.New()
		productID2 := uuid.New()
		productID3 := uuid.New()

		now := time.Now()

		product1 := fixtureCreateProduct(
			ctx,
			t,
			connection,
			productID1,
			receptionID,
			"электроника",
			now.Add(-2*time.Hour),
		)
		product2 := fixtureCreateProduct(
			ctx,
			t,
			connection,
			productID2,
			receptionID,
			"одежда",
			now.Add(-1*time.Hour),
		)
		product3 := fixtureCreateProduct(ctx, t, connection, productID3, receptionID, "обувь", now)

		limit := 1
		page := 0
		productsFound, err := products.Search(ctx, connection, nil, nil, &page, &limit)
		require.NoError(t, err)
		require.Len(t, productsFound, 1)
		require.Equal(t, product1.ID, productsFound[0].ID)

		page = 1
		productsFound, err = products.Search(ctx, connection, nil, nil, &page, &limit)
		require.NoError(t, err)
		require.Len(t, productsFound, 1)
		require.Equal(t, product2.ID, productsFound[0].ID)

		limit = 3
		page = 0
		to := now.Add(4 * time.Hour)
		from := now.Add(-4 * time.Hour)
		productsFound, err = products.Search(ctx, connection, &from, nil, &page, &limit)
		require.NoError(t, err)
		require.Len(t, productsFound, 3)

		productsFound, err = products.Search(ctx, connection, nil, &to, &page, &limit)
		require.NoError(t, err)
		require.Len(t, productsFound, 3)

		productsFound, err = products.Search(ctx, connection, &from, &to, &page, &limit)
		require.NoError(t, err)
		require.Len(t, productsFound, 3)
		require.Equal(t, product1.ID, productsFound[0].ID)

		limit = 1
		productsFound, err = products.Search(ctx, connection, nil, nil, nil, &limit)
		require.NoError(t, err)
		require.Len(t, productsFound, limit)
		require.Equal(t, product3.ID, productsFound[2].ID)
	})
}

func TestProductSearchErrors(t *testing.T) {
	connection := mocks.NewMockConnection(t)

	from := time.Now()
	to := from.Add(-time.Hour)
	_, err := repository.NewProduct().Search(t.Context(), connection, &from, &to, nil, nil)
	require.ErrorIs(t, err, repository.ErrSearchProduct)
	require.ErrorContains(t, err, "from must be less than to")

	page := -1
	limit := 1
	_, err = repository.NewProduct().Search(t.Context(), connection, nil, nil, &page, &limit)
	require.ErrorIs(t, err, repository.ErrSearchProduct)
	require.ErrorContains(t, err, "invalid page")

	page = 1
	limit = 0
	_, err = repository.NewProduct().Search(t.Context(), connection, nil, nil, &page, &limit)
	require.ErrorIs(t, err, repository.ErrSearchProduct)
	require.ErrorContains(t, err, "invalid limit")

	limit = 1
	_, err = repository.NewProduct().Search(t.Context(), connection, nil, nil, &page, nil)
	require.ErrorIs(t, err, repository.ErrSearchProduct)
	require.ErrorContains(t, err, "page without limit")
}

func TestProductUnitSearch(t *testing.T) {
	connection := mocks.NewMockConnection(t)

	connection.EXPECT().
		SelectContext(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("some error")).
		Once()

	_, err := repository.NewProduct().Search(t.Context(), connection, nil, nil, nil, nil)
	require.ErrorIs(t, err, repository.ErrSearchProduct)
	require.ErrorContains(t, err, "some error")
}

func TestProductUnitCreate(t *testing.T) {
	connection := mocks.NewMockConnection(t)

	connection.EXPECT().
		ExecContext(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(0, errors.New("some error")).
		Once()

	err := repository.NewProduct().Create(t.Context(), connection, domain.Product{})
	require.ErrorIs(t, err, repository.ErrCreateProduct)
	require.ErrorContains(t, err, "some error")
}

func TestProductUnitDelete(t *testing.T) {
	connection := mocks.NewMockConnection(t)

	connection.EXPECT().ExecContext(mock.Anything, mock.Anything, mock.Anything).
		Return(0, errors.New("some error")).
		Once()

	var receptionID domain.ReceptionID
	err := repository.NewProduct().DeleteLast(t.Context(), connection, receptionID)
	require.ErrorIs(t, err, repository.ErrDeleteProduct)
	require.ErrorContains(t, err, "some error")
}

func fixtureCreateProduct(
	ctx context.Context,
	t *testing.T,
	connection domain.Connection,
	id domain.ReceptionID,
	receptionID domain.ReceptionID,
	productType domain.ProductType,
	createdAt time.Time,
) domain.Product {
	product := domain.Product{
		ID:          id,
		ReceptionID: receptionID,
		Type:        productType,
		CreatedAt:   createdAt,
	}
	require.NoError(t, repository.NewProduct().Create(ctx, connection, product))

	return product
}

func TestUnitProducts_Search(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		from, to     *time.Time
		page, limit  *int
		prepareMocks func(*mocks.MockConnection)
		check        func(*testing.T, []domain.Product, error)
	}{
		{
			name: "Success - no params",
			prepareMocks: func(connection *mocks.MockConnection) {
				const expectedQuery = "select id, reception_id, type, created_at from products order by created_at"
				connection.EXPECT().
					SelectContext(mock.Anything, mock.Anything, expectedQuery).
					Return(nil).
					Once()
			},
			check: func(t *testing.T, _ []domain.Product, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			connection := mocks.NewMockConnection(t)
			if test.prepareMocks != nil {
				test.prepareMocks(connection)
			}

			products, err := repository.NewProduct().
				Search(t.Context(), connection, test.from, test.to, test.page, test.limit)

			if test.check != nil {
				test.check(t, products, err)
			}
		})
	}
}
