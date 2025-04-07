package infras

import (
	"ecommerce/common"
	"ecommerce/gen/category"
	"ecommerce/module/product/domain/query"
	"ecommerce/module/product/repository/grpcclient"
	"net/http"

	"github.com/gin-gonic/gin"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
	"google.golang.org/grpc"
)

type httpService struct {
	sctx sctx.ServiceContext
	grpcClient *grpc.ClientConn
}

func NewHttpService(sctx sctx.ServiceContext, grpcClient *grpc.ClientConn) *httpService {
	return &httpService{
		sctx: sctx, 
		grpcClient: grpcClient,
	}
}

func (s *httpService) handleListProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var param query.ListProductQueryParam
		c.Bind(&param)

		// Create list product use case
		listProductQuery := query.NewListProductQuery(s.sctx)

		// Create Category repo (gRPC category)
		categoryServiceClient := category.NewCategoryServiceClient(s.grpcClient)
		categoryClientGRPC := grpcclient.NewCategoryGRPCClient(categoryServiceClient)

		// Create list product wrapper use case
		listProductWrapper := query.NewListProductWrapper(listProductQuery, categoryClientGRPC)

		// Call use case's method
		results, err := listProductWrapper.Execute(c.Request.Context(), &param)
		if err != nil {
			common.WriteErrorResponse(c, err)
		}

		c.JSON(http.StatusOK, core.SuccessResponse(results, param.Paging, param.ListProductFilterParam))

	}
}

func (s *httpService) Routes(g *gin.RouterGroup){
	products := g.Group("/products")
	{
		products.GET("/", s.handleListProduct())
	}
}
