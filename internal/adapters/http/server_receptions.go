package http

import (
	"context"

	"avito_pvz/internal/domain"
	oapi "avito_pvz/internal/generated/oapi"
	"avito_pvz/internal/infra/pointer"
)

//nolint:revive,staticcheck // Method name must comply with generated code naming requirements.
func (s *Server) PostPvzPvzIdCloseLastReception(
	ctx context.Context,
	request oapi.PostPvzPvzIdCloseLastReceptionRequestObject,
) (oapi.PostPvzPvzIdCloseLastReceptionResponseObject, error) {
	reception, err := s.receptions.Close(ctx, s.GetCurrentUserFromCtx(ctx), request.PvzId)
	if err == domain.ErrNotAuthorized {
		//nolint:nilerr // generated code expects error in response.
		return oapi.PostPvzPvzIdCloseLastReception403JSONResponse{
			Message: "Доступ запрещен",
		}, nil
	}

	if err != nil {
		//nolint:nilerr // generated code expects error in response.
		return oapi.PostPvzPvzIdCloseLastReception400JSONResponse{
			Message: "Неверный запрос или приемка уже закрыта",
		}, nil
	}

	return oapi.PostPvzPvzIdCloseLastReception200JSONResponse{
		Id:       pointer.Ref(reception.ID),
		PvzId:    reception.PVZID,
		Status:   oapi.ReceptionStatus(reception.Status),
		DateTime: reception.CreatedAt,
	}, nil
}

//nolint:revive,staticcheck // Method name must comply with generated code naming requirements.
func (s *Server) PostPvzPvzIdDeleteLastProduct(
	ctx context.Context,
	request oapi.PostPvzPvzIdDeleteLastProductRequestObject,
) (oapi.PostPvzPvzIdDeleteLastProductResponseObject, error) {
	err := s.receptions.DeleteLastProduct(ctx, s.GetCurrentUserFromCtx(ctx), request.PvzId)
	if err == domain.ErrNotAuthorized {
		//nolint:nilerr // generated code expects error in response.
		return oapi.PostPvzPvzIdDeleteLastProduct403JSONResponse{
			Message: "Доступ запрещен",
		}, nil
	}

	if err != nil {
		//nolint:nilerr // generated code expects error in response.
		return oapi.PostPvzPvzIdDeleteLastProduct400JSONResponse{
			Message: "Неверный запрос, нет активной приемки или нет товаров для удаления",
		}, nil
	}

	return oapi.PostPvzPvzIdDeleteLastProduct200Response{}, nil
}

func (s *Server) PostProducts(
	ctx context.Context,
	request oapi.PostProductsRequestObject,
) (oapi.PostProductsResponseObject, error) {
	product, err := s.receptions.CreateProduct(
		ctx,
		s.GetCurrentUserFromCtx(ctx),
		request.Body.PvzId,
		domain.ProductType(request.Body.Type),
	)
	if err == domain.ErrNotAuthorized {
		//nolint:nilerr // generated code expects error in response.
		return oapi.PostProducts403JSONResponse{
			Message: "Доступ запрещен",
		}, nil
	}

	if err != nil {
		//nolint:nilerr // generated code expects error in response.
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

func (s *Server) PostReceptions(
	ctx context.Context,
	request oapi.PostReceptionsRequestObject,
) (oapi.PostReceptionsResponseObject, error) {
	reception, err := s.receptions.Create(ctx, s.GetCurrentUserFromCtx(ctx), request.Body.PvzId)

	if err == domain.ErrNotAuthorized {
		//nolint:nilerr // generated code expects error in response.
		return oapi.PostReceptions403JSONResponse{
			Message: "Доступ запрещен",
		}, nil
	}

	if err != nil {
		//nolint:nilerr // generated code expects error in response.
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
