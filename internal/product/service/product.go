package service

import (
	"context"
	"github.com/aaanger/ecommerce/internal/product/model"
	"github.com/aaanger/ecommerce/internal/product/repository"
	pb "github.com/aaanger/ecommerce/pkg/proto/gen/product"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockery --name=IProductService

type IProductService interface {
	CreateProduct(req *model.ProductReq) (*model.Product, error)
	GetAllProducts() ([]model.Product, error)
	GetProductByID(id int) (*model.Product, error)
	UpdateProduct(id int, input model.UpdateProduct) error
	DeleteProduct(id int) error
	ReserveProducts(ctx context.Context, req *pb.ReserveProductsReq) (*pb.ReserveProductsRes, error)
	UnreserveProducts(ctx context.Context, req *pb.ReserveProductsReq) (*pb.ReserveProductsRes, error)
}

type ProductService struct {
	repo repository.IProductRepository
}

func NewProductService(repo repository.IProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) CreateProduct(req *model.ProductReq) (*model.Product, error) {
	product, err := s.repo.CreateProduct(req)
	if err != nil {
		logrus.Errorf("Product create error: %s", err)
		return nil, err
	}
	return product, nil
}

func (s *ProductService) GetAllProducts() ([]model.Product, error) {
	return s.repo.GetAllProducts()
}

func (s *ProductService) GetProductByID(id int) (*model.Product, error) {
	return s.repo.GetProductByID(id)
}

func (s *ProductService) UpdateProduct(id int, input model.UpdateProduct) error {
	return s.repo.UpdateProduct(id, input)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.DeleteProduct(id)
}

func (s *ProductService) ReserveProducts(ctx context.Context, req *pb.ReserveProductsReq) (*pb.ReserveProductsRes, error) {
	for _, item := range req.Products {
		product, err := s.GetProductByID(int(item.ProductID))
		if err != nil {
			return nil, err
		}
		if product.Amount < int(item.Quantity) {
			return nil, status.Errorf(codes.FailedPrecondition, "not enough amount for product %s", product.Name)
		}
		if product.InStock == false {
			return nil, status.Errorf(codes.FailedPrecondition, "product %s is not in stock", product.Name)
		}

		updatedAmount := product.Amount - int(item.Quantity)
		inStock := updatedAmount > 0

		err = s.UpdateProduct(product.ID, model.UpdateProduct{
			Amount:  &updatedAmount,
			InStock: &inStock,
		})
		if err != nil {
			return nil, err
		}
	}
	return &pb.ReserveProductsRes{Success: true}, nil
}

func (s *ProductService) UnreserveProducts(ctx context.Context, req *pb.ReserveProductsReq) (*pb.ReserveProductsRes, error) {
	for _, item := range req.Products {
		product, err := s.GetProductByID(int(item.ProductID))
		if err != nil {
			return nil, err
		}

		updatedAmount := product.Amount + int(item.Quantity)
		inStock := updatedAmount > 0
		err = s.UpdateProduct(product.ID, model.UpdateProduct{
			Amount:  &updatedAmount,
			InStock: &inStock,
		})
		if err != nil {
			return nil, err
		}
	}
	return &pb.ReserveProductsRes{Success: true}, nil
}
