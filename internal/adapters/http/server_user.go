package http

import (
	"context"
	"strings"

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
	generatedID := uuid.New()
	password := generatedID.String()
	email := strings.ReplaceAll(password, "-", "") + "@email.foo"
	userRole := domain.UserRole(request.Body.Role)

	user, err := s.users.Create(ctx, email, password, userRole)
	if err != nil {
		//nolint:nilerr // generated code expects error in response.
		return oapi.PostDummyLogin400JSONResponse{
			Message: "Неверный запрос",
		}, nil
	}

	return oapi.PostDummyLogin200JSONResponse(
		user.Token,
	), nil
}

func (s *Server) PostLogin(
	ctx context.Context,
	request oapi.PostLoginRequestObject,
) (oapi.PostLoginResponseObject, error) {
	token, err := s.users.FindTokenByEmailAndPassword(
		ctx,
		request.Body.Password,
		string(request.Body.Email),
	)
	if err != nil {
		//nolint:nilerr // generated code expects error in response.
		return oapi.PostLogin401JSONResponse{
			Message: "Неверные учетные данные",
		}, nil
	}
	_, err = s.users.LoginByToken(ctx, token)
	if err != nil {
		//nolint:nilerr // generated code expects error in response.
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
	user, err := s.users.Create(
		ctx,
		request.Body.Password,
		string(request.Body.Email),
		domain.UserRole(request.Body.Role),
	)
	if err != nil {
		//nolint:nilerr // generated code expects error in response.
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
