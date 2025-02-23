package service

import (
	"errors"
	"github.com/aaanger/ecommerce/internal/order/model"
	"github.com/aaanger/ecommerce/internal/order/repository"
	productModel "github.com/aaanger/ecommerce/internal/product/model"
	productRepository "github.com/aaanger/ecommerce/internal/product/repository"
	"github.com/sirupsen/logrus"
	_ "github.com/vektra/mockery/mockery"
)

//go:generate mockery --name=IOrderService

type IOrderService interface {
	CreateOrder(userID int, lines *model.CreateOrderReq) (*model.Order, error)
	GetOrderByID(userID, orderID int) (*model.Order, error)
	GetAllOrders(userID int) ([]model.Order, error)
	UpdateOrderStatus(userID, orderID int, status string) (*model.Order, error)
	CancelOrder(userID, orderID int) (*model.Order, error)
}

type OrderService struct {
	repo        repository.IOrderRepository
	productRepo productRepository.IProductRepository
}

func NewOrderService(repo repository.IOrderRepository, productRepo productRepository.IProductRepository) *OrderService {
	return &OrderService{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *OrderService) CreateOrder(userID int, req *model.CreateOrderReq) (*model.Order, error) {
	var lines []model.OrderLine

	for _, line := range req.Lines {
		lines = append(lines, model.OrderLine{
			ProductID: line.ProductID,
			Quantity:  line.Quantity,
		})
	}

	logrus.Info(lines)

	productMap := make(map[int]*productModel.Product)
	for i := range lines {
		product, err := s.productRepo.GetProductByID(lines[i].ProductID)
		if err != nil {
			return nil, err
		}
		productMap[product.ID] = product
		lines[i].Price = product.Price * float64(lines[i].Quantity)
	}

	order, err := s.repo.CreateOrder(userID, lines)
	if err != nil {
		return nil, err
	}

	for i := range order.Lines {
		lines[i].Product = productMap[lines[i].ProductID]
	}

	return order, nil
}

func (s *OrderService) GetOrderByID(userID, orderID int) (*model.Order, error) {
	order, err := s.repo.GetOrderByID(userID, orderID)
	if err != nil {
		return nil, err
	}

	for i := range order.Lines {
		product, err := s.productRepo.GetProductByID(order.Lines[i].ProductID)
		if err != nil {
			return nil, err
		}
		order.Lines[i].Product = product
	}

	return order, nil
}

func (s *OrderService) GetAllOrders(userID int) ([]model.Order, error) {
	return s.repo.GetAllOrders(userID)
}

func (s *OrderService) UpdateOrderStatus(userID, orderID int, status string) (*model.Order, error) {
	order, err := s.repo.GetOrderByID(userID, orderID)
	if err != nil {
		return nil, err
	}

	if order.Status == model.StatusOrderDelivered || order.Status == model.StatusOrderCanceled {
		return nil, errors.New("failed to update status: order delivered or canceled")
	}

	if status == model.StatusOrderDelivering || status == model.StatusOrderDelivered {
		err = s.repo.UpdateOrder(userID, orderID, status)
		if err != nil {
			return nil, err
		}
		return order, nil
	}

	return nil, errors.New("invalid status")
}

func (s *OrderService) CancelOrder(userID, orderID int) (*model.Order, error) {
	order, err := s.repo.GetOrderByID(userID, orderID)
	if err != nil {
		return nil, err
	}

	if order.Status == model.StatusOrderDelivered || order.Status == model.StatusOrderCanceled {
		return nil, errors.New("invalid order status")
	}

	err = s.repo.UpdateOrder(userID, orderID, model.StatusOrderCanceled)
	if err != nil {
		return nil, err
	}
	order.Status = model.StatusOrderCanceled

	return order, nil
}
