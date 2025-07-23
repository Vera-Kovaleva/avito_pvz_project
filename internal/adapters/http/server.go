package http

import (
	"context"

	"avito_pvz/internal/domain"
	oapi "avito_pvz/internal/generated/oapi"
)

type Server struct {
	pvzs       domain.PVZsInterface
	receptions domain.ReceptionsInterface
	users      domain.UsersInterface
}

var _ oapi.StrictServerInterface = (*Server)(nil)

func NewServer(
	pvzs domain.PVZsInterface,
	receptions domain.ReceptionsInterface,
	users domain.UsersInterface,
) *Server {
	return &Server{
		pvzs:       pvzs,
		receptions: receptions,
		users:      users,
	}
}

func (s *Server) GetCurrentUserFromCtx(ctx context.Context) domain.AuthenticatedUser {
	authUser, ok := ctx.Value(domain.CtxCurUserKey).(domain.AuthenticatedUser)

	if !ok {
		return nil
	}

	return authUser
}
