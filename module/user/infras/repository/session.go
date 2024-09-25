package repository

import (
	"context"
	"ecommerce/common"
	"ecommerce/module/user/domain"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	SESSION_TABLE_NAME = "user_sessions"
)

type mySQLSession struct {
	db *gorm.DB
}

func NewMysqlSession(db *gorm.DB) *mySQLSession{
	return &mySQLSession{db:db}
}

func (repo *mySQLSession) Create(ctx context.Context, data *domain.Session) error {
	// Transform session entity to dto to create new session
	dto := SessionDTO{
		Id: data.GetID(),
		UserId: data.GetUserID(),
		RefreshToken: data.GetRefreshToken(),
		AccessExpAt: data.GetAccessExpAt(),
		RefreshExpAt: data.GetRefreshExpAt(),
	}

	// Create new session
	if err := repo.db.Table(SESSION_TABLE_NAME).Create(&dto).Error; err != nil {
		return err
	}
	return nil
}

func (repo *mySQLSession) Find(ctx context.Context, sessionID string) (*domain.Session, error){
	// Find session by id
	var session SessionUpdateDTO
	if err := repo.db.Table(SESSION_TABLE_NAME).Where("id = ?", common.UUID(uuid.MustParse(sessionID))).First(&session).Error; err != nil {
		// If record is not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		
		return nil, err
	}

	sessionEntity, err := session.ToEntity()
	if err != nil {
		return nil, err
	}
	return sessionEntity, nil
}

