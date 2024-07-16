package mysql

import (
	"context"
	"ecommerce/module/product/domain"
)

func (repo MysqlRepository) CreateProduct(ctx context.Context, prod *domain.ProductCreationDTO) error {
	if err := repo.db.Create(&prod).Error; err != nil {
		return err
	}
	return nil
} 