package handler

import (
	"github.com/aaanger/ecommerce/internal/user/model"
	"github.com/aaanger/ecommerce/internal/user/service"
	"github.com/aaanger/ecommerce/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandler struct {
	service service.IUserService
}

func NewUserHandler(service service.IUserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var req model.UserReq

	err := c.BindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input parameters")
		return
	}

	user, err := h.service.Register(&req)
	if err != nil {
		logrus.Errorf("Service register error: %s", err)
		response.Error(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	res := model.RegisterRes{
		ID:    user.ID,
		Email: user.Email,
	}

	response.JSON(c, http.StatusOK, res)
}

func (h *UserHandler) SignIn(c *gin.Context) {
	var req model.UserReq

	err := c.BindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input parameters")
		return
	}

	user, accessToken, refreshToken, err := h.service.Login(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	res := model.LoginRes{
		ID:           user.ID,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.JSON(c, http.StatusOK, res)
}
