package query

import (
	"context"
	"ecommerce/common"

	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
)

type CategoryDTO struct {
	Id    common.UUID `gorm:"column:id" json:"id"`
	Title string      `gorm:"column:title" json:"title"`
}

func (CategoryDTO) TableName() string { return "categories" }

type findCategoriesByIDsQuery struct {
	sctx sctx.ServiceContext
}

func NewFindCategoriesByIDsQuery(sctx sctx.ServiceContext) findCategoriesByIDsQuery {
	return findCategoriesByIDsQuery{
		sctx: sctx,
	}
}

func (q *findCategoriesByIDsQuery) Execute(ctx context.Context, ids []common.UUID) ([]CategoryDTO, error) {
	db := q.sctx.MustGet(common.KeyGormComponent).(common.GormCompContext).GetDB()
	var categories []CategoryDTO

	if err := db.Table(CategoryDTO{}.TableName()).
		Where("id in (?)", ids).
		Find(&categories).Error; err != nil {
			return nil, core.ErrBadRequest.WithError("cannot list categories").WithDebug(err.Error())
	}
	return categories, nil
}
