package grpcclient

import (
	"context"
	"ecommerce/common"
	"ecommerce/gen/category"
	"ecommerce/module/product/domain/query"
	"fmt"

	"github.com/google/uuid"
)

type categoryGRPCClient struct {
	client category.CategoryServiceClient
}

func NewCategoryGRPCClient(client category.CategoryServiceClient) *categoryGRPCClient {
	return &categoryGRPCClient{client: client}
}

func (c *categoryGRPCClient) FindCategoriesByIDs(ctx context.Context, categoryIDs []common.UUID) ([]query.CategoryDTO, error){
	// Convert uuid slice into string slice
	ids := make([]string, len(categoryIDs))

	for i:= range(ids) {
		ids[i] = categoryIDs[i].String()
	}

	// Create Find category request parameter
	req := category.FindCategoriesReq{Ids :ids}

	// Call gRPC client method
	resp, err := c.client.FindCategoriesByIDs(ctx, &req)

	if err != nil {
		return nil, fmt.Errorf("error finding categories: %w", err)
	}

	// Convert response into []query.CategoryDTO
	categories := make([]query.CategoryDTO, len(resp.Data))
	
	for i, dto := range(resp.Data) {
		categories[i] = query.CategoryDTO{
			Id : common.UUID(uuid.MustParse(dto.Id)),
			Title: dto.Title,
		}
	}
	return categories, nil
}