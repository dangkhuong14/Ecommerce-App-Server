package usecase

import (
	"ecommerce/common"
	"ecommerce/module/product/domain"

	"context"

	"strings"
)

// Driving side adapter of use case
type CreateProductUseCase interface {
	CreateProduct(ctx context.Context, prod *domain.ProductCreationDTO) error
}

func NewCreateProductUseCase(repo CreateProductRepository) CreateNewProductUseCase {
	return CreateNewProductUseCase{
		repo: repo,
	}
}

// Use case
type CreateNewProductUseCase struct {
	repo CreateProductRepository
}

func (uc CreateNewProductUseCase) CreateProduct(ctx context.Context, prod *domain.ProductCreationDTO) error {
	// business logic
	*prod.Name = strings.TrimSpace(*prod.Name)

	if *prod.Name == "" {
		return domain.ErrProductNameCannotBeBlank
	}

	// Generate product's id
	prod.ID = common.GenNewUUID()

	if err := uc.repo.CreateProduct(ctx, prod); err != nil {
		return err
	}

	return nil
}

// Repository adapter
type CreateProductRepository interface {
	CreateProduct(ctx context.Context, prod *domain.ProductCreationDTO) error
}
