package cmd

import (
	"fmt"
	"os"

	"log"

	"github.com/spf13/cobra"

	"net/http"

	"ecommerce/cmd/consumer"
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	productservice "ecommerce/module/product/infras"

	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/component/gormc"

	categorygrpcservice "ecommerce/module/category/infras/grpcservice"
	categoryservice "ecommerce/module/category/infras/httpservice"
)

func newService() sctx.ServiceContext {
	newSctx := sctx.NewServiceContext(sctx.WithName("G11"),
		sctx.WithComponent(gormc.NewGormDB(common.KeyGormComponent, "")),
		sctx.WithComponent(component.NewNATSComponent(common.KeyNatsComponent)),
		sctx.WithComponent(component.NewJWT(common.KeyJwtComponent)),
		sctx.WithComponent(component.NewAWSS3Provider(common.KeyAwsS3Component)),
		sctx.WithComponent(component.NewConfigComponent(common.KeyConfigComponent)),
	)
	return newSctx
}

var rootCmd = &cobra.Command{
	Use: "app",
	Short: "Start main service",
	Run: func(cmd *cobra.Command, args []string){
		
		// gin.SetMode(gin.ReleaseMode)

		// Create service context with component initialized
		sv := newService()

		// Print env variables
		// sv.OutEnv()

		if err := sv.Load(); err != nil {
			log.Fatalln("Error loading service: ", err)
		}

		// Get config component from service context
		configComp := (sv.MustGet(common.KeyConfigComponent)).(common.ConfigCompContext)

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

		// <------------------------Category service-------------------------->
		catService := categoryservice.NewHttpService(sv) 
		catService.Routes(v1)

		// category's gRPC server
		go func() {
			categoryGRPCPort := configComp.GetCategoryGRPCPort()
			gRPCCategoryServer := categorygrpcservice.NewGRPCCategoryServer(sv, categoryGRPCPort)
			gRPCCategoryServer.Start()
		}()

		// Create Category GRPC client connection
		opts := grpc.WithTransportCredentials(insecure.NewCredentials())

		cc, err := grpc.NewClient(configComp.GetCategoryGrpcUrl(), opts)

		if err != nil {
			log.Fatal(err)
		}

		// <------------------------Product service-------------------------->
		productService := productservice.NewHttpService(sv, cc)
		productService.Routes(v1)

		// HTTP server
		r.Run(":3000")
	},
}

func Execute() {
	// Add Out env command before calling Execute method
	rootCmd.AddCommand(outEnvCmd)

	consumerCmd := &cobra.Command{Use: "consumer", Short: "Start consumer"}
	consumerCmd.AddCommand(consumer.SetImgActiveAfterAvtChangeCmd)

	rootCmd.AddCommand(consumerCmd)
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} 