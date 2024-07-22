package controller

import (
	"context"
	"ecommerce/module/product/domain"
)

type CreateProductUseCase interface {
	CreateProduct(ctx context.Context, prod *domain.ProductCreationDTO) error
}

type APIController struct {
	createUseCase CreateProductUseCase
}

func NewAPIController(createUseCase CreateProductUseCase) APIController {
	return APIController{
		createUseCase: createUseCase,
	}
}
