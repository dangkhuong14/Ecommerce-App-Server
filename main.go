package main

import (
	"log"
	"os"

	"net/http"

	"ecommerce/common"
	"ecommerce/component"
	"ecommerce/module/product/controller"
	"ecommerce/module/product/domain/usecase"
	mysqlRepo "ecommerce/module/product/repository/mysql"
	"ecommerce/module/user/infras/httpservice"
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

	if err != nil {
		log.Fatalln("Cannot connect to MySQL:", err)
	}
	log.Println("Connected to MySQL:", db)

	// Gin API Ping
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
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
			products.POST("", api.CreateProductAPI(db))
		}
	}

	// Set up User service dependencies
	jwt_secret := os.Getenv("JWT_SECRET")

	tokenProvider := component.NewJWTProvider(
		jwt_secret,
		component.DefaultExpireTokenInSeconds,
		component.DefaultExpireRefreshInSeconds,
	)

	userRepo := repository.NewMysqlUser(db)
	sessionRepo := repository.NewMysqlSession(db)
	userUseCase := userusecase.NewUseCase(userRepo, sessionRepo, &common.Hasher{}, tokenProvider)
	httpservice.NewUserService(userUseCase).Routes(v1)


	r.Run(":3000")
}
