package image

import (
	"context"
	"ecommerce/common"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	ImageTbName = "images"
)

type imageRepo struct {
	db *gorm.DB
}

func NewImageRepo(db *gorm.DB) *imageRepo {
	return &imageRepo{db: db}
}

func (repo *imageRepo) Create(ctx context.Context, entity *Image) error {
	if err := repo.db.Table(ImageTbName).Create(&entity).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (repo *imageRepo) Find(ctx context.Context, id common.UUID) (*common.Image, error) {
	imageDTO := common.Image{}
	if err := repo.db.Table(ImageTbName).Where("id = ?", id).First(&imageDTO).Error; err != nil {

		// if err is gorm record not found error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(common.ErrRecordNotFound)
		}
		return nil, errors.WithStack(err)
	}
	return &imageDTO, nil
}

func (repo *imageRepo) SetImageStatusActivated(ctx context.Context, id common.UUID) error {
	result := repo.db.Table(ImageTbName).Where("id = ?", id).Updates(Image{Status: StatusActivated})

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
