package grpcorder

import (
	"context"
	grpc2 "github.com/aaanger/ecommerce/internal/server/grpc"
	pb "github.com/aaanger/ecommerce/pkg/proto/gen/product"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type OrderGRPCClient struct {
	Client pb.ProductServiceClient
}

func NewClient(ctx context.Context, log *zap.Logger, addr string, retriesCount int, timeout time.Duration) (*OrderGRPCClient, error) {
	retryOpts := []retry.CallOption{
		retry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		retry.WithMax(uint(retriesCount)),
		retry.WithPerRetryTimeout(timeout),
	}

	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.PayloadReceived, logging.PayloadSent),
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			logging.UnaryClientInterceptor(grpc2.InterceptorLogger(log), logOpts...),
			retry.UnaryClientInterceptor(retryOpts...)))
	if err != nil {
		return nil, err
	}

	client := pb.NewProductServiceClient(conn)

	return &OrderGRPCClient{
		Client: client,
	}, nil
}
