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
		ID        ReceptionID
		PVZID     PVZID
		Status    ReceptionStatus
		CreatedAt time.Time
	}

	Product struct {
		ID          ProductID
		ReceptionID ReceptionID
		Type        ProductType
		CreatedAt   time.Time
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
