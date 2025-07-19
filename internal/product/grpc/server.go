package productgrpc

import (
	"context"
	"database/sql"
	"github.com/aaanger/ecommerce/internal/product/repository"
	"github.com/aaanger/ecommerce/internal/product/service"
	pb "github.com/aaanger/ecommerce/proto/gen/product"
	"google.golang.org/grpc"
)

type ProductGRPCHandler struct {
	pb.UnimplementedProductServiceServer
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

	pb.RegisterProductServiceServer(srv, grpcHandler)
}

func (h *ProductGRPCHandler) ReserveProducts(ctx context.Context, req *pb.ReserveProductsReq) (*pb.ReserveProductsRes, error) {
	return h.service.ReserveProducts(ctx, req)
}
