package repository

import (
	"database/sql/driver"
	"ecommerce/common"
	"ecommerce/module/user/domain"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
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
	Avatar    *AvatarDTO  `gorm:"column:avatar"`
}

func (dto UserDTO) ToEntity() (*domain.User, error) {
	entityAvatar, err := dto.Avatar.toEntityAvatar()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return domain.NewUser(
		dto.Id,
		dto.FirstName,
		dto.LastName,
		dto.Email,
		dto.Password,
		dto.Salt,
		dto.Status,
		domain.GetRole(dto.Role),
		entityAvatar,
	)
}

type AvatarDTO struct {
	ImageID   string `json:"image_id"`
	ImageName string `json:"image_name"`
	ImageCDN  string `json:"image_cdn"`
}

// Implement driver.Valuer to store Avatar in JSON form
func (a *AvatarDTO) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Implement sql.Scanner to scan data in mysql to Avatar struct
func (a *AvatarDTO) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, a)
}

func (avatarDTO *AvatarDTO) toEntityAvatar() (*domain.Avatar, error) {
	if avatarDTO == nil {
		return nil, nil
	}
	imageUUID, err := common.ParseUUID(avatarDTO.ImageID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &domain.Avatar{
		ImageID:   imageUUID,
		ImageName: avatarDTO.ImageName,
		ImageCDN:  avatarDTO.ImageCDN,
	}, nil
}

// Convert domain Avatar to AvatarDTO
func toAvatarDTO(avatar *domain.Avatar) *AvatarDTO {
	if avatar == nil {
		return nil
	}
	return &AvatarDTO{
		ImageID:   avatar.ImageID.String(),
		ImageName: avatar.ImageName,
		ImageCDN:  avatar.ImageCDN,
	}
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
