package main

import (
	"context"
	cartHandler "github.com/aaanger/ecommerce/internal/cart/handler"
	orderHandler "github.com/aaanger/ecommerce/internal/order/handler"
	productHandler "github.com/aaanger/ecommerce/internal/product/handler"
	userHandler "github.com/aaanger/ecommerce/internal/user/handler"
	"github.com/aaanger/ecommerce/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	logrus.SetFormatter(new(logrus.JSONFormatter))

	err := initConfig()
	if err != nil {
		logrus.Fatalf("Error initializing config: %s", err)
	}

	err = godotenv.Load()
	if err != nil {
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
		logrus.Fatalf("Error loading database: %s", err)
	}

	router := gin.Default()

	userHandler.UserRoutes(router, db)
	productHandler.ProductRoutes(router, db)
	cartHandler.CartRoutes(router, db)
	orderHandler.OrderRoutes(router, db)

	srv := new(Server)

	go func() {
		err = srv.Run(viper.GetString("port"), router)
		if err != nil {
			logrus.Fatalf("Error running the server: %s", err)
		}
	}()

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
