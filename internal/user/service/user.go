package service

import (
	"fmt"
	"github.com/aaanger/ecommerce/internal/user/model"
	"github.com/aaanger/ecommerce/internal/user/repository"
	"github.com/aaanger/ecommerce/pkg/jwt"
)

//go:generate mockgen -source=user.go -destination=mocks/mock.go

type IUserService interface {
	Register(req *model.UserReq) (*model.User, error)
	Login(req *model.UserReq) (*model.User, string, string, error)
	GetEmail(userID int) string
}

type UserService struct {
	repo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Register(req *model.UserReq) (*model.User, error) {
	user, err := s.repo.CreateUser(req.Email, req.Password, "user")
	if err != nil {
		return nil, fmt.Errorf("service user register: %w", err)
	}

	return user, nil
}

func (s *UserService) Login(req *model.UserReq) (*model.User, string, string, error) {
	user, err := s.repo.AuthUser(req.Email, req.Password)
	if err != nil {
		return nil, "", "", err
	}

	accessToken := jwt.GenerateAccessToken(user.ID, user.Email, user.Role)
	refreshToken := jwt.GenerateRefreshToken(user.ID, user.Email, user.Role)

	return user, accessToken, refreshToken, nil
}

func (s *UserService) GetEmail(userID int) string {
	return s.repo.GetEmail(userID)
}
