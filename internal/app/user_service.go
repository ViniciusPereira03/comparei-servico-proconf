package app

import (
	"comparei-servico-promer/internal/domain/users"
	users_interface "comparei-servico-promer/internal/domain/users/interface"
	"log"
)

type UserService struct {
	userRepo users_interface.UserRepository
}

func NewUserService(userRepo users_interface.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(user *users.User) error {
	log.Println("EXEC: service.CreateUser")
	user, err := s.userRepo.CreateUser(user)
	return err
}
