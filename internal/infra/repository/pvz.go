package repository

import (
	"context"
	"errors"

	"avito_pvz/internal/domain"
)

var _ domain.PVZsRepository = (*PVZ)(nil)

var (
	errPVZ       = errors.New("pvzs error")
	ErrCreatePVZ = errors.Join(errPVZ, errors.New("create failed"))
	ErrFindByIDs = errors.Join(errPVZ, errors.New("find by IDs failed"))
	ErrFindAll   = errors.Join(errPVZ, errors.New("find all failed"))
)

type PVZ struct{}

func NewPVZ() *PVZ {
	return &PVZ{}
}

func (p *PVZ) Create(ctx context.Context, connection domain.Connection, pvz domain.PVZ) error {
	const query = `insert into pvz
    (id, city, registered_at)
	values
    ($1, $2, default)`

	_, err := connection.ExecContext(ctx, query, pvz.ID, pvz.City)
	if err != nil {
		return errors.Join(ErrCreatePVZ, err)
	}

	return nil
}

func (p *PVZ) FindAll(ctx context.Context, connection domain.Connection) ([]domain.PVZ, error) {
	const query = `select id, city, registered_at from pvz`

	var pvzs []domain.PVZ
	err := connection.SelectContext(ctx, &pvzs, query)
	if err != nil {
		return nil, errors.Join(ErrFindAll, err)
	}

	return pvzs, nil
}

func (p *PVZ) FindByIDs(ctx context.Context, connection domain.Connection, pvzIDs []domain.PVZID) ([]domain.PVZ, error) {
	const query = `select id, city, registered_at from pvz where id = any($1)`

	var pvzs []domain.PVZ
	err := connection.SelectContext(ctx, &pvzs, query, pvzIDs)
	if err != nil {
		return nil, errors.Join(ErrFindByIDs, err)
	}

	return pvzs, nil
}
