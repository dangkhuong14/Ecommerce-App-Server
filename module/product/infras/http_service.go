package infras

import (
	"ecommerce/common"
	"ecommerce/module/product/domain/query"
	"net/http"

	"github.com/gin-gonic/gin"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
)

type httpService struct {
	sctx sctx.ServiceContext
}

func NewHttpService(sctx sctx.ServiceContext) httpService {
	return httpService{sctx: sctx}
}

func (s httpService) handleListProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var param query.ListProductQueryParam
		c.Bind(&param)

		// Create new use case
		listProductQuery := query.NewListProductQuery(s.sctx)
		// Call use case's method
		results, err := listProductQuery.Execute(c.Request.Context(), &param)
		if err != nil {
			common.WriteErrorResponse(c, err)
		}

		c.JSON(http.StatusOK, core.SuccessResponse(results, param.Paging, param.ListProductFilterParam))

	}
}

func (s httpService) Routes(g *gin.RouterGroup){
	products := g.Group("/products")
	{
		products.GET("/", s.handleListProduct())
	}
}
