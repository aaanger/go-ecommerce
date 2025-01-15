package service

import (
	"github.com/aaanger/ecommerce/internal/product/model"
	"github.com/aaanger/ecommerce/internal/product/repository"
	"github.com/sirupsen/logrus"
)

//go:generate mockery --name=IProductService

type IProductService interface {
	CreateProduct(req *model.ProductReq) (*model.Product, error)
	GetAllProducts() ([]model.Product, error)
	GetProductByID(id int) (*model.Product, error)
	UpdateProduct(id int, input model.UpdateProduct) error
	DeleteProduct(id int) error
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
