package repository

import (
	"context"
	"errors"
	"time"

	"avito_pvz/internal/domain"
)

var _ domain.ProductsRepository = (*Product)(nil)

var (
	errProduct        = errors.New("Product repository error")
	ErrProductsCreate = errors.Join(errProduct, errors.New("create failed"))
	ErrProductsDelete = errors.Join(errProduct, errors.New("delete failed"))
)

type Product struct{}

func (p *Product) Create(ctx context.Context, connection domain.Connection, product domain.Product) error {
	const query = `insert into Product
    (id, reception_id, product_type, created_at)
	values
    ($1, $2, $3, default)`

	_, err := connection.ExecContext(ctx, query, product.ID, product.ReceptionID, product.Type)
	if err != nil {
		return errors.Join(ErrPVZCreate, err)
	}

	return nil
}

func (p *Product) DeleteLast(ctx context.Context, connection domain.Connection, receptionID domain.ReceptionID) error {
	const query = ` 
	delete from products 
	where id = (select id from product 
	where reception_id = $1 order by created_at desc limit 1)`

	_, err := connection.ExecContext(ctx, query, receptionID)
	if err != nil {
		return errors.Join(ErrProductsDelete, err)
	}

	return nil
}

func (p *Product) Search(ctx context.Context, connection domain.Connection, from *time.Time, to *time.Time, page *int, limit *int) ([]domain.Product, error) {
	var products []domain.Product

	const query = `
	select * from products 
	where created_at between $1 and $2 
	order by created_at desc 
	offset $3 limit $4`

	err := connection.SelectContext(ctx, &products, query, from, to, page, limit)

	return products, err
}

func NewProduct() *Product {
	return &Product{}
}
