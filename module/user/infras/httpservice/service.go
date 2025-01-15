package httpservice

import (
	"ecommerce/common"
	"ecommerce/middleware"
	"ecommerce/module/user/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
)

type Service struct {
	uc usecase.UseCase
	sctx sctx.ServiceContext
}

func NewUserService(uc usecase.UseCase, sctx sctx.ServiceContext) Service {
	return Service{
		uc:uc,
		sctx: sctx,
	}
}

func (s Service) handleRegistration() gin.HandlerFunc{
	return func(c *gin.Context) {
		var dto = usecase.EmailPasswordRegistrationDTO{}
		// Parse body to user dto
		if err := c.Bind(&dto); err != nil {
			common.WriteErrorResponse(c, core.ErrBadRequest.WithWrap(err).WithDebug(err.Error()))
			return
		}
		// Call use case registration
		if err := s.uc.Register(c, dto); err != nil {
			common.WriteErrorResponse(c, err)
			return
		}
		// Response to client
		c.JSON(http.StatusOK, core.ResponseData(true))
	}
}

func (s Service) handleEmailPasswordLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse body data
		var dto  usecase.EmailPasswordLoginDTO

		if err := c.Bind(&dto); err != nil {
			common.WriteErrorResponse(c, core.ErrBadRequest.WithWrap(err).WithDebug(err.Error()))
			return
		}

		// Get token
		token, err := s.uc.LoginEmailPassword(c, dto)
		if err != nil {
			common.WriteErrorResponse(c, core.ErrInternalServerError.WithWrap(err).WithDebug(err.Error()))
			return
		}
		// Response token
		c.JSON(http.StatusOK, core.ResponseData(token))
	}
}

func (s Service) HandleRevokeToken() gin.HandlerFunc {
	return func(c *gin.Context){
		// Get requester from Gin context
		requester, ok := c.Get(common.KeyRequester)
		if !ok {
			common.WriteErrorResponse(c, core.ErrForbidden)
			return
		}
		// type assertion
		r, ok := requester.(common.Requester)
		if !ok {
			common.WriteErrorResponse(c, core.ErrInternalServerError.WithDebug("Can not assert type of Requester from parameter that is gotten from gin context"))
			return
		}

		// Revoke token
		err := s.uc.RevokeToken(c, r.TokenId())
		if err != nil {
			common.WriteErrorResponse(c, core.ErrInternalServerError.WithDebug(err.Error()))
			return
		}
		c.JSON(http.StatusOK, core.ResponseData(true))
	}
}

func (s Service) handleRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var refreshTokenDTO usecase.RefreshTokenDTO

		if err := c.BindJSON(&refreshTokenDTO); err != nil {
			common.WriteErrorResponse(c, core.ErrBadRequest.WithWrap(err).WithDebug(err.Error()))
			return
		}

		tokenResponse, err := s.uc.RefreshToken(c, refreshTokenDTO.RefreshToken)
		if err != nil {
			common.WriteErrorResponse(c, core.ErrInternalServerError.WithWrap(err).WithDebug(err.Error()))
			return
		}

		c.JSON(http.StatusOK, core.ResponseData(tokenResponse))
		
	}
}

func (s Service) Routes(g *gin.RouterGroup) {
	g.POST("/register",middleware.Recovery(), s.handleRegistration())
	g.POST("/login", s.handleEmailPasswordLogin())
	g.POST("/refresh-token", s.handleRefreshToken())
}

