package http

import (
	"context"
	"time"

	"avito_pvz/internal/domain"
	oapi "avito_pvz/internal/generated/oapi"
	"avito_pvz/internal/infra/pointer"
)

func (s *Server) PostPvzPvzIdCloseLastReception(
	ctx context.Context,
	request oapi.PostPvzPvzIdCloseLastReceptionRequestObject,
) (oapi.PostPvzPvzIdCloseLastReceptionResponseObject, error) {
	reception, err := s.receptions.Close(ctx, nil, domain.PVZID(request.PvzId))
	if err != nil {
		return oapi.PostPvzPvzIdCloseLastReception400JSONResponse{
			Message: "Неверный запрос или приемка уже закрыта",
		}, nil
	}

	return oapi.PostPvzPvzIdCloseLastReception200JSONResponse{
		Id:       pointer.Ref(reception.ID),
		PvzId:    reception.PVZID,
		Status:   oapi.ReceptionStatus(reception.Status),
		DateTime: time.Time(reception.CreatedAt),
	}, nil
}

func (s *Server) PostPvzPvzIdDeleteLastProduct(
	ctx context.Context,
	request oapi.PostPvzPvzIdDeleteLastProductRequestObject,
) (oapi.PostPvzPvzIdDeleteLastProductResponseObject, error) {
	err := s.receptions.DeleteLastProduct(ctx, nil, domain.PVZID(request.PvzId))

	if err != nil {
		return oapi.PostPvzPvzIdDeleteLastProduct400JSONResponse{
			Message: "Неверный запрос, нет активной приемки или нет товаров для удаления",
		}, nil
	}

	return oapi.PostPvzPvzIdDeleteLastProduct200Response{}, nil
}

// PostProducts implements api.StrictServerInterface.
func (s *Server) PostProducts(
	ctx context.Context,
	request oapi.PostProductsRequestObject,
) (oapi.PostProductsResponseObject, error) {
	product, err := s.receptions.CreateProduct(ctx, nil, domain.PVZID(request.Body.PvzId), domain.ProductType(request.Body.Type))
	if err != nil {
		return oapi.PostProducts400JSONResponse{
			Message: "Неверный запрос или нет активной приемки",
		}, nil
	}
	return oapi.PostProducts201JSONResponse{
		DateTime:    pointer.Ref(product.CreatedAt),
		Id:          pointer.Ref(product.ID),
		ReceptionId: product.ReceptionID,
		Type:        oapi.ProductType(product.Type),
	}, nil
}

// PostReceptions implements api.StrictServerInterface.
func (s *Server) PostReceptions(
	ctx context.Context,
	request oapi.PostReceptionsRequestObject,
) (oapi.PostReceptionsResponseObject, error) {
	reception, err := s.receptions.Create(ctx, nil, domain.PVZID(request.Body.PvzId))
	if err != nil {
		return oapi.PostReceptions400JSONResponse{
			Message: "Неверный запрос или есть незакрытая приемка",
		}, nil
	}
	return oapi.PostReceptions201JSONResponse{
		DateTime: reception.CreatedAt,
		Id:       pointer.Ref(reception.ID),
		PvzId:    reception.PVZID,
		Status:   oapi.ReceptionStatus(reception.Status),
	}, nil
}
