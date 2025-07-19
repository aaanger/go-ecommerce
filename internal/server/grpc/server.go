package grpc

import (
	"context"
	"database/sql"
	"fmt"
	productgrpc "github.com/aaanger/ecommerce/internal/product/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

type Server struct {
	engine *grpc.Server
	log    *zap.Logger
	port   int
}

func InterceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		zapFields := make([]zapcore.Field, 0, len(fields))
		for _, f := range fields {
			if zf, ok := f.(zapcore.Field); ok {
				zapFields = append(zapFields, zf)
			}
		}
		switch level {
		case logging.LevelDebug:
			l.Debug(msg, zapFields...)
		case logging.LevelInfo:
			l.Info(msg, zapFields...)
		case logging.LevelWarn:
			l.Warn(msg, zapFields...)
		case logging.LevelError:
			l.Error(msg, zapFields...)
		default:
			l.Info(msg, zapFields...)
		}
	})
}

func NewServer(log *zap.Logger, db *sql.DB, port int) *Server {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent),
	}
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			log.Error("Recovered from panic", zap.Any("panic", p))
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...)))

	productgrpc.RegisterProductGRPCServer(grpcServer, db)

	return &Server{
		engine: grpcServer,
		log:    log,
		port:   port,
	}
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("grpc server run: %w", err)
	}

	err = s.engine.Serve(lis)
	if err != nil {
		return fmt.Errorf("grpc server run: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.engine.GracefulStop()
}
