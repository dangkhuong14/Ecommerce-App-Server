package common

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        UUID      `gorm:"column:id" json:"id"`
	Status    string    `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func GenNewBaseModel() BaseModel {
	now := time.Now().UTC()

	return BaseModel{
		ID:        GenNewUUID(),
		Status:    "activated",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func GenNewUUID() UUID {
	newUUID, _ := uuid.NewV7()
	return UUID(newUUID)
}

func ParseUUID(s string) (UUID, error) {
	newUUID, err := uuid.Parse(s)
	if err != nil {
		return UUID{}, err
	}
	
	return UUID(newUUID), nil
	
}
