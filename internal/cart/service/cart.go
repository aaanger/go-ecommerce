package service

import (
	"errors"
	"github.com/aaanger/ecommerce/internal/cart/model"
	"github.com/aaanger/ecommerce/internal/cart/repository"
	productRepository "github.com/aaanger/ecommerce/internal/product/repository"
	"go.uber.org/zap"
)

//go:generate mockery --name=ICartService

type ICartService interface {
	GetCartByUserID(userID int, sessionID string) (*model.Cart, error)
	AddProduct(userID, productID, quantity int, sessionID string) (*model.Cart, error)
	DeleteProduct(userID, productID int, sessionID string) (*model.Cart, error)
}

type CartService struct {
	repo        repository.ICartRepository
	redisRepo   repository.IRedisCartRepository
	productRepo productRepository.IProductRepository
	log         *zap.Logger
}

func NewCartService(repo repository.ICartRepository, redisRepo repository.IRedisCartRepository, productRepo productRepository.IProductRepository, log *zap.Logger) *CartService {
	return &CartService{
		repo:        repo,
		redisRepo:   redisRepo,
		productRepo: productRepo,
		log:         log,
	}
}

func (s *CartService) GetCartByUserID(userID int, sessionID string) (*model.Cart, error) {
	if userID == 0 {
		cart, err := s.redisRepo.GetCart(sessionID)
		if err != nil {
			return nil, err
		}
		return cart, nil
	}

	return s.repo.GetCartByUserID(userID)
}

func (s *CartService) AddProduct(userID, productID, quantity int, sessionID string) (*model.Cart, error) {
	log := s.log.With(
		zap.String("service", "cart"),
		zap.String("layer", "service"),
		zap.String("method", "AddProduct"),
		zap.Int("userID", userID))

	var totalPrice float64

	product, err := s.productRepo.GetProductByID(productID)
	if err != nil {
		log.Error("Get product error", zap.Error(err))
		return nil, err
	}
	if product.InStock == false {
		return nil, errors.New("product is not in stock")
	}

	if userID == 0 {
		cart, err := s.redisRepo.GetCart(sessionID)
		if err != nil {
			log.Error("Redis get cart error", zap.Error(err))
			return nil, err
		}
		err = s.redisRepo.AddProduct(sessionID, productID, quantity)
		if err != nil {
			log.Error("Redis add product error", zap.Error(err))
			return nil, err
		}

		cart.Lines = append(cart.Lines, model.CartLine{
			ProductID: productID,
			Product:   product,
			Quantity:  quantity,
		})

		for _, line := range cart.Lines {
			totalPrice += line.Product.Price
		}

		cart.TotalPrice = totalPrice
		return cart, nil
	}

	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		cartID, err := s.repo.CreateCart(userID)
		if err != nil {
			return nil, err
		}
		cart = &model.Cart{
			ID:     cartID,
			UserID: userID,
		}
	}

	err = s.repo.AddProduct(cart.ID, productID, quantity)
	if err != nil {
		return nil, err
	}

	cart.Lines = append(cart.Lines, model.CartLine{
		ProductID: productID,
		Product:   product,
		Quantity:  quantity,
	})

	for _, line := range cart.Lines {
		totalPrice += line.Product.Price
	}

	cart.TotalPrice = totalPrice

	return cart, nil
}

func (s *CartService) DeleteProduct(userID, productID int, sessionID string) (*model.Cart, error) {
	if userID == 0 {
		cart, err := s.redisRepo.GetCart(sessionID)
		if err != nil {

			return nil, err
		}
		err = s.redisRepo.DeleteProduct(sessionID, productID)
		if err != nil {
			return nil, err
		}

		for i, line := range cart.Lines {
			if line.ProductID == productID {
				cart.Lines = append(cart.Lines[:i], cart.Lines[i+1:]...)
				break
			}
		}

		return cart, nil
	}

	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	err = s.repo.DeleteProduct(cart.ID, productID)
	if err != nil {
		return nil, err
	}

	for i, line := range cart.Lines {
		if line.ProductID == productID {
			cart.Lines = append(cart.Lines[:i], cart.Lines[i+1:]...)
			break
		}
	}

	return cart, nil
}
