package http

import (
	"context"

	"avito_pvz/internal/domain"
	oapi "avito_pvz/internal/generated/oapi"
	"avito_pvz/internal/infra/pointer"
)

func (s *Server) PostPvz(
	ctx context.Context,
	request oapi.PostPvzRequestObject,
) (oapi.PostPvzResponseObject, error) {
	pvz, err := s.pvzs.Create(ctx, nil, domain.PVZCity(request.Body.City))
	if err != nil {
		return oapi.PostPvz400JSONResponse{
			Message: "Неверный запрос",
		}, nil
	}

	return oapi.PostPvz201JSONResponse{
		City:             oapi.PVZCity(pvz.City),
		Id:               pointer.Ref(pvz.ID),
		RegistrationDate: pointer.Ref(pvz.RegisteredAt),
	}, nil
}

// GetPvz implements api.StrictServerInterface.
func (s *Server) GetPvz(
	ctx context.Context,
	request oapi.GetPvzRequestObject,
) (oapi.GetPvzResponseObject, error) {
	return oapi.GetPvz200JSONResponse{}, nil
}
