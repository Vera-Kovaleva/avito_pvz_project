package repository

import (
	"context"
	"errors"

	"avito_pvz/internal/domain"
)

var _ domain.UsersRepository = (*Users)(nil)

var (
	errUsers                   = errors.New("users repository error")
	ErrUsersCreate             = errors.Join(errUsers, errors.New("create failed"))
	ErrUsersRead               = errors.Join(errUsers, errors.New("read failed"))
	ErrUsersUpdate             = errors.Join(errUsers, errors.New("update failed"))
	ErrUsersDelete             = errors.Join(errUsers, errors.New("delete failed"))
	ErrUsersUpdateTokenByEmail = errors.Join(errUsers, errors.New("update token by email failed"))
)

type Users struct{}

func NewUsers() *Users {
	return &Users{}
}

func (r Users) Create(ctx context.Context, connection domain.Connection, user domain.User) error {
	const query = `
insert into users
    (id, email, role, password_hash, token)
values
    ($1, $2, $3, $4, $5)`

	_, err := connection.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.Role,
		user.PasswordHash,
		user.Token,
	)
	if err != nil {
		return errors.Join(ErrUsersCreate, err)
	}

	return nil
}

// ReadByEmail implements domain.UsersRepository.
func (r *Users) ReadByEmail(context.Context, domain.Connection, string) (domain.User, error) {
	panic("unimplemented")
}

// Update implements domain.UsersRepository.
func (r *Users) Update(context.Context, domain.Connection, domain.User) error {
	panic("unimplemented")
}

// UpdateTokenByEmail implements domain.UsersRepository.
func (r *Users) UpdateTokenByEmail(context.Context, domain.Connection, string, string) error {
	panic("unimplemented")
}
