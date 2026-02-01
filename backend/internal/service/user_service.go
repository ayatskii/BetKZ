package service

import "BetKZ/internal/domain"

type UserService struct {
	repo interface {
		Create(domain.User) domain.User
	}
}

func NewUserService(r interface {
	Create(domain.User) domain.User
}) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) CreateUser(email string, balance float64) domain.User {
	return s.repo.Create(domain.User{
		Email:   email,
		Balance: balance,
	})
}
