package controller

import (
	"context"
	"ecommerce/common"
	"ecommerce/module/product/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (api APIController) CreateProductAPI(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy requester từ context
		requester := c.MustGet(common.KeyRequester).(common.Requester)

		// Thêm requester vào context
		ctx := context.WithValue(c.Request.Context(), common.KeyRequester, requester)

		// Parse product data from body
		var productData domain.ProductCreationDTO

		if err := c.Bind(&productData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := api.createUseCase.CreateProduct(ctx, &productData); err != nil {
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