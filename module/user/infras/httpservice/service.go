package httpservice

import (
	"ecommerce/module/user/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service struct {
	uc usecase.UseCase
}

func NewUserService(uc usecase.UseCase) Service {
	return Service{
		uc:uc,
	}
}

func (s Service) handleRegistration() gin.HandlerFunc{
	return func(c *gin.Context) {
		var dto = usecase.EmailPasswordRegistrationDTO{}
		// Parse body to user dto
		if err := c.Bind(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Call use case registration
		if err := s.uc.Register(c, dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Response to client
		c.JSON(http.StatusOK, gin.H{"data": true})
	}
}

func (s Service) Routes(g *gin.RouterGroup) {
	g.POST("/register", s.handleRegistration())
}

