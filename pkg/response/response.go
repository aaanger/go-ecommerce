package response

import (
	"github.com/gin-gonic/gin"
)

func JSON(c *gin.Context, code int, data any) {
	c.JSON(code, data)
}
