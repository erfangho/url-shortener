package service

import (
	"errors"

	"github.com/erfangho/url-shortener/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepositoryInterface interface {
	Create(user *model.User) error
	FindByUserName(username string) (*model.User, error)
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

func (s *UserService) FindByUserName(username string) (*model.User, error) {
	user, err := s.repo.FindByUserName(username)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (s *UserService) Authenticate(user *model.User, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
}
