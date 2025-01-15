package middleware

import (
	"errors"
	"github.com/aaanger/ecommerce/pkg/jwt"
	"github.com/aaanger/ecommerce/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func UserIdentity(c *gin.Context) {
	header := c.GetHeader("Authorization")

	if header == "" {
		response.Error(c, http.StatusUnauthorized, "empty auth header")
		c.Abort()
		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 {
		response.Error(c, http.StatusUnauthorized, "invalid auth header")
		c.Abort()
		return
	}

	userID, role, err := jwt.ParseToken(headerParts[1])
	if err != nil {
		logrus.WithError(err).Warn("invalid token")
		response.Error(c, http.StatusUnauthorized, "invalid token")
		c.Abort()
		return
	}

	logrus.Infof("Authenticated user id = %d, role = %s", userID, role)
	c.Set("userID", userID)
	c.Set("role", role)
}

func GetUserID(c *gin.Context) (int, error) {
	id, ok := c.Get("userID")
	if !ok {
		return 0, errors.New("GetUserID error: user id not found")
	}

	userID, ok := id.(int)
	if !ok {
		return 0, errors.New("GetUserID error: invalid type of user id")
	}

	return userID, nil
}

func ModeratorIdentity(c *gin.Context) {
	role, ok := c.Get("role")
	if !ok {
		response.Error(c, http.StatusUnauthorized, "role not found")
		c.Abort()
		return
	}

	if role != "moderator" {
		response.Error(c, http.StatusForbidden, "moderator role required")
		c.Abort()
		return
	}
}
