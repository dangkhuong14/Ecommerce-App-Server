package image

import (
	"context"

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
