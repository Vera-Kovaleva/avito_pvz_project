package http_test

import (
	"context"
	"testing"
	"time"

	"avito_pvz/internal/adapters/http"
	"avito_pvz/internal/domain"
	"avito_pvz/internal/generated/mocks"
	oapi "avito_pvz/internal/generated/oapi"

	"github.com/google/uuid"
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			connection := mocks.NewMockConnectionProvider(t)
			productRepo := mocks.NewMockProductsRepository(t)
			receptionRepo := mocks.NewMockReceptionsRepository(t)

			if test.prepareMocks != nil {
				test.prepareMocks(connection, receptionRepo)
			}

			server := http.NewServer(
				nil,
				domain.NewReceptionService(connection, receptionRepo, productRepo),
				nil,
			)

			response, err := server.PostPvzPvzIdCloseLastReception(t.Context(), test.request)
			test.check(t, response, err)
		})
	}
}
