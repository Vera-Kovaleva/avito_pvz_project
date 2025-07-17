package http

import (
	"context"

	"avito_pvz/internal/domain"
	oapi "avito_pvz/internal/generated/oapi"
	"avito_pvz/internal/infra/pointer"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
)

func (s *Server) PostDummyLogin(
	ctx context.Context,
	request oapi.PostDummyLoginRequestObject,
) (oapi.PostDummyLoginResponseObject, error) {
	return oapi.PostDummyLogin200JSONResponse(uuid.New().String() + ":" + string(request.Body.Role)), nil
}

func (s *Server) PostLogin(
	ctx context.Context,
	request oapi.PostLoginRequestObject,
) (oapi.PostLoginResponseObject, error) {
	token, err := s.users.FindTokenByEmailAndPassword(ctx, request.Body.Password, string(request.Body.Email))
	if err != nil {
		return oapi.PostLogin401JSONResponse{
			Message: "Неверные учетные данные",
		}, nil
	}
	_, err = s.users.LoginByToken(ctx, token)
	if err != nil {
		return oapi.PostLogin401JSONResponse{
			Message: "Неверные учетные данные",
		}, nil
	}
	return oapi.PostLogin200JSONResponse(token), nil
}

func (s *Server) PostRegister(
	ctx context.Context,
	request oapi.PostRegisterRequestObject,
) (oapi.PostRegisterResponseObject, error) {
	user, err := s.users.Create(ctx, request.Body.Password, string(request.Body.Email), domain.UserRole(request.Body.Role))

	if err != nil {
		return oapi.PostRegister400JSONResponse{
			Message: "Неверный запрос",
		}, nil
	}

	return oapi.PostRegister201JSONResponse{
		Email: types.Email(user.Email),
		Id:    pointer.Ref(user.ID),
		Role:  oapi.UserRole(user.Role),
	}, nil
}
