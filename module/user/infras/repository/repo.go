package repository

import (
	"context"
	"ecommerce/common"
	"ecommerce/module/user/domain"
	"errors"

	"gorm.io/gorm"
)

const(
	TbName = "users"
)

type mysqlUser struct {
	db *gorm.DB
}

func NewMysqlUser(db *gorm.DB) mysqlUser {
	return mysqlUser{db: db}
}

func (repo mysqlUser) Create(ctx context.Context, data *domain.User) error {
	// Transform entity into dto
	dto := UserDTO{
		Id: data.GetID(),
		FirstName: data.GetFirstName(),
		LastName: data.GetLastName(),
		Password: data.GetPassword(),
		Salt: data.GetSalt(),
		Role: data.GetRole().String(),
	}

	if err := repo.db.Table(TbName).Create(dto).Error; err != nil {
		return err
	}
	return nil
}

func (repo mysqlUser) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var dto UserDTO

	if err := repo.db.Table(TbName).Where("email = ?", email).First(&dto).Error; err != nil {
		// If record is not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}

		return nil, err
	}
	// Transform dto into entity
	return dto.ToEntity()
}
