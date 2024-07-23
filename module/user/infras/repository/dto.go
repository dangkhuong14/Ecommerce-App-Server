package repository

import (
	"ecommerce/common"
	"ecommerce/module/user/domain"
)

type UserDTO struct {
	Id        common.UUID `gorm:"column:id"`
	FirstName string      `gorm:"column:first_name"`
	LastName  string      `gorm:"column:last_name"`
	Password  string      `gorm:"column:password"`
	Salt      string      `gorm:"column:salt"`
	Role      string      `gorm:"column:role"`
}

func (dto UserDTO) ToEntity() (*domain.User, error){
	return domain.NewUser(
		common.GenNewUUID(),
		dto.FirstName,
		dto.LastName,
		dto.Password,
		dto.Salt,
		domain.GetRole(dto.Role),
	)
}