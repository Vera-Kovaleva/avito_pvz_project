package repository

import (
	"context"
	"errors"

	"avito_pvz/internal/domain"
)

var _ domain.ReceptionsRepository = (*Reception)(nil)

var (
	errReception           = errors.New("resseptions error")
	ErrCreateReception     = errors.Join(errReception, errors.New("create failed"))
	ErrFindActiveReception = errors.Join(errReception, errors.New("find active failed"))
	ErrCloseReception      = errors.Join(errReception, errors.New("close failed"))
	ErrFindByIDsReception  = errors.Join(errReception, errors.New("find by IDs failed"))
)

type Reception struct{}

func NewReceptions() *Reception {
	return &Reception{}
}

func (r *Reception) Create(
	ctx context.Context,
	connection domain.Connection,
	reception domain.Reception,
) error {
	const query = `insert into receptions
    (id, created_at, pvz_id, status)
	values
    ($1, default, $2, $3)`

	_, err := connection.ExecContext(ctx, query, reception.ID, reception.PVZID, domain.InProgress)
	if err != nil {
		return errors.Join(ErrCreateReception, err)
	}

	return nil
}

func (r *Reception) Close(
	ctx context.Context,
	connection domain.Connection,
	receptionID domain.ReceptionID,
) error {
	const query = `update receptions set status = $2 where id = $1`

	_, err := connection.ExecContext(ctx, query, receptionID, domain.Close)
	if err != nil {
		return errors.Join(ErrCloseReception, err)
	}

	return nil
}

func (r *Reception) FindActive(
	ctx context.Context,
	connection domain.Connection,
	pvzID domain.PVZID,
) (domain.Reception, error) {
	const query = `select id, created_at, pvz_id, status from receptions where pvz_id = $1 and status = 'in_progress'`

	var reception domain.Reception
	err := connection.GetContext(ctx, &reception, query, pvzID)
	if err != nil {
		return reception, errors.Join(ErrFindActiveReception, err)
	}

	return reception, nil
}

func (r *Reception) FindByIDs(
	ctx context.Context,
	connection domain.Connection,
	receptionIDs []domain.ReceptionID,
) ([]domain.Reception, error) {
	const query = `select id, created_at, pvz_id, status from receptions where id = any($1)`

	var receptions []domain.Reception
	err := connection.SelectContext(ctx, &receptions, query, receptionIDs)
	if err != nil {
		return nil, errors.Join(ErrFindByIDsReception, err)
	}

	return receptions, nil
}
