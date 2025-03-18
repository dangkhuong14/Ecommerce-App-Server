package infras

import (
	"ecommerce/common"
	"ecommerce/module/product/domain/query"
	"ecommerce/module/product/repository/rpchttp"
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

		// Get category rpc url in config component
		urlCategoryRPC := s.sctx.MustGet(common.KeyConfigComponent).(common.ConfigCompContext).GetURLRPCCategory()

		// Create list product use case
		listProductQuery := query.NewListProductQuery(s.sctx)

		// Create Category repo (rpc category)
		rpcCategoryRepo := rpchttp.NewRPCFindCategoriesByIDs(urlCategoryRPC)

		// Create list product wrapper use case
		listProductWrapper := query.NewListProductWrapper(listProductQuery, rpcCategoryRepo)

		// Call use case's method
		results, err := listProductWrapper.Execute(c.Request.Context(), &param)
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
