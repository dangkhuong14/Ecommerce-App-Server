package controller

import (
	"net/http"
	"ecommerce/module/product/domain"
	"ecommerce/module/product/domain/usecase"
	"ecommerce/module/product/repository/mysql"

	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

func CreateProductAPI(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse product data from body
		var productData domain.ProductCreationDTO

		if err := c.Bind(&productData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Create mysql repository
		mysqlRepo := mysql.NewMysqlRepository(db)
		// mysqlRepository.CreateProduct(c, productData)

		// Construct and call use case
		createProductUseCase := usecase.NewCreateProductUseCase(mysqlRepo)
		if err := createProductUseCase.CreateProduct(c.Request.Context(), &productData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Response to client
		c.JSON(http.StatusCreated, gin.H{
			"data": productData.ID.String(),
		})
	}
		
}