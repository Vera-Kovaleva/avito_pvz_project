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
		_ = fixtureCreateProduct(ctx, t, connection, productID1, receptionID, "электроника", now.Add(-2*time.Hour))
		product2 := fixtureCreateProduct(ctx, t, connection, productID2, receptionID, "одежда", now.Add(-1*time.Hour))
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
		require.Equal(t, product2.ID, productsFound[0].ID)
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

		_ = fixtureCreateProduct(ctx, t, connection, productID1, receptionID, "электроника", now.Add(-2*time.Hour))
		product2 := fixtureCreateProduct(ctx, t, connection, productID2, receptionID, "одежда", now.Add(-1*time.Hour))
		product3 := fixtureCreateProduct(ctx, t, connection, productID3, receptionID, "обувь", now)

		limit := 1
		page := 1
		productsFound, err := products.Search(ctx, connection, nil, nil, &page, &limit)
		require.NoError(t, err)
		require.Len(t, productsFound, 1)
		require.Equal(t, product2.ID, productsFound[0].ID)

		page = 0
		productsFound, err = products.Search(ctx, connection, nil, nil, &page, &limit)
		require.NoError(t, err)
		require.Len(t, productsFound, 1)
		require.Equal(t, product3.ID, productsFound[0].ID)

		limit = 3
		to := now.Add(90 * time.Minute)
		from := now.Add(-90 * time.Minute)
		productsFound, err = products.Search(ctx, connection, &from, &to, &page, &limit)
		require.NoError(t, err)
		require.Len(t, productsFound, 2)
		require.Equal(t, product3.ID, productsFound[0].ID)
		require.Equal(t, product2.ID, productsFound[1].ID)

		from = now.Add(-3 * time.Hour)
		productsFound, err = products.Search(ctx, connection, &from, nil, &page, &limit)
		require.NoError(t, err)
		require.Len(t, productsFound, 3)

		productsFound, err = products.Search(ctx, connection, nil, &to, &page, &limit)
		require.NoError(t, err)
		require.Len(t, productsFound, 3)

		err = products.DeleteLast(ctx, connection, receptionID)
		require.NoError(t, err)

		limit = 1
		productsFound, err = products.Search(ctx, connection, nil, nil, nil, &limit)
		require.NoError(t, err)
		require.Len(t, productsFound, limit)
		require.Equal(t, product2.ID, productsFound[0].ID)

		from = now.Add(-90 * time.Minute)
		productsFound, err = products.Search(ctx, connection, &from, nil, nil, &limit)
		require.NoError(t, err)
		require.Len(t, productsFound, 1)
		require.Equal(t, product2.ID, productsFound[0].ID)
	})
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

func fixtureCreateProduct(ctx context.Context, t *testing.T, connection domain.Connection, id domain.ReceptionID, receptionID domain.ReceptionID, productType domain.ProductType, createdAt time.Time) domain.Product {
	product := domain.Product{
		ID:          id,
		ReceptionID: receptionID,
		Type:        productType,
		CreatedAt:   createdAt,
	}
	require.NoError(t, repository.NewProduct().Create(ctx, connection, product))

	return product
}
