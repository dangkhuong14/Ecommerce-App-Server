package usecase

import (
	"ecommerce/module/product/domain"

	"golang.org/x/net/context"

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
	prod.Name = strings.TrimSpace(prod.Name)

	if prod.Name == "" {
		return domain.ErrProductNameCannotBeBlank
	}

	if err := uc.repo.CreateProduct(ctx, prod); err != nil {
		return err
	}

	return nil
}

// Repository adapter
type CreateProductRepository interface {
	CreateProduct(ctx context.Context, prod *domain.ProductCreationDTO) error
}
