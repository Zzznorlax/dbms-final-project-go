package domain

import (
	"dbmsbackend/util"
)

type UserType string

const (
	Buyer  UserType = "買家"
	Seller UserType = "賣家"
)

type User struct {
	ID             int
	Name           string
	Kind           UserType
	Email          string
	Phone          string
	HashedPassword string
	Salt           string
}

type UserDTO struct {
	ID    int
	Name  string
	Kind  UserType `json:"type"`
	Email string
	Phone string
}

type UserRespDTO struct {
	Data UserDTO `json:"data"`
}

func NewUser(name string, email string, phone string, kind UserType, password string) *User {

	salt, err := util.GenerateRandomString(24)

	if err != nil {
		salt = ""
	}

	hash := util.HashPassword(email, password, salt)

	user := User{
		Name:           name,
		Email:          email,
		Kind:           kind,
		Phone:          phone,
		Salt:           salt,
		HashedPassword: hash,
	}

	return &user
}

func (user *User) VerifyPassword(password string) bool {
	hash := util.HashPassword(user.Email, password, user.Salt)

	return hash == user.HashedPassword
}

func (user *User) ToDTO() *UserRespDTO {
	return &UserRespDTO{
		Data: UserDTO{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Phone,
			Kind:  user.Kind,
		},
	}
}
