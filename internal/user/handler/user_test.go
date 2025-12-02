package handler

import (
	"bytes"
	"github.com/aaanger/ecommerce/internal/user/model"

	"errors"
	mock_service "github.com/aaanger/ecommerce/internal/user/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestHandler_SignUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockIUserService, user *model.UserReq)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            *model.UserReq
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"test@test.com","password":"123"}`,
			inputUser: &model.UserReq{Email: "test@test.com", Password: "123"},
			mockBehavior: func(s *mock_service.MockIUserService, user *model.UserReq) {
				s.EXPECT().Register(user).Return(&model.User{ID: 0, Email: user.Email, Password: user.Password}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"Successfully registered":{"id":0,"email":"test@test.com"}}`,
		},
		{
			name:      "Empty fields",
			inputBody: `{"password":"123"}`,
			mockBehavior: func(s *mock_service.MockIUserService, user *model.UserReq) {

			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"Error":"Invalid input parameters"}`,
		},
		{
			name:      "Service failure",
			inputBody: `{"email":"test@test.com","password":"123"}`,
			inputUser: &model.UserReq{Email: "test@test.com", Password: "123"},
			mockBehavior: func(s *mock_service.MockIUserService, user *model.UserReq) {
				s.EXPECT().Register(user).Return(nil, errors.New("Something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"Error":"Something went wrong"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockIUserService(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			handler := NewUserHandler(auth)

			r := gin.New()
			r.POST("/signup", handler.SignUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/signup", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_SignIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockIUserService, user *model.UserReq)

	testCases := []struct {
		name                 string
		inputBody            string
		inputUser            *model.UserReq
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"test@test.com","password":"123"}`,
			inputUser: &model.UserReq{Email: "test@test.com", Password: "123"},
			mockBehavior: func(s *mock_service.MockIUserService, user *model.UserReq) {
				s.EXPECT().Login(user).Return(&model.User{ID: 0, Email: user.Email, Password: user.Password}, "access_token", "refresh_token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"Logged in":{"id":0,"email":"test@test.com","access_token":"access_token","refresh_token":"refresh_token"}}`,
		},
		{
			name:      "Empty fields",
			inputBody: `{"password":"123"}`,
			mockBehavior: func(s *mock_service.MockIUserService, user *model.UserReq) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"Error":"Invalid input parameters"}`,
		},
		{
			name:      "Service failure",
			inputBody: `{"email":"test@test.com","password":"123"}`,
			inputUser: &model.UserReq{Email: "test@test.com", Password: "123"},
			mockBehavior: func(s *mock_service.MockIUserService, user *model.UserReq) {
				s.EXPECT().Login(user).Return(nil, "access_token", "refresh_token", errors.New("Something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"Error":"Something went wrong"}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockIUserService(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			handler := NewUserHandler(auth)

			r := gin.New()
			r.POST("/signin", handler.SignIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/signin", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
