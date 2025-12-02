package repository

import (
	"encoding/json"
	"errors"
	"github.com/aaanger/ecommerce/internal/cart/model"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"time"
)

const (
	TTL = 24 * time.Hour
)

type IRedisCartRepository interface {
	GetCart(sessionID string) (*model.Cart, error)
	AddProduct(sessionID string, productID, quantity int) error
	DeleteProduct(sessionID string, productID int) error
}

type RedisCartRepository struct {
	db  *redis.Client
	ttl time.Duration
	log *zap.Logger
}

func NewRedisCartRepository(client *redis.Client, ttl time.Duration, log *zap.Logger) *RedisCartRepository {
	return &RedisCartRepository{
		db:  client,
		ttl: ttl,
		log: log,
	}
}

func (r *RedisCartRepository) GetCart(sessionID string) (*model.Cart, error) {
	data, err := r.db.Get("cart:" + sessionID).Result()
	if errors.Is(err, redis.Nil) {
		return &model.Cart{}, nil
	} else if err != nil {
		return nil, err
	}

	var cart model.Cart

	err = json.Unmarshal([]byte(data), &cart)
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *RedisCartRepository) AddProduct(sessionID string, productID, quantity int) error {
	log := r.log.With(
		zap.String("storage", "redis"),
		zap.String("method", "AddProduct"))

	var cart model.Cart

	data, err := r.db.Get("cart:" + sessionID).Result()
	if errors.Is(err, redis.Nil) {
		cart = model.Cart{
			Lines: []model.CartLine{},
		}
	} else if err != nil {
		log.Error("failed to get cart", zap.Error(err))
		return err
	}

	if data == "" {
		cart = model.Cart{
			Lines: []model.CartLine{},
		}
	} else {
		err = json.Unmarshal([]byte(data), &cart)
		if err != nil {
			log.Error("json unmarshal error", zap.Error(err), zap.String("data", data))
			return err
		}
	}

	cart.Lines = append(cart.Lines, model.CartLine{
		ProductID: productID,
		Quantity:  quantity,
	})

	encodedCart, err := json.Marshal(cart)
	if err != nil {
		log.Error("json marshal error", zap.Error(err))
		return err
	}

	return r.db.Set("cart:"+sessionID, encodedCart, r.ttl).Err()
}

func (r *RedisCartRepository) DeleteProduct(sessionID string, productID int) error {
	var cart model.Cart

	data, err := r.db.Get("cart:" + sessionID).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(data), &cart)
	if err != nil {
		return err
	}

	updatedLines := make([]model.CartLine, 0, len(cart.Lines))
	for _, line := range cart.Lines {
		if line.ProductID != productID {
			updatedLines = append(updatedLines, line)
		}
		if line.ProductID == productID {
			cart.TotalPrice -= line.Product.Price
		}
	}

	cart.Lines = updatedLines

	encoded, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	return r.db.Set("cart:"+sessionID, encoded, r.ttl).Err()
}
