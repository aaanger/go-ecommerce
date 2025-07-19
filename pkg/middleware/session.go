package middleware

import (
	"github.com/aaanger/ecommerce/pkg/cookie"
	"github.com/aaanger/ecommerce/pkg/lib"
	"github.com/aaanger/ecommerce/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	bytesPerToken = 32
)

func SessionMiddleware(c *gin.Context) {
	sessionToken, err := c.Cookie(cookie.CookieSession)
	if err != nil {
		newToken, err := lib.String(bytesPerToken)
		if err != nil {
			response.Error(c, http.StatusBadGateway, "502 Bad Gateway")
		}
		cookie.SetCookie(c.Writer, cookie.CookieSession, newToken)
	}

	c.Set(cookie.CookieSession, sessionToken)
}
