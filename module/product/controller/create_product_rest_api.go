package controller

import (
	"net/http"
	"ecommerce/module/product/domain"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

func (api APIController) CreateProductAPI(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse product data from body
		var productData domain.ProductCreationDTO

		if err := c.Bind(&productData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := api.createUseCase.CreateProduct(c.Request.Context(), &productData); err != nil {
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