package productusecase

import (
	"ecommerce/module/product/productdomain"

	"golang.org/x/net/context"

	"strings"
)

// Driving side adapter of use case
type CreateProductUseCase interface {
	CreateProduct(ctx context.Context, prod *productdomain.ProductCreationDTO) error
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

func (uc CreateNewProductUseCase) CreateProduct(ctx context.Context, prod *productdomain.ProductCreationDTO) error {
	prod.Name = strings.TrimSpace(prod.Name)

	if prod.Name == "" {
		return productdomain.ErrProductNameCannotBeBlank
	}

	if err := uc.repo.CreateProduct(ctx, prod); err != nil {
		return err
	}

	return nil
}

// Repository adapter
type CreateProductRepository interface {
	CreateProduct(ctx context.Context, prod *productdomain.ProductCreationDTO) error
}
