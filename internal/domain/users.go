package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var _ UsersInterface = (*UserService)(nil)

var (
	errUser                        = errors.New("user service error")
	ErrAvitoServiceCreateUser      = errors.Join(errUser, errors.New("create user failed"))
	ErrFindTokenByEmailAndPassword = errors.Join(errUser, errors.New("find token failed"))
	ErrInvalidPasswordUser         = errors.Join(errUser, errors.New("invalid password"))
	ErrInvalidToken                = errors.Join(errUser, errors.New("invalid token"))
)

type UserService struct {
	provider               ConnectionProvider
	userRepo               UsersRepository
	hashPassword           func(string) (string, error)
	compareHashAndPassword func(string, string) error
	generateToken          func(UserID, UserRole) (string, error)
	authenticateByToken    func(string) (AuthenticatedUser, error)
}

func NewUserService(
	provider ConnectionProvider,
	userRepo UsersRepository,
	hashPassword func(string) (string, error),
	compareHashAndPassword func(string, string) error,
	generateToken func(UserID, UserRole) (string, error),
	authenticateByToken func(string) (AuthenticatedUser, error),
) *UserService {
	return &UserService{
		provider:               provider,
		userRepo:               userRepo,
		hashPassword:           hashPassword,
		compareHashAndPassword: compareHashAndPassword,
		generateToken:          generateToken,
		authenticateByToken:    authenticateByToken,
	}
}

// откуда generate token ??
// login by token все??

func (s *UserService) Create(
	ctx context.Context,
	email string,
	password string,
	userRole UserRole,
) (User, error) {
	var user User
	var err error

	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return user, errors.Join(ErrAvitoServiceCreateUser, err)
	}

	userID := UserID(uuid.New())

	token, err := s.generateToken(userID, userRole)
	if err != nil {
		return user, errors.Join(ErrAvitoServiceCreateUser, err)
	}

	err = s.provider.ExecuteTx(ctx, func(ctx context.Context, connection Connection) error {
		user = User{
			ID:           userID,
			Email:        email,
			Role:         userRole,
			PasswordHash: hashedPassword,
			Token:        token,
		}
		return s.userRepo.Create(ctx, connection, user)
	})
	err = s.provider.Execute(ctx, func(ctx context.Context, connection Connection) error {
		user, err = s.userRepo.ReadByEmail(ctx, connection, email)
		return err
	})
	if err != nil {
		return user, errors.Join(ErrAvitoServiceCreateUser, err)
	}

	return user, nil
}

func (s *UserService) FindTokenByEmailAndPassword(
	ctx context.Context,
	email string,
	passwordHash string,
) (string, error) {
	var user User
	err := s.provider.Execute(ctx, func(ctx context.Context, connection Connection) error {
		var err error
		user, err = s.userRepo.ReadByEmail(ctx, connection, email)

		return err
	})
	if err != nil {
		return "", errors.Join(ErrFindTokenByEmailAndPassword, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passwordHash)); err != nil {
		return "", errors.Join(ErrInvalidPasswordUser, err)
	}

	return user.Token, nil
}

func (s *UserService) LoginByToken(ctx context.Context, token string) (AuthenticatedUser, error) {
	authUser, err := s.authenticateByToken(token)
	if err != nil {
		return nil, errors.Join(ErrInvalidToken, err)
	}

	return authUser, err
}
