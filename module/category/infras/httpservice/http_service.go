package httpservice

import (
	"ecommerce/common"
	"ecommerce/module/category/query"
	"net/http"

	"github.com/gin-gonic/gin"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
)

type categoryHttpService struct {
	sctx sctx.ServiceContext
}

func NewHttpService(sctx sctx.ServiceContext) categoryHttpService {
	return categoryHttpService{sctx: sctx}
}

func (s categoryHttpService) handleFindCategoriesByIDs() gin.HandlerFunc {
	return func(c *gin.Context) {
		var param struct {
			Ids []common.UUID `json:"ids"`
		}

		if err := c.BindJSON(&param); err != nil {
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		// Create list product use case
		findCategoriesByIDsQuery := query.NewFindCategoriesByIDsQuery(s.sctx)

		// Call use case's method
		results, err := findCategoriesByIDsQuery.Execute(c.Request.Context(), param.Ids)
		if err != nil {
			common.WriteErrorResponse(c, err)
		}

		c.JSON(http.StatusOK, core.ResponseData(results))

	}
}

func (s categoryHttpService) Routes(g *gin.RouterGroup){
	category := g.Group("category")

	rpc := category.Group("rpc")
	{
		rpc.POST("/query-categories-ids", s.handleFindCategoriesByIDs())
	}
}