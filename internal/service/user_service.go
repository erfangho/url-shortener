package service

import (
	"github.com/erfangho/url-shortener/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryInterface interface {
	Create(user *model.User) error
}

type UserService struct {
	repo UserRepositoryInterface
}

func NewUserService(repo UserRepositoryInterface) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(name, username, password string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     name,
		Username: username,
		Password: string(hashedPassword),
	}

	err = s.repo.Create(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
