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
		//nolint:nilerr // generated code expects error in response.
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

func (s *Server) GetPvz(
	ctx context.Context,
	request oapi.GetPvzRequestObject,
) (oapi.GetPvzResponseObject, error) {
	all, err := s.pvzs.FindPVZReceptionProducts(
		ctx,
		nil,
		request.Params.StartDate,
		request.Params.EndDate,
		request.Params.Page,
		request.Params.Limit,
	)
	if err != nil {
		return nil, err
	}

	type RespReception = struct {
		Products  *[]oapi.Product `json:"products,omitempty"`
		Reception *oapi.Reception `json:"reception,omitempty"`
	}

	type RespItem = struct {
		Pvz        *oapi.PVZ        `json:"pvz,omitempty"`
		Receptions *[]RespReception `json:"receptions,omitempty"`
	}

	var response oapi.GetPvz200JSONResponse

	for _, pvzData := range all {
		pvz := &oapi.PVZ{
			Id:               pointer.Ref(pvzData.PVZ.ID),
			City:             oapi.PVZCity(pvzData.PVZ.City),
			RegistrationDate: pointer.Ref(pvzData.PVZ.RegisteredAt),
		}

		var receptions []RespReception
		for _, receptionData := range pvzData.Receptions {
			reception := &oapi.Reception{
				DateTime: receptionData.Reception.CreatedAt,
				Id:       pointer.Ref(receptionData.Reception.ID),
				PvzId:    receptionData.Reception.PVZID,
				Status:   oapi.ReceptionStatus(receptionData.Reception.Status),
			}

			var products []oapi.Product
			for _, productData := range receptionData.Products {
				products = append(products, oapi.Product{
					DateTime:    pointer.Ref(productData.CreatedAt),
					Id:          pointer.Ref(productData.ID),
					ReceptionId: productData.ReceptionID,
					Type:        oapi.ProductType(productData.Type),
				})
			}

			receptions = append(receptions, RespReception{
				Reception: reception,
				Products:  pointer.Ref(products),
			})
		}

		response = append(response, RespItem{
			Pvz:        pvz,
			Receptions: pointer.Ref(receptions),
		})
	}

	return response, nil
}
