package handler

import (
	"database/sql"
	"github.com/aaanger/ecommerce/internal/product/repository"
	"github.com/aaanger/ecommerce/internal/product/service"
	"github.com/aaanger/ecommerce/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(r *gin.Engine, db *sql.DB) {
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	h := NewProductHandler(svc)

	p := r.Group("/products", middleware.UserIdentity)

	p.POST("/create", middleware.ModeratorIdentity, h.CreateProduct)
	p.GET("/", h.GetProducts)
	p.GET("/:id", h.GetProductByID)
	p.PUT("/:id", middleware.ModeratorIdentity, h.UpdateProduct)
	p.DELETE("/:id", middleware.ModeratorIdentity, h.DeleteProduct)

}
