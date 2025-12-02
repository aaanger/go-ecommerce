package productgrpc

import (
	"context"
	"database/sql"
	"github.com/aaanger/ecommerce/internal/product/repository"
	"github.com/aaanger/ecommerce/internal/product/service"
	"github.com/aaanger/ecommerce/pkg/proto/gen/product"
	"google.golang.org/grpc"
)

type ProductGRPCHandler struct {
	product.UnimplementedProductServiceServer
	service service.IProductService
}

func NewProductGRPCServer(service service.IProductService) *ProductGRPCHandler {
	return &ProductGRPCHandler{
		service: service,
	}
}

func RegisterProductGRPCServer(srv *grpc.Server, db *sql.DB) {
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)

	grpcHandler := NewProductGRPCServer(svc)

	product.RegisterProductServiceServer(srv, grpcHandler)
}

func (h *ProductGRPCHandler) ReserveProducts(ctx context.Context, req *product.ReserveProductsReq) (*product.ReserveProductsRes, error) {
	return h.service.ReserveProducts(ctx, req)
}
