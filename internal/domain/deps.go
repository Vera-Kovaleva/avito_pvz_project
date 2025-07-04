package domain

import (
	"context"
	"errors"
	"io"
	"time"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrPVZNotFound       = errors.New("PVZ not found")
	ErrReceptionNotFound = errors.New("reception not found")
)

type (
	Connection interface {
		GetContext(context.Context, any, string, ...any) error
		SelectContext(context.Context, any, string, ...any) error
		ExecContext(context.Context, string, ...any) (int64, error)
	}

	ConnectionProvider interface {
		Execute(context.Context, func(context.Context, Connection) error) error
		ExecuteTx(context.Context, func(context.Context, Connection) error) error
		io.Closer
	}
)

type (
	UsersRepository interface {
		Create(context.Context, Connection, User) error
		ReadByEmail(context.Context, Connection, string) (User, error)
		Update(context.Context, Connection, User) error
		UpdateTokenByEmail(context.Context, Connection, string, string) error
	}

	PVZsRepository interface {
		Create(context.Context, Connection, PVZ) error
		FindByIDs(context.Context, Connection, []PVZID) ([]PVZ, error)
		FindAll(context.Context, Connection) ([]PVZ, error)
	}

	ReceptionsRepository interface {
		Create(context.Context, Connection, Reception) error
		FindActive(context.Context, Connection, PVZID) (Reception, error)
		FindByIDs(context.Context, Connection, []ReceptionID) ([]Reception, error)
		Close(context.Context, Connection, ReceptionID) error
	}

	ProductsRepository interface {
		Create(context.Context, Connection, Product) error
		DeleteLast(context.Context, Connection, ReceptionID) error
		Search(
			ctx context.Context,
			connection Connection,
			from, to *time.Time,
			page, limit *int,
		) ([]Product, error)
	}
)
