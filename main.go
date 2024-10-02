package main

import (
	"log"
	"os"

	"net/http"

	// "ecommerce/common"
	"ecommerce/component"
	"ecommerce/middleware"
	"ecommerce/module/product/controller"
	"ecommerce/module/product/domain/usecase"
	mysqlRepo "ecommerce/module/product/repository/mysql"
	"ecommerce/module/user/infras/httpservice"

	userbuilder "ecommerce/builder"
	"ecommerce/module/user/infras/repository"
	userusecase "ecommerce/module/user/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Checking that an environment variable is present or not.
	mysqlConnStr, ok := os.LookupEnv("MYSQL_CONNECTION")
	if !ok {
		log.Fatalln("Missing MySQL connection string")
	}

	dsn := mysqlConnStr
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db = db.Debug()

	if err != nil {
		log.Fatalln("Cannot connect to MySQL:", err)
	}
	log.Println("Connected to MySQL:", db)

	// Set up auth middleware dependencies
	// Set up User service dependencies
	jwt_secret := os.Getenv("JWT_SECRET")

	tokenProvider := component.NewJWTProvider(
		jwt_secret,
		component.DefaultExpireTokenInSeconds,
		component.DefaultExpireRefreshInSeconds,
	)

	userRepo := repository.NewMysqlUser(db)
	sessionRepo := repository.NewMysqlSession(db)

	introspectTokenUC := userusecase.NewIntrospectTokenUC(userRepo, sessionRepo, tokenProvider)

	// Gin API Ping
	r := gin.Default()

	// r.Use(middleware.RequireAuth(introspectTokenUC))

	r.GET("/ping", middleware.RequireAuth(introspectTokenUC), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
		// requester := c.MustGet(common.KeyRequester).(common.Requester)
		// log.Println(requester)
	})

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



	// use case with normal constructor

	// userRepo := repository.NewMysqlUser(db)
	// sessionRepo := repository.NewMysqlSession(db)
	// userUseCase := userusecase.NewUseCase(userRepo, sessionRepo, &common.Hasher{}, tokenProvider)

	// use case with simple builder
	// userUseCaseWithBuilder := userusecase.NewUseCaseWithBuilder(userbuilder.NewSimpleBuilder(db, tokenProvider))

	// use case with complex builder
	userUseCaseWithCmplxBuilder := userusecase.NewUseCaseWithBuilder(userbuilder.NewCmplxBuilder(userbuilder.NewSimpleBuilder(db, tokenProvider)))
	userService := httpservice.NewUserService(userUseCaseWithCmplxBuilder)
	userService.Routes(v1)

	// revoke token dependencies
	
	v1.DELETE("/revoke-token", middleware.RequireAuth(introspectTokenUC), userService.HandleRevokeToken())

	r.Run(":3000")
}
