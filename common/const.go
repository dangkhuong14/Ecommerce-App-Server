package common

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const (
	KeyRequester       = "Requester"
	KeyGormComponent   = "gorm"
	KeyJwtComponent    = "jwt"
	KeyAwsS3Component  = "aws_s3"
	KeyConfigComponent = "config"
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

type ImageSaver interface {
	SaveFileUploaded(ctx context.Context, data []byte, dst string) error
	GetName() string
	GetDomain() string
}

type ConfigCompContext interface {
	GetURLRPCCategory() string
	GetCategoryGRPCPort() int
}

type Paging struct {
	Page  int `json:"page" form:"page"`
	Total int `json:"total"`
	Limit int `json:"limit" form:"limit"`
}

// Use pointer receiver because process will change instance value
func (p *Paging) Process() {
	if p.Limit < 1 {
		p.Limit = 10
	}

	if p.Limit > 200 {
		p.Limit = 200
	}

	if p.Page < 1 {
		p.Page = 1
	}
}
