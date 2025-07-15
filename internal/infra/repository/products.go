package repository

import (
	"context"
	"errors"
	"strconv"
	"strings"
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
	}
	if page != nil && *page < 0 {
		return nil, errors.Join(ErrSearchProduct, errors.New("invalid page"))
	}
	if limit != nil && *limit <= 0 {
		return nil, errors.Join(ErrSearchProduct, errors.New("invalid limit"))
	}
	if page != nil && limit == nil {
		return nil, errors.Join(ErrSearchProduct, errors.New("page without limit"))
	}

	var conditions []string
	var limits string
	var args []any

	arg := func(v any) string {
		args = append(args, v)

		return "$" + strconv.Itoa(len(args))
	}

	query := `select id, reception_id, type, created_at from products`

	if from != nil {
		conditions = append(conditions, arg(*from)+" <= created_at")
	}
	if to != nil {
		conditions = append(conditions, "created_at <= "+arg(*to))
	}
	if limit != nil {
		if page != nil {
			limits += " offset " + arg((*page)*(*limit))
		}

		limits += " limit " + arg(*limit)
	}

	if 0 < len(conditions) {
		query += " where " + strings.Join(conditions, " and ")
	}
	query += " order by created_at"
	if limits != "" {
		query += limits
	}

	var products []domain.Product
	err := connection.SelectContext(ctx, &products, query, args...)
	if err != nil {
		return nil, errors.Join(ErrSearchProduct, err)
	}
	return products, nil

}
