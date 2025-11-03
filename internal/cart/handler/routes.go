package handler

import (
	"database/sql"
	"github.com/aaanger/ecommerce/internal/cart/repository"
	"github.com/aaanger/ecommerce/internal/cart/service"
	productRepository "github.com/aaanger/ecommerce/internal/product/repository"
	"github.com/aaanger/ecommerce/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func CartRoutes(r *gin.Engine, db *sql.DB, log *zap.Logger, redisClient *redis.Client) {
	repo := repository.NewCartRepository(db)
	redisRepo := repository.NewRedisCartRepository(redisClient, repository.TTL, log)
	productRepo := productRepository.NewProductRepository(db)
	svc := service.NewCartService(repo, redisRepo, productRepo, log)
	h := NewCartHandler(svc, log)

	cart := r.Group("/cart", middleware.SessionMiddleware)

	cart.GET("/", h.GetCart)
	cart.POST("/add", h.AddProduct)
	cart.DELETE("/", h.DeleteProduct)

}
