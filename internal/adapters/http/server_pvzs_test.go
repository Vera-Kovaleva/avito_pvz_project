package http_test

import (
	"context"
	"errors"
	"testing"

	"avito_pvz/internal/adapters/http"
	"avito_pvz/internal/domain"
	"avito_pvz/internal/generated/mocks"
	oapi "avito_pvz/internal/generated/oapi"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServer_PostPvz(t *testing.T) {
	t.Parallel()

	testCity := oapi.Москва

	tests := []struct {
		name         string
		request      oapi.PostPvzRequestObject
		prepareMocks func(*mocks.MockConnectionProvider, *mocks.MockPVZsRepository, *mocks.MockMetrics)
		check        func(*testing.T, oapi.PostPvzResponseObject, error)
	}{
		{
			name:    "Success",
			request: oapi.PostPvzRequestObject{Body: &oapi.PVZ{City: testCity}},
			prepareMocks: func(provider *mocks.MockConnectionProvider, repo *mocks.MockPVZsRepository, m *mocks.MockMetrics) {
				provider.EXPECT().
					ExecuteTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, nil)
					})
				repo.EXPECT().
					Create(mock.Anything, mock.Anything, mock.Anything).
					Return(nil)

				m.EXPECT().IncPVZs().Return().Once()
			},
			check: func(t *testing.T, response oapi.PostPvzResponseObject, err error) {
				require.NoError(t, err)
				require.IsType(
					t,
					oapi.PostPvz201JSONResponse{},
					response.(oapi.PostPvz201JSONResponse),
				)
			},
		},
		{
			name:    "ServiceError",
			request: oapi.PostPvzRequestObject{Body: &oapi.PVZ{City: testCity}},
			prepareMocks: func(provider *mocks.MockConnectionProvider, repo *mocks.MockPVZsRepository, m *mocks.MockMetrics) {
				provider.EXPECT().
					ExecuteTx(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, f func(context.Context, domain.Connection) error) error {
						return f(ctx, nil)
					})
				repo.EXPECT().
					Create(mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("some error"))
			},
			check: func(t *testing.T, response oapi.PostPvzResponseObject, err error) {
				require.NoError(t, err)
				res, _ := response.(oapi.PostPvz400JSONResponse)
				assert.Equal(t, "Неверный запрос", res.Message)
				require.IsType(
					t,
					oapi.PostPvz400JSONResponse{},
					response.(oapi.PostPvz400JSONResponse),
				)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			conn := mocks.NewMockConnectionProvider(t)
			pvzRepo := mocks.NewMockPVZsRepository(t)
			productRepo := mocks.NewMockProductsRepository(t)
			receptionRepo := mocks.NewMockReceptionsRepository(t)
			metrics := mocks.NewMockMetrics(t)

			if test.prepareMocks != nil {
				test.prepareMocks(conn, pvzRepo, metrics)
			}

			server := http.NewServer(
				domain.NewPVZService(conn, pvzRepo, productRepo, receptionRepo, metrics),
				nil,
				nil,
			)

			response, err := server.PostPvz(t.Context(), test.request)
			test.check(t, response, err)
		})
	}
}
