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
}

// Constructor function for User
func NewUser(
	id common.UUID, firstName string, lastName string,
	email string, password string, salt string, status string, role Role,
) (*User, error) {
	return &User{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		password:  password,
		salt:      salt,
		role:      role,
		status:    status,
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
