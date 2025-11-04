package main

import (
	"context"
	cartHandler "github.com/aaanger/ecommerce/internal/cart/handler"
	orderHandler "github.com/aaanger/ecommerce/internal/order/handler"
	grpcorder "github.com/aaanger/ecommerce/internal/order/handler/grpc/product"
	"github.com/aaanger/ecommerce/internal/order/service"
	payment "github.com/aaanger/ecommerce/internal/payment/client"
	productHandler "github.com/aaanger/ecommerce/internal/product/handler"
	"github.com/aaanger/ecommerce/internal/server/grpc"
	userHandler "github.com/aaanger/ecommerce/internal/user/handler"
	"github.com/aaanger/ecommerce/pkg/db"
	"github.com/aaanger/ecommerce/pkg/email"
	"github.com/aaanger/ecommerce/pkg/kafka"
	"github.com/aaanger/ecommerce/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (srv *Server) Run(port string, handler http.Handler) error {
	srv.httpServer = &http.Server{
		Addr:    port,
		Handler: handler,
	}

	return srv.httpServer.ListenAndServe()
}

func (srv *Server) Shutdown(ctx context.Context) error {
	return srv.httpServer.Shutdown(ctx)
}

func main() {
	logCfg := zap.NewProductionConfig()
	logCfg.OutputPaths = []string{"stdout", "var/log/ecom.log"}
	logCfg.ErrorOutputPaths = []string{"stderr"}

	logger, err := logCfg.Build()
	if err != nil {
		log.Fatalf("")
	}
	defer logger.Sync()

	if err := initConfig(); err != nil {
		logrus.Fatalf("Error initializing config: %s", err)
	}

	if err = godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env file: %s", err)
	}

	db, err := db.Open(db.PostgresConfig{
		os.Getenv("PSQL_HOST"),
		os.Getenv("PSQL_PORT"),
		os.Getenv("PSQL_USER"),
		os.Getenv("PSQL_PASSWORD"),
		os.Getenv("PSQL_DBNAME"),
		os.Getenv("PSQL_SSLMODE"),
	})
	if err != nil {
		logrus.Fatalf("Error loading PostgreSQL database: %s", err)
	}

	redisClient, err := redis.NewRedisClient(redis.RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if err != nil {
		logrus.Fatalf("Error loading Redis database: %s", err)
	}

	writer, reader := kafka.NewKafkaConnection(kafka.KafkaConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   service.CreateOrderTopic,
		GroupID: "1",
	})

	producer := kafka.NewProducer(writer, logger)
	consumer := kafka.NewConsumer(reader, logger)

	logger.Debug("Kafka producer and consumer loaded successfully")

	emailService, err := email.NewEmailService(
		os.Getenv("EMAIL_SENDER"),
		email.SMTPConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     os.Getenv("SMTP_PORT"),
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
		})
	if err != nil {
		logrus.Fatalf("Error initializing email service: %s", err)
	}
	orderConsumer := service.NewOrderConsumer(emailService, logger)

	go func() {
		productGrpcServer := grpc.NewServer(logger, db, 9090)
		productGrpcServer.MustRun()
	}()

	grpcClient, err := grpcorder.NewClient(context.Background(), logger, "localhost:9090", 3, 5*time.Second)
	if err != nil {
		logger.Error("error starting grpc client", zap.Error(err))
	}

	paymentClient := payment.NewClient(os.Getenv("SHOP_ID"), os.Getenv("SHOP_SECRET_KEY"))

	router := gin.Default()

	userHandler.UserRoutes(router, db)
	productHandler.ProductRoutes(router, db)
	cartHandler.CartRoutes(router, db, logger, redisClient)
	orderHandler.OrderRoutes(router, db, producer, grpcClient, paymentClient, orderConsumer, logger)

	srv := new(Server)

	go func() {
		err = srv.Run(viper.GetString("port"), router)
		if err != nil {
			logrus.Fatalf("Error running the server: %s", err)
		}
	}()

	go consumer.Consume(context.Background(), orderConsumer.HandleOrderCreated, 5)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	err = srv.Shutdown(context.Background())
	if err != nil {
		logrus.Errorf("Error shutting down the server: %s", err)
	}

	err = db.Close()
	if err != nil {
		logrus.Errorf("Error closing database: %s", err)
	}
}

func initConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("pkg/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
