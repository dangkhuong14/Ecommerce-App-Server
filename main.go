package main

import (
	"log"

	"net/http"

	"ecommerce/common"
	"ecommerce/component"
	"ecommerce/middleware"
	"ecommerce/module/image"
	"ecommerce/module/product/controller"
	"ecommerce/module/product/domain/usecase"
	mysqlRepo "ecommerce/module/product/repository/mysql"
	"ecommerce/module/user/infras/httpservice"

	userbuilder "ecommerce/builder"
	"ecommerce/module/user/infras/repository"
	userusecase "ecommerce/module/user/usecase"

	"github.com/gin-gonic/gin"

	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/component/gormc"
)

func newService() sctx.ServiceContext {
	newSctx := sctx.NewServiceContext(sctx.WithName("G11"),
		sctx.WithComponent(gormc.NewGormDB(common.KeyGormComponent, "")),
		sctx.WithComponent(component.NewJWT(common.KeyJwtComponent)),
		sctx.WithComponent(component.NewAWSS3Provider(common.KeyAwsS3Component)),
	)
	return newSctx
}

func main() {

	// gin.SetMode(gin.ReleaseMode)

	// Create service context with component initialized
	sv := newService()

	sv.OutEnv()

	if err := sv.Load(); err != nil {
		log.Fatalln("Error loading service: ", err)
	}

	db := sv.MustGet(common.KeyGormComponent).(common.GormCompContext).GetDB()
	tokenProvider := sv.MustGet(common.KeyJwtComponent).(common.TokenProvider)

	userRepo := repository.NewMysqlUser(db)
	sessionRepo := repository.NewMysqlSession(db)

	introspectTokenUC := userusecase.NewIntrospectTokenUC(userRepo, sessionRepo, tokenProvider)

	r := gin.Default()

	// r.Use(middleware.RequireAuth(introspectTokenUC))
	r.Use(middleware.Recovery())

	// Gin API Ping
	r.GET("/ping", middleware.RequireAuth(introspectTokenUC), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
		// requester := c.MustGet(common.KeyRequester).(common.Requester)
		// log.Println(requester)
	})

	// <------------------------Product service-------------------------->
	// Set up dependencies
	mysqlRepo := mysqlRepo.NewMysqlRepository(db)
	createProductUseCase := usecase.NewCreateProductUseCase(mysqlRepo)
	api := controller.NewAPIController(createProductUseCase)

	// Create Product
	v1 := r.Group("/v1")
	{
		products := v1.Group("products")
		{
			products.POST("", middleware.RequireAuth(introspectTokenUC), api.CreateProductAPI(db))
		}
	}

	// <------------------------User service-------------------------->

	// use case with normal constructor
	// userRepo := repository.NewMysqlUser(db)
	// sessionRepo := repository.NewMysqlSession(db)
	// userUseCase := userusecase.NewUseCase(userRepo, sessionRepo, &common.Hasher{}, tokenProvider)

	// use case with simple builder
	// userUseCaseWithBuilder := userusecase.NewUseCaseWithBuilder(userbuilder.NewSimpleBuilder(db, tokenProvider))


	// use case with complex builder
	userUseCaseWithCmplxBuilder := userusecase.NewUseCaseWithBuilder(userbuilder.NewCmplxBuilder(userbuilder.NewSimpleBuilder(db, tokenProvider)))
	userService := httpservice.NewUserService(userUseCaseWithCmplxBuilder, sv)
	userService.Routes(v1)


	v1.DELETE("/revoke-token", middleware.RequireAuth(introspectTokenUC), userService.HandleRevokeToken())
	v1.POST("/change-avatar", middleware.RequireAuth(introspectTokenUC), userService.HandleChangeAvatar())
	
	// <------------------------Image service-------------------------->
	imageService := image.NewImageService(sv)
	imageService.Routes(v1)

	r.Run(":3000")
}
