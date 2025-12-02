package service

import (
	"context"
	"errors"
	"fmt"
	grpcorder "github.com/aaanger/ecommerce/internal/order/handler/grpc/product"
	"github.com/aaanger/ecommerce/internal/order/model"
	"github.com/aaanger/ecommerce/internal/order/repository"
	payment "github.com/aaanger/ecommerce/internal/payment/client"
	paymentModel "github.com/aaanger/ecommerce/internal/payment/model"
	productModel "github.com/aaanger/ecommerce/internal/product/model"
	productRepository "github.com/aaanger/ecommerce/internal/product/repository"
	"github.com/aaanger/ecommerce/pkg/kafka"
	pb "github.com/aaanger/ecommerce/pkg/proto/gen/product"
	_ "github.com/vektra/mockery/mockery"
	"go.uber.org/zap"
	"strconv"
)

const (
	CreateOrderTopic = "order_created"
)

//go:generate mockery --name=IOrderService

type IOrderService interface {
	CreateOrder(ctx context.Context, userID int, userEmail string, lines *model.CreateOrderReq) (*model.CreateOrderRes, error)
	ConfirmOrder(ctx context.Context, orderID int) error
	CancelOrder(ctx context.Context, orderID int) error
	GetOrderByID(orderID int) (*model.Order, error)
	GetAllOrders(userID int) ([]model.Order, error)
	UpdateOrderStatus(orderID int, status string) (*model.Order, error)
	ReserveProducts(ctx context.Context, lines []model.OrderLineReq) error
}

type OrderService struct {
	repo          repository.IOrderRepository
	productRepo   productRepository.IProductRepository
	grpcClient    *grpcorder.OrderGRPCClient
	paymentClient *payment.Client
	producer      *kafka.Producer
	log           *zap.Logger
}

func NewOrderService(repo repository.IOrderRepository, productRepo productRepository.IProductRepository, grpcClient *grpcorder.OrderGRPCClient, paymentClient *payment.Client, producer *kafka.Producer, log *zap.Logger) *OrderService {
	return &OrderService{
		repo:          repo,
		productRepo:   productRepo,
		grpcClient:    grpcClient,
		paymentClient: paymentClient,
		producer:      producer,
		log:           log,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID int, userEmail string, req *model.CreateOrderReq) (*model.CreateOrderRes, error) {
	log := s.log.With(
		zap.String("service", "order"),
		zap.String("layer", "service"),
		zap.String("method", "CreateOrder"),
		zap.Int("userID", userID))

	var lines []model.OrderLine

	for _, line := range req.Lines {
		lines = append(lines, model.OrderLine{
			ProductID: line.ProductID,
			Quantity:  line.Quantity,
		})
	}

	productMap := make(map[int]*productModel.Product)
	for i := range lines {
		log.Debug("Fetching product data", zap.Int("productID", lines[i].ProductID))
		product, err := s.productRepo.GetProductByID(lines[i].ProductID)
		if err != nil {
			log.Error("Error fetching product data", zap.Error(err), zap.Int("productID", lines[i].ProductID))
			return nil, err
		}
		productMap[product.ID] = product
		lines[i].Price = product.Price * float64(lines[i].Quantity)
	}

	if err := s.ReserveProducts(ctx, req.Lines); err != nil {
		return nil, err
	}

	log.Debug("Starting creating order")
	order, err := s.repo.CreateOrder(userID, userEmail, lines)
	if err != nil {
		log.Error("Error creating order", zap.Error(err))
		return nil, err
	}

	for i := range order.Lines {
		order.Lines[i].Product = productMap[lines[i].ProductID]
	}

	paymentReq := &paymentModel.CreatePaymentReq{
		Amount: paymentModel.Amount{
			Value:    fmt.Sprintf("%.2f", order.TotalPrice),
			Currency: "RUB",
		},
		Capture: true,
		Confirmation: paymentModel.ConfirmationReq{
			Type:      "redirect",
			ReturnURL: "http://localhost:3000/payment/success",
		},
		Metadata: map[string]string{
			"order_id": strconv.Itoa(order.ID),
		},
		Description: fmt.Sprintf("Заказ №%d", order.ID),
	}

	paymentRes, err := s.paymentClient.CreatePayment(ctx, paymentReq)
	if err != nil {
		log.Error("Failed to create payment", zap.Error(err))
		return nil, err
	}

	return &model.CreateOrderRes{
		Order:   order,
		Payment: paymentRes,
	}, nil
}

func (s *OrderService) ConfirmOrder(ctx context.Context, orderID int) error {
	log := s.log.With(
		zap.String("service", "order"),
		zap.String("layer", "service"),
		zap.String("method", "ConfirmOrder"),
		zap.Int("orderID", orderID))

	order, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		log.Error("failed to get order by id", zap.Error(err))
		return err
	}

	if order.Status != model.StatusPending {
		log.Error("order is already paid")
		return fmt.Errorf("order is already paid")
	}

	if err := s.repo.UpdateOrder(orderID, model.StatusCreated); err != nil {
		log.Error("failed to update order status", zap.Error(err))
		return err
	}

	log.Info("Producing order to Kafka",
		zap.Int("orderID", order.ID),
		zap.String("userEmail", order.UserEmail),
		zap.Any("order", order),
	)
	if err = s.producer.Produce(ctx, strconv.Itoa(order.ID), order, 3); err != nil {
		log.Error("Kafka produce error, topic `order_created`", zap.Int("orderID", order.ID))
		return err
	}
	log.Info("Kafka message produced in topic `order_created`", zap.Any("order", order))

	log.Info("Order successfully confirmed", zap.Int("orderID", order.ID), zap.Any("order", order))

	return nil
}

