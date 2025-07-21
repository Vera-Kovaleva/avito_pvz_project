package domain

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwordBCryptoCost  = 12
	passwordHashBytesLen = 16
)

func GenerateToken(userID UserID, userRole UserRole) (string, error) {
	return userID.String() + ":" + string(userRole), nil
}

type authenticatedUser struct {
	ID   UserID
	Role UserRole
}

func (u *authenticatedUser) GetUserID() UserID {
	return u.ID
}

func (u *authenticatedUser) GetUserRole() UserRole {
	return u.Role
}

func AuthenticateByToken(token string) (AuthenticatedUser, error) {
	partsCount := 2
	parts := strings.Split(token, ":")
	if len(parts) != partsCount {
		return nil, errors.New("invalid token format")
	}

	userIDStr := parts[0]
	roleStr := parts[1]

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	return &authenticatedUser{
		ID:   userID,
		Role: UserRole(roleStr),
	}, nil
}

func HashPassword(pasword string) (string, error) {
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(pasword), passwordBCryptoCost)
	if err != nil {
		return "", err
	}
	passwordHash := string(passwordHashBytes)
	return passwordHash, nil
}

func CompareHashAndPassword(password string, passwordHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return errors.Join(err, errors.New("wrong password"))
	}
	return nil
}
