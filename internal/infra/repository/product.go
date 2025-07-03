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
)

type Product struct{}

// Create implements domain.ProductsRepository.
func (p *Product) Create(context.Context, domain.Connection, domain.Product) error {
	panic("unimplemented")
}

// DeleteLast implements domain.ProductsRepository.
func (p *Product) DeleteLast(context.Context, domain.Connection, domain.ReceptionID) error {
	panic("unimplemented")
}

// Search implements domain.ProductsRepository.
func (p *Product) Search(ctx context.Context, connection domain.Connection, from *time.Time, to *time.Time, page *int, limit *int) ([]domain.Product, error) {
	panic("unimplemented")
}

func NewProduct() *Product {
	return &Product{}
}
