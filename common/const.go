package common

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const (
	KeyRequester     = "Requester"
	KeyGormComponent = "gorm"
	KeyJwtComponent  = "jwt"
)

type GormCompContext interface {
	GetDB() *gorm.DB
}


type TokenProvider interface {
	IssueToken(ctx context.Context, id, sub string) (token string, err error)
	TokenExpireInSeconds() int
	RefreshExpireInSeconds() int
	ParseToken(ctx context.Context, tokenString string) (claims *jwt.RegisteredClaims, err error)
}