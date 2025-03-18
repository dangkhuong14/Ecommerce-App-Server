package query

import (
	"context"
	"ecommerce/common"

	"github.com/viettranx/service-context/core"
)

type ListProductWrapper struct {
	next listProductQuery
	categoryRepo CategoryRepository
}

func NewListProductWrapper(next listProductQuery, categoryRepo CategoryRepository) *ListProductWrapper {
	return &ListProductWrapper{
		next: next,
		categoryRepo: categoryRepo,
	}
}

func (w *ListProductWrapper) Execute(ctx context.Context, param *ListProductQueryParam) ([]ProductDTO, error) {
	// Gọi hàm wrapped
	products, err := w.next.Execute(ctx, param)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, nil // Không cần tiếp tục xử lý nếu không có sản phẩm
	}

	// Tạo map để loại bỏ trùng lặp category ID
	categoryIDSet := make(map[common.UUID]struct{}, len(products))
	for _, product := range products {
		categoryIDSet[product.CatId] = struct{}{}
	}

	// Chuyển map thành slice để truy vấn
	categoryIDs := make([]common.UUID, 0, len(categoryIDSet))
	for id := range categoryIDSet {
		categoryIDs = append(categoryIDs, id)
	}

	// Fetch category qua category service (RPC)
	categories, err := w.categoryRepo.FindCategoriesByIDs(ctx, categoryIDs)
	if err != nil {
		return nil, core.ErrInternalServerError.
			WithDebug(err.Error()).
			WithError("error fetching product categories")
	}

	// Tạo map category để ánh xạ nhanh
	categoryMap := make(map[common.UUID]*CategoryDTO, len(categories))

	for i, category := range categories {
		categoryMap[category.Id] = &categories[i]
	}

	// Gán category vào từng product
	for i := range products {
		if category, ok := categoryMap[products[i].CatId]; ok {
			products[i].Category = category
		}
	}

	return products, nil

}

type CategoryRepository interface {
	FindCategoriesByIDs(ctx context.Context, categoryIDs []common.UUID) ([]CategoryDTO, error)
}