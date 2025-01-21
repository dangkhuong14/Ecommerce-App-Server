package repository

import (
	"context"
	"ecommerce/common"
	"ecommerce/module/user/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
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
		Id:        data.GetID(),
		FirstName: data.GetFirstName(),
		LastName:  data.GetLastName(),
		Email:     data.GetEmail(),
		Password:  data.GetPassword(),
		Salt:      data.GetSalt(),
		Role:      data.GetRole().String(),
		Status:    data.GetStatus(),
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

func (repo mysqlUser) Find(ctx context.Context, userID string) (*domain.User, error) {
	// Find user by id
	var user UserDTO
	if err := repo.db.Table(TbName).Where("id = ?", common.UUID(uuid.MustParse(userID))).First(&user).Error; err != nil {
		// If record is not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}

		return nil, err
	}

	userEntity, err := user.ToEntity()
	if err != nil {
		return nil, err
	}
	return userEntity, nil
}

func (repo mysqlUser) UpdateAvatar(ctx context.Context, user *domain.User, avatar *domain.Avatar) error {
	// Transform entity into dto
	dto := toAvatarDTO(avatar)

	result := repo.db.Table(TbName).
		Where("id = ?", user.GetID()).
		Update("avatar", dto)
	
	// if err is gorm record not found error
	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	// check if no rows were affected
	if result.RowsAffected == 0 {
		return errors.WithStack(common.ErrRecordNotFound)
	}

	return nil
}
