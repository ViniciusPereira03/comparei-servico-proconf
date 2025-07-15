package users_interface

import "comparei-servico-promer/internal/domain/users"

type UserRepository interface {
	CreateUser(user *users.User) (*users.User, error)
	GetUser(id string) (*users.User, error)
	UpdateUser(u *users.User) error
	DeleteUser(u *users.User) error
}
