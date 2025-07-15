package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type (
	UserID          = uuid.UUID
	UserRole        string
	PVZID           = uuid.UUID
	PVZCity         string
	ReceptionID     = uuid.UUID
	ReceptionStatus string
	ProductID       = uuid.UUID
	ProductType     string

	User struct {
		ID           UserID   `db:"id"`
		Email        string   `db:"email"`
		Role         UserRole `db:"role"`
		PasswordHash string   `db:"password_hash"`
		Token        string   `db:"token"`
	}

	PVZ struct {
		ID           PVZID
		City         PVZCity
		RegisteredAt time.Time
	}

	Reception struct {
		ID        ReceptionID     `db:"id"`
		PVZID     PVZID           `db:"pvz_id"`
		Status    ReceptionStatus `db:"status"`
		CreatedAt time.Time       `db:"created_at"`
	}

	Product struct {
		ID          ProductID   `db:"id"`
		ReceptionID ReceptionID `db:"reception_id"`
		Type        ProductType `db:"type"`
		CreatedAt   time.Time   `db:"created_at"`
	}

	PVZReceptionsProducts struct {
		PVZ        PVZ
		Receptions []ReceptionsProducts
	}

	ReceptionsProducts struct {
		Reception Reception
		Products  []Product
	}

	AuthenticatedUser interface {
		GetUserID() UserID
		GetUserRole() UserRole
	}
)

const (
	Employee  UserRole = "employee"
	Moderator UserRole = "moderator"
)

const (
	Msk PVZCity = "Москва"
	SPb PVZCity = "Санкт-Петербург"
	Kzn PVZCity = "Казань"
)

const (
	InProgress ReceptionStatus = "in_progress"
	Close      ReceptionStatus = "close"
)

const (
	Electronics ProductType = "электроника"
	Clothes     ProductType = "одежда"
	Shoes       ProductType = "обувь"
)

type (
	UsersInterface interface {
		Create(context.Context, string, string, UserRole) (User, error)
		FindTokenByEmailAndPassword(context.Context, string, string) (string, error)
		LoginByToken(context.Context, string) (AuthenticatedUser, error)
	}

	PVZsInterface interface {
		Create(context.Context, AuthenticatedUser, PVZCity) (PVZ, error)
		FindPVZReceptionProducts(
			context.Context,
			AuthenticatedUser,
			*time.Time,
			*time.Time,
			*int,
			*int,
		) ([]PVZReceptionsProducts, error)
		FindAll(context.Context) ([]PVZ, error)
	}

	ReceptionsInterface interface {
		Create(context.Context, AuthenticatedUser, PVZID) (Reception, error)
		CreateProduct(context.Context, AuthenticatedUser, PVZID, ProductType) (Product, error)
		DeleteLastProduct(context.Context, AuthenticatedUser, PVZID) error
		Close(context.Context, AuthenticatedUser, PVZID) (Reception, error)
	}
)
