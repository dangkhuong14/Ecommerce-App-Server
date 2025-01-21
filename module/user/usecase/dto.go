package usecase

import "ecommerce/common"

type EmailPasswordRegistrationDTO struct {
	FirstName string `json:"first_name" binding:"required"`  // Không được rỗng
	LastName  string `json:"last_name" binding:"required"`   // Không được rỗng
	Email     string `json:"email" binding:"required,email"` // Không được rỗng, phải đúng định dạng email
	Password  string `json:"password" binding:"required"`    // Không được rỗng
}

type EmailPasswordLoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenResponseDTO struct {
	AccessToken       string `json:"access_token"`
	AccessTokenExpIn  int    `json:"token_exp_in"`
	RefreshToken      string `json:"refresh_token"`
	RefreshTokenExpIn int    `json:"refresh_token_exp"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type SingleAvatarChangeDTO struct {
	ImageID   string           `json:"image_id" binding:"required"`
	Requester common.Requester `json:"-" binding:"-"`
	CDNDomain  string           `json:"-" binding:"-"`
}
