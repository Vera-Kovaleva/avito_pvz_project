package repository

import (
	"context"
	"errors"

	"avito_pvz/internal/domain"
)

var _ domain.PVZsRepository = (*PVZ)(nil)

var (
	errPVZ          = errors.New("pvzs error")
	ErrPVZCreate    = errors.Join(errPVZ, errors.New("create failed"))
	ErrPVZFindByIDs = errors.Join(errPVZ, errors.New("find by IDs failed"))
	ErrPVZFindAll   = errors.Join(errPVZ, errors.New("find all failed"))
)

type PVZ struct{}

func NewPVZ() *PVZ {
	return &PVZ{}
}

func (p *PVZ) Create(ctx context.Context, connection domain.Connection, pvz domain.PVZ) error {
	const query = `insert into pvz
    (id, city, registered_at)
	values
    ($1, $2, $3)`

	_, err := connection.ExecContext(ctx, query, pvz.ID, pvz.City, pvz.RegisteredAt)
	if err != nil {
		return errors.Join(ErrPVZCreate, err)
	}

	return nil
}

func (p *PVZ) FindAll(ctx context.Context, connection domain.Connection) ([]domain.PVZ, error) {
	const query = `select id, city, registered_at from pvz`

	var pvzs []domain.PVZ
	err := connection.SelectContext(ctx, &pvzs, query)
	if err != nil {
		return nil, errors.Join(ErrPVZFindAll, err)
	}

	return pvzs, nil
}

func (p *PVZ) FindByIDs(
	ctx context.Context,
	connection domain.Connection,
	pvzIDs []domain.PVZID,
) ([]domain.PVZ, error) {
	const query = `select id, city, registered_at from pvz where id = any($1)`

	var pvzs []domain.PVZ
	err := connection.SelectContext(ctx, &pvzs, query, pvzIDs)
	if err != nil {
		return nil, errors.Join(ErrPVZFindByIDs, err)
	}

	return pvzs, nil
}
