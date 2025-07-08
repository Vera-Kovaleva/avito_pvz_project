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
	ErrUsersReadByEmail        = errors.Join(errUsers, errors.New("read by email failed"))
	ErrUsersUpdate             = errors.Join(errUsers, errors.New("update failed"))
	ErrUsersUpdateTokenByEmail = errors.Join(errUsers, errors.New("update token by email failed"))
)

type Users struct{}

func NewUsers() *Users {
	return &Users{}
}

func (r *Users) Create(ctx context.Context, connection domain.Connection, user domain.User) error {
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

func (r *Users) ReadByEmail(ctx context.Context, connection domain.Connection, email string) (domain.User, error) {
	const query = `select id, email, role, password_hash, token from users where email = $1`

	var user domain.User
	err := connection.GetContext(ctx, &user, query, email)
	if err != nil {
		return user, errors.Join(ErrUsersReadByEmail, err)
	}

	return user, nil
}

func (r *Users) Update(ctx context.Context, connection domain.Connection, user domain.User) error {
	const query = `update users set email = $2, role =$3, password_hash = $4, token = $5 where id = $1`

	_, err := connection.ExecContext(ctx, query, user.ID, user.Email, user.Role, user.PasswordHash, user.Token)
	if err != nil {
		return errors.Join(ErrUsersUpdate, err)
	}

	return nil
}

func (r *Users) UpdateTokenByEmail(ctx context.Context, connection domain.Connection, email string, token string) error {
	const query = `update users set token = $1 where email = $2`

	_, err := connection.ExecContext(ctx, query, token, email)
	if err != nil {
		return errors.Join(ErrUsersUpdateTokenByEmail, err)
	}

	return nil
}
