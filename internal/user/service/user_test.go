package service

import (
	"errors"
	"github.com/aaanger/ecommerce/internal/user/model"
	"github.com/aaanger/ecommerce/internal/user/repository/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserServiceSuite struct {
	suite.Suite
	repo    *mocks.IUserRepository
	service *UserService
}

func (suite *UserServiceSuite) SetupTest() {
	suite.repo = mocks.NewIUserRepository(suite.T())
	suite.service = NewUserService(suite.repo)
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}

// ====================================================================================================================

func (suite *UserServiceSuite) TestService_RegisterSuccess() {
	req := &model.UserReq{
		Email:    "test",
		Password: "test",
	}

	suite.repo.On("CreateUser", req.Email, req.Password, "user").Return(&model.User{
		ID:    1,
		Email: "test",
		Role:  "user",
	}, nil)

	user, err := suite.service.Register(req)
	suite.NotNil(user)
	suite.Nil(err)
}

func (suite *UserServiceSuite) TestService_RegisterFailure() {
	req := &model.UserReq{
		Email:    "test",
		Password: "test",
	}

	suite.repo.On("CreateUser", req.Email, req.Password, "user").Return(nil, errors.New("error"))

	user, err := suite.service.Register(req)
	suite.Nil(user)
	suite.NotNil(err)
}

// ====================================================================================================================

func (suite *UserServiceSuite) TestService_LoginSuccess() {
	req := &model.UserReq{
		Email:    "test",
		Password: "test",
	}

	suite.repo.On("AuthUser", req.Email, req.Password).Return(&model.User{
		ID:       1,
		Email:    "test",
		Password: "test",
		Role:     "user",
	}, nil)

	user, accessToken, refreshToken, err := suite.service.Login(req)

	suite.NotNil(user)
	suite.NotNil(accessToken)
	suite.NotNil(refreshToken)
	suite.Nil(err)
}

func (suite *UserServiceSuite) TestService_LoginFailure() {
	req := &model.UserReq{
		Email:    "test",
		Password: "test",
	}

	suite.repo.On("AuthUser", req.Email, req.Password).Return(nil, errors.New("error"))

	user, accessToken, refreshToken, err := suite.service.Login(req)

	suite.Nil(user)
	suite.Equal("", accessToken)
	suite.NotNil("", refreshToken)
	suite.NotNil(err)
}
