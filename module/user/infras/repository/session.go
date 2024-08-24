package repository

import(
	"context"
	"ecommerce/module/user/domain"
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

