package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"avito_pvz/internal/domain"
)

var _ domain.ProductsRepository = (*Product)(nil)

var (
	errProduct       = errors.New("products repository error")
	ErrCreateProduct = errors.Join(errProduct, errors.New("create failed"))
	ErrDeleteProduct = errors.Join(errProduct, errors.New("delete failed"))
	ErrSearchProduct = errors.Join(errProduct, errors.New("search failed"))
)

type Product struct{}

func NewProduct() *Product {
	return &Product{}
}

func (p *Product) Create(
	ctx context.Context,
	connection domain.Connection,
	product domain.Product,
) error {
	const query = `insert into products
    (id, reception_id, type, created_at)
	values
    ($1, $2, $3, $4)`

	_, err := connection.ExecContext(ctx, query, product.ID, product.ReceptionID, product.Type, product.CreatedAt)
	if err != nil {
		return errors.Join(ErrCreateProduct, err)
	}

	return nil
}

func (p *Product) DeleteLast(
	ctx context.Context,
	connection domain.Connection,
	receptionID domain.ReceptionID,
) error {
	const query = ` 
	delete from products 
	where id = (select id from products
	where reception_id = $1 order by created_at desc limit 1)`

	_, err := connection.ExecContext(ctx, query, receptionID)
	if err != nil {
		return errors.Join(ErrDeleteProduct, err)
	}

	return nil
}

func (p *Product) Search(
	ctx context.Context,
	connection domain.Connection,
	from *time.Time,
	to *time.Time,
	page *int,
	limit *int,
) ([]domain.Product, error) {

	if from != nil && to != nil && to.Before(*from) {
		return nil, errors.Join(ErrSearchProduct, errors.New("from must be less than to"))
	} else if *page < 0 {
		return nil, errors.Join(ErrSearchProduct, errors.New("invalud page"))
	} else if *limit <= 0 {
		return nil, errors.Join(ErrSearchProduct, errors.New("invalid limit"))
	} else if page != nil && limit == nil {
		return nil, errors.Join(ErrSearchProduct, errors.New("page without limit"))
	}

	const baseQuery = `select id, reception_id, type, created_at from products`
	var args []interface{}
	var conditions = ""

	if from != nil && to != nil {
		conditions = "created_at between $1 and $2"
		args = append(args, *from, *to)

	} else if from != nil {
		conditions = "created_at > $1"
		args = append(args, *from)

	} else if to != nil {
		conditions = "created_at < $1"
		args = append(args, *to)
	}

	query := baseQuery
	if len(conditions) > 0 {
		query += " where " + conditions
	}
	query += " order by created_at desc"

	argPos := len(args) + 1
	if page != nil && limit != nil {
		query += " offset $" + strconv.Itoa(argPos) + " limit $" + strconv.Itoa(argPos+1)
		args = append(args, (*page)*(*limit), *limit)
	} else if limit != nil {
		query += " limit $" + strconv.Itoa(argPos)
		args = append(args, *limit)
	}

	var products []domain.Product
	err := connection.SelectContext(ctx, &products, query, args...)
	if err != nil {
		return nil, errors.Join(ErrSearchProduct, err)
	}
	return products, nil

}
