package app

import (
	"comparei-servico-proconf/internal/domain/users"
	users_interface "comparei-servico-proconf/internal/domain/users/interface"
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
func (s *UserService) UpdateLevelUser(u *users.User) error {
	log.Println("EXEC: service.UpdateLevelUser")
	err := s.userRepo.UpdateLevelUser(u)
	return err
}

