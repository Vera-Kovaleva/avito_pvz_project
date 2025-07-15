package http

import (
	"avito_pvz/internal/domain"
	oapi "avito_pvz/internal/generated/oapi"
)

type Server struct {
	pvzs domain.PVZsInterface
}

var _ oapi.StrictServerInterface = (*Server)(nil)
