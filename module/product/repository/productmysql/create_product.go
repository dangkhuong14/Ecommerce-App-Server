package productmysql

import (
	"context"
	"ecommerce/module/product/productdomain"
)

func (repo MysqlRepository) CreateProduct(ctx context.Context, prod *productdomain.ProductCreationDTO) error {
	if err := repo.db.Create(&prod).Error; err != nil {
		return err
	}
	return nil
} 