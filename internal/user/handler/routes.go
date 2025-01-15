package handler

import (
	"database/sql"
	"github.com/aaanger/ecommerce/internal/user/repository"
	"github.com/aaanger/ecommerce/internal/user/service"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, db *sql.DB) {
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	h := NewUserHandler(svc)

	r.POST("/signup", h.SignUp)
	r.POST("/signin", h.SignIn)
}
