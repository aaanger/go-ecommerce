package handler

import (
	"database/sql"
	"github.com/aaanger/ecommerce/internal/cart/repository"
	"github.com/aaanger/ecommerce/internal/cart/service"
	productRepository "github.com/aaanger/ecommerce/internal/product/repository"
	"github.com/aaanger/ecommerce/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func CartRoutes(r *gin.Engine, db *sql.DB) {
	repo := repository.NewCartRepository(db)
	productRepo := productRepository.NewProductRepository(db)
	svc := service.NewCartService(repo, productRepo)
	h := NewCartHandler(svc)

	cart := r.Group("/cart", middleware.UserIdentity)

	cart.GET("/", h.GetCart)
	cart.POST("/add", h.AddProduct)
	cart.DELETE("/", h.DeleteProduct)

}
