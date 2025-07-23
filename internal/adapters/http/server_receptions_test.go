package http_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"avito_pvz/internal/adapters/http"
	"avito_pvz/internal/domain"
	"avito_pvz/internal/generated/mocks"
	oapi "avito_pvz/internal/generated/oapi"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServer_PostPvzPvzIdCloseLastReception(t *testing.T) {
	t.Parallel()

	pvzID := uuid.New()

	reseption := domain.Reception{
		ID:        uuid.New(),
		PVZID:     pvzID,
		Status:    domain.Close,
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name         string
		request      oapi.PostPvzPvzIdCloseLastReceptionRequestObject
		prepareMocks func(*mocks.MockConnectionProvider, *mocks.MockReceptionsRepository)
		check        func(*testing.T, oapi.PostPvzPvzIdCloseLastReceptionResponseObject, error)
	}{
		{
			name:    "Success",
			request: oapi.PostPvzPvzIdCloseLastReceptionRequestObject{PvzId: pvzID},
			prepareMocks: func(provider *mocks.MockConnectionProvider, repo *mocks.MockReceptionsRepository) {
				provider.EXPECT().
					Execute(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, nil)
					})
				provider.EXPECT().
					ExecuteTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, nil)
					})
				repo.EXPECT().
					FindActive(mock.Anything, mock.Anything, mock.Anything).
					Return(reseption, nil)
				repo.EXPECT().
					Close(mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
			},
			check: func(t *testing.T, response oapi.PostPvzPvzIdCloseLastReceptionResponseObject, err error) {
				require.NoError(t, err)
				require.IsType(
					t,
					oapi.PostPvzPvzIdCloseLastReception200JSONResponse{},
					response.(oapi.PostPvzPvzIdCloseLastReception200JSONResponse),
				)
			},
		},
		{
			name:    "Error",
			request: oapi.PostPvzPvzIdCloseLastReceptionRequestObject{PvzId: pvzID},
			prepareMocks: func(provider *mocks.MockConnectionProvider, repo *mocks.MockReceptionsRepository) {
				provider.EXPECT().
					Execute(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, nil)
					})
				provider.EXPECT().
					ExecuteTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, nil)
					})
				repo.EXPECT().
					FindActive(mock.Anything, mock.Anything, mock.Anything).
					Return(reseption, nil)
				repo.EXPECT().
					Close(mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("some error"))
			},
			check: func(t *testing.T, response oapi.PostPvzPvzIdCloseLastReceptionResponseObject, err error) {
				require.NoError(t, err)
				res, _ := response.(oapi.PostPvzPvzIdCloseLastReception400JSONResponse)
				assert.Contains(t, res.Message, "Неверный запрос или приемка уже закрыта")
				require.IsType(
					t,
					oapi.PostPvzPvzIdCloseLastReception400JSONResponse{},
					response.(oapi.PostPvzPvzIdCloseLastReception400JSONResponse),
				)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			connection := mocks.NewMockConnectionProvider(t)
			productRepo := mocks.NewMockProductsRepository(t)
			receptionRepo := mocks.NewMockReceptionsRepository(t)
			metrics := mocks.NewMockMetrics(t)

			if test.prepareMocks != nil {
				test.prepareMocks(connection, receptionRepo)
			}

			server := http.NewServer(
				nil,
				domain.NewReceptionService(connection, receptionRepo, productRepo, metrics),
				nil,
			)

			response, err := server.PostPvzPvzIdCloseLastReception(t.Context(), test.request)
			test.check(t, response, err)
		})
	}
}

