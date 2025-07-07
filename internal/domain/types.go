package domain

import (
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
		ID           UserID
		Email        string
		Role         UserRole
		PasswordHash string
		Token        string
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
