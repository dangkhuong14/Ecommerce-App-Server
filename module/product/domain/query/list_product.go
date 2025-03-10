package query

import (
	"context"
	"ecommerce/common"

	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
)

// Consider query as usecase

type ProductDTO struct {
	common.BaseModel
	CatId    common.UUID  `gorm:"column:category_id" json:"category_id"`
	Name     string       `gorm:"column:name" json:"name"`
	Type     string       `gorm:"column:type" json:"type"`
	Category *CategoryDTO `gorm:"foreignKey:CatId" json:"category,omitempty"`
}

type CategoryDTO struct {
	Id    common.UUID `gorm:"column:id" json:"id"`
	Title string      `gorm:"column:title" json:"title"`
}

func (CategoryDTO) TableName() string {return "categories"}

type ListProductFilterParam struct {
	CateID string `form:"cate_id"`
}

type ListProductQueryParam struct {
	common.Paging
	ListProductFilterParam
}

type listProductQuery struct {
	sctx sctx.ServiceContext
}

func NewListProductQuery(sctx sctx.ServiceContext) listProductQuery {
	return listProductQuery{sctx: sctx}
}

func (q *listProductQuery) Execute(ctx context.Context, param *ListProductQueryParam) ([]ProductDTO, error) {
	var products []ProductDTO

	// Filter by Category ID in query string
	db := q.sctx.MustGet(common.KeyGormComponent).(common.GormCompContext).GetDB()
	db = db.Table("products")

	if param.CateID != "" {
		cateUUID, err := common.ParseUUID(param.CateID)
		if err != nil {
			return nil, core.ErrBadRequest.WithDebug(err.Error()).WithError("invalid category id")
		}
		db = db.Where("category_id = ?", cateUUID)
	}

	// Paging

	// Đếm tổng số sản phẩm theo điều kiện filter đã áp dụng
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, core.ErrInternalServerError.WithDebug(err.Error())
	}
	param.Total = int(count)

	// Kiem tra cac tham so offset, limit
	param.Process()

	db = db.Preload("Category")

	// Tính offset dựa trên trang và limit
	offset := param.Limit * (param.Page - 1)

	// Thực hiện truy vấn lấy dữ liệu có phân trang
	if err := db.Order("id desc").Offset(offset).Limit(param.Limit).Find(&products).Error; err != nil {
		return nil, core.ErrBadRequest.WithError("cannot list product").WithDebug(err.Error())
	}

	return products, nil
}
