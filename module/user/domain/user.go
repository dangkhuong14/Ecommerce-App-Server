package domain

import (
	"ecommerce/common"
	"strings"
)

type User struct {
	id        common.UUID
	firstName string
	lastName  string
	email     string
	password  string
	salt      string
	role      Role
	status    string
	avatar    *Avatar
}

type Avatar struct {
	ImageID   common.UUID
	ImageName string
	ImageCDN  string
}

func NewAvatar(imageID common.UUID, imageName string, imageCDN string) (*Avatar, error) {
	return &Avatar{
		ImageID:   imageID,
		ImageName: imageName,
		ImageCDN:  imageCDN,
	}, nil
}

// Constructor function for User
func NewUser(
	id common.UUID, firstName string, lastName string,
	email string, password string, salt string, status string, role Role, avatar *Avatar) (*User, error) {
	return &User{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		password:  password,
		salt:      salt,
		role:      role,
		status:    status,
		avatar:    avatar,
	}, nil
}

func (u *User) GetID() common.UUID {
	return u.id
}

func (u *User) GetFirstName() string {
	return u.firstName
}

func (u *User) GetLastName() string {
	return u.lastName
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) GetSalt() string {
	return u.salt
}

func (u *User) GetRole() Role {
	return u.role
}

func (u *User) GetStatus() string {
	return u.status
}

type Role int

const (
	RoleUser  = iota
	RoleAdmin = 1
)

func (r Role) String() string {
	return [2]string{"user", "admin"}[r]
}

func GetRole(s string) Role {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "admin":
		return RoleAdmin
	default:
		return RoleUser
	}
}

func (u *User) GetAvatar() *Avatar {
	return u.avatar
}

func (u *User) SetAvatar(avatar *Avatar) {
	u.avatar = avatar
}