func TestServer_PostPvzPvzIdDeleteLastProduct(t *testing.T) {
	t.Parallel()

	pvzID := uuid.New()

	reseption := domain.Reception{
		ID:        uuid.New(),
		PVZID:     pvzID,
		Status:    domain.Close,
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name         string
		request      oapi.PostPvzPvzIdDeleteLastProductRequestObject
		prepareMocks func(*mocks.MockConnectionProvider, *mocks.MockReceptionsRepository, *mocks.MockProductsRepository)
		check        func(*testing.T, oapi.PostPvzPvzIdDeleteLastProductResponseObject, error)
	}{
		{
			name:    "Success",
			request: oapi.PostPvzPvzIdDeleteLastProductRequestObject{PvzId: pvzID},
			prepareMocks: func(provider *mocks.MockConnectionProvider, repoReception *mocks.MockReceptionsRepository, repoProduct *mocks.MockProductsRepository) {
				provider.EXPECT().
					Execute(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, nil)
					})
				provider.EXPECT().
					ExecuteTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, nil)
					})
				repoReception.EXPECT().
					FindActive(mock.Anything, mock.Anything, mock.Anything).
					Return(reseption, nil)
				repoProduct.EXPECT().
					DeleteLast(mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
			},
			check: func(t *testing.T, response oapi.PostPvzPvzIdDeleteLastProductResponseObject, err error) {
				require.NoError(t, err)
				require.IsType(
					t,
					oapi.PostPvzPvzIdDeleteLastProduct200Response{},
					response.(oapi.PostPvzPvzIdDeleteLastProduct200Response),
				)
			},
		},
		{
			name:    "Error",
			request: oapi.PostPvzPvzIdDeleteLastProductRequestObject{PvzId: pvzID},
			prepareMocks: func(provider *mocks.MockConnectionProvider, repoReception *mocks.MockReceptionsRepository, repoProduct *mocks.MockProductsRepository) {
				provider.EXPECT().
					Execute(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, nil)
					})
				provider.EXPECT().
					ExecuteTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, nil)
					})
				repoReception.EXPECT().
					FindActive(mock.Anything, mock.Anything, mock.Anything).
					Return(reseption, nil)
				repoProduct.EXPECT().
					DeleteLast(mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("some error"))
			},
			check: func(t *testing.T, response oapi.PostPvzPvzIdDeleteLastProductResponseObject, err error) {
				require.NoError(t, err)
				require.IsType(
					t,
					oapi.PostPvzPvzIdDeleteLastProduct400JSONResponse{},
					response.(oapi.PostPvzPvzIdDeleteLastProduct400JSONResponse),
				)
				res, _ := response.(oapi.PostPvzPvzIdDeleteLastProduct400JSONResponse)
				assert.Contains(
					t,
					res.Message,
					"Неверный запрос, нет активной приемки или нет товаров для удаления",
				)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			connection := mocks.NewMockConnectionProvider(t)
			productRepo := mocks.NewMockProductsRepository(t)
			receptionRepo := mocks.NewMockReceptionsRepository(t)
			metrics := mocks.NewMockMetrics(t)

			if test.prepareMocks != nil {
				test.prepareMocks(connection, receptionRepo, productRepo)
			}

			server := http.NewServer(
				nil,
				domain.NewReceptionService(connection, receptionRepo, productRepo, metrics),
				nil,
			)

			response, err := server.PostPvzPvzIdDeleteLastProduct(t.Context(), test.request)
			test.check(t, response, err)
		})
	}
}

func TestServer_PostProducts(t *testing.T) {
	t.Parallel()

	pvzID := uuid.New()

	reseption := domain.Reception{
		ID:        uuid.New(),
		PVZID:     pvzID,
		Status:    domain.Close,
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name         string
		request      oapi.PostProductsRequestObject
		prepareMocks func(*mocks.MockConnectionProvider, *mocks.MockReceptionsRepository, *mocks.MockProductsRepository, *mocks.MockMetrics)
		check        func(oapi.PostProductsResponseObject, error)
	}{
		{
			name: "Success",
			request: oapi.PostProductsRequestObject{Body: &oapi.PostProductsJSONRequestBody{
				PvzId: uuid.New(),
				Type:  oapi.PostProductsJSONBodyType(oapi.ProductTypeОбувь),
			}},
			prepareMocks: func(provider *mocks.MockConnectionProvider, repoReception *mocks.MockReceptionsRepository, repoProduct *mocks.MockProductsRepository, m *mocks.MockMetrics) {
				provider.EXPECT().
					Execute(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, nil)
					})
				provider.EXPECT().
					ExecuteTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, nil)
					})
				m.EXPECT().IncProducts().Return().Once()
				repoReception.EXPECT().
					FindActive(mock.Anything, mock.Anything, mock.Anything).
					Return(reseption, nil)
				repoProduct.EXPECT().
					Create(mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
			},
			check: func(response oapi.PostProductsResponseObject, err error) {
				require.NoError(t, err)
				require.IsType(
					t,
					oapi.PostProducts201JSONResponse{},
					response.(oapi.PostProducts201JSONResponse),
				)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			connection := mocks.NewMockConnectionProvider(t)
			repoReception := mocks.NewMockReceptionsRepository(t)
			repoProduct := mocks.NewMockProductsRepository(t)
			metrics := mocks.NewMockMetrics(t)

			if test.prepareMocks != nil {
				test.prepareMocks(connection, repoReception, repoProduct, metrics)
			}

			server := http.NewServer(
				nil,
				domain.NewReceptionService(connection, repoReception, repoProduct, metrics),
				nil,
			)

			response, err := server.PostProducts(t.Context(), test.request)
			test.check(response, err)
		})
	}
}
