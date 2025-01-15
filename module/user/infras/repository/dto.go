package repository

import (
	"ecommerce/common"
	"ecommerce/module/user/domain"
	"time"
)

type UserDTO struct {
	Id        common.UUID `gorm:"column:id"`
	FirstName string      `gorm:"column:first_name;not null"`
	LastName  string      `gorm:"column:last_name;not null"` 
	Email     string      `gorm:"column:email;not null"`     
	Password  string      `gorm:"column:password;not null"`  
	Salt      string      `gorm:"column:salt"`
	Role      string      `gorm:"column:role;not null"`
	Status    string      `gorm:"column:status"`
}

func (dto UserDTO) ToEntity() (*domain.User, error) {
	return domain.NewUser(
		dto.Id,
		dto.FirstName,
		dto.LastName,
		dto.Email,
		dto.Password,
		dto.Salt,
		dto.Status,
		domain.GetRole(dto.Role),
	)
}

type SessionDTO struct {
	Id           common.UUID `gorm:"column:id"`
	UserId       common.UUID `gorm:"column:user_id"`
	RefreshToken string      `gorm:"column:refresh_token"`
	AccessExpAt  time.Time   `gorm:"column:access_exp_at"`
	RefreshExpAt time.Time   `gorm:"column:refresh_exp_at"`
}

type SessionUpdateDTO struct {
	Id           common.UUID `gorm:"column:id"`
	UserId       common.UUID `gorm:"column:user_id"`
	AcessToken   string      `gorm:"column: access_token"`
	RefreshToken string      `gorm:"column:refresh_token"`
	AccessExpAt  time.Time   `gorm:"column:access_exp_at"`
	RefreshExpAt time.Time   `gorm:"column:refresh_exp_at"`
}

func (sdto SessionUpdateDTO) ToEntity() (*domain.Session, error) {
	return domain.NewSession(
		sdto.Id,
		sdto.UserId,
		sdto.RefreshToken,
		sdto.AccessExpAt,
		sdto.RefreshExpAt,
	), nil
}
