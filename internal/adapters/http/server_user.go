package http

import (
	"context"

	oapi "avito_pvz/internal/generated/oapi"
)

// PostDummyLogin implements api.StrictServerInterface.
func (s *Server) PostDummyLogin(
	ctx context.Context,
	request oapi.PostDummyLoginRequestObject,
) (oapi.PostDummyLoginResponseObject, error) {
	panic("unimplemented")
}

// PostLogin implements api.StrictServerInterface.
func (s *Server) PostLogin(
	ctx context.Context,
	request oapi.PostLoginRequestObject,
) (oapi.PostLoginResponseObject, error) {
	panic("unimplemented")
}

// PostRegister implements api.StrictServerInterface.
func (s *Server) PostRegister(
	ctx context.Context,
	request oapi.PostRegisterRequestObject,
) (oapi.PostRegisterResponseObject, error) {
	panic("unimplemented")
}
