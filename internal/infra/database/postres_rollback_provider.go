package database

import (
	"context"
	"errors"

	"avito_pvz/internal/domain"
)

var _ domain.ConnectionProvider = (*PostgresRollbackProvider)(nil)

type PostgresRollbackProvider struct {
	provider domain.ConnectionProvider
}

func NewPostgresRollbackProvider(provider domain.ConnectionProvider) *PostgresRollbackProvider {
	return &PostgresRollbackProvider{
		provider: provider,
	}
}

func (p PostgresRollbackProvider) Execute(
	ctx context.Context,
	receiver func(context.Context, domain.Connection) error,
) error {
	return p.ExecuteTx(ctx, receiver)
}

func (p PostgresRollbackProvider) ExecuteTx(
	ctx context.Context,
	receiver func(context.Context, domain.Connection) error,
) error {
	errRollback := errors.New(
		"this provider always rollbacks transactions and used for testing purposes",
	)
	err := p.provider.ExecuteTx(ctx, func(ctx context.Context, connection domain.Connection) error {
		if err := receiver(ctx, connection); err != nil {
			return err
		}

		return errRollback
	})
	if errors.Is(err, errRollback) {
		err = nil
	}

	return err
}

func (p PostgresRollbackProvider) Close() error {
	return p.provider.Close()
}