func (s *OrderService) GetOrderByID(orderID int) (*model.Order, error) {
	order, err := s.repo.GetOrderByID(orderID)
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

func (s *OrderService) UpdateOrderStatus(orderID int, status string) (*model.Order, error) {
	order, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}

	if order.Status == model.StatusDelivered || order.Status == model.StatusCanceled {
		return nil, errors.New("failed to update status: order delivered or canceled")
	}

	if status == model.StatusDelivering || status == model.StatusDelivered {
		err = s.repo.UpdateOrder(orderID, status)
		if err != nil {
			return nil, err
		}
		return order, nil
	}

	return nil, errors.New("invalid status")
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID int) error {
	order, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	if order.Status == model.StatusDelivered || order.Status == model.StatusCanceled {
		return errors.New("invalid order status")
	}

	err = s.UnreserveProducts(ctx, order.Lines)
	if err != nil {
		return err
	}
	err = s.repo.UpdateOrder(orderID, model.StatusCanceled)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) ReserveProducts(ctx context.Context, lines []model.OrderLineReq) error {
	var products []*pb.ReservedProduct

	for _, line := range lines {
		products = append(products, &pb.ReservedProduct{
			ProductID: int32(line.ProductID),
			Quantity:  int32(line.Quantity),
		})
	}

	res, err := s.grpcClient.Client.ReserveProducts(ctx, &pb.ReserveProductsReq{
		Products: products,
	})
	if err != nil {
		return err
	}
	if !res.Success {
		return fmt.Errorf("reservation failed")
	}

	return nil
}

func (s *OrderService) UnreserveProducts(ctx context.Context, lines []model.OrderLine) error {
	var products []*pb.ReservedProduct

	for _, line := range lines {
		products = append(products, &pb.ReservedProduct{
			ProductID: int32(line.ProductID),
			Quantity:  int32(line.Quantity),
		})
	}

	res, err := s.grpcClient.Client.UnreserveProducts(ctx, &pb.ReserveProductsReq{
		Products: products,
	})
	if err != nil {
		return err
	}
	if !res.Success {
		return fmt.Errorf("unreservation failed")
	}

	return nil
}
