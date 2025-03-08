package service

import "github.com/2Cheetah/MedGuardianBot/internal/domain"

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) CreateUser(user *domain.User) error {
	return us.repo.CreateUser(user)
}
