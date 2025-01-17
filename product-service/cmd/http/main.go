package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	_ "product-service/docs"
	"product-service/internal/adapter/auth/jwt"
	"product-service/internal/adapter/broker/rabbitmq"
	"product-service/internal/adapter/config"
	httpLib "product-service/internal/adapter/handler/http"
	"product-service/internal/adapter/logger"
	"product-service/internal/adapter/storage/mongodb"
	"product-service/internal/adapter/storage/mongodb/repository"
	"product-service/internal/adapter/storage/redis"
	"product-service/internal/core/service"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// @title						product
// @version					1.0
// @description				Product Service
//
// @contact.name				Emmanuel Jonathan
// @contact.url				https://github.com/emmrys-jay
// @contact.email				jonathanemma121@gmail.com
//
// @host						localhost:8080
// @BasePath					/api/v1
// @schemes					http https
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and the access token.
func main() {
	// Load environment variables
	config := config.Setup()

	// Set logger
	l := logger.Get()
	zap.ReplaceGlobals(l)

	l.Info("Starting the application",
		zap.String("app", config.App.Name),
		zap.String("env", config.App.Env))

	// Init database
	ctx := context.Background()
	db, err := mongodb.New(ctx, &config.Database)
	if err != nil {
		l.Error("Error initializing database connection", zap.Error(err))
		os.Exit(1)
	}
	defer db.Close()

	l.Info("Successfully connected to the database")

	// Init cache service
	cache, err := redis.New(ctx, &config.Redis)
	if err != nil {
		l.Error("Error initializing cache connection", zap.Error(err))
		// os.Exit(1) // Cache is not being used at the moment
	}
	defer cache.Close()

	l.Info("Successfully connected to the cache server")

	// Init token service
	tokenService := jwt.New(&config.Token)

	// Message Queue Producer
	producer, err := rabbitmq.New(ctx, &config.Rabbitmq)
	if err != nil {
		l.Error("Error initializing Message Queue producer", zap.Error(err))
		os.Exit(1)
	}

	l.Info("Successfully connected to the message queue and created producer")

	// Dependency injection
	// Ping
	pingRepo := repository.NewPingRepository(db)
	pingService := service.NewPingService(pingRepo, cache)
	pingHandler := httpLib.NewPingHandler(pingService, validator.New())

	// Product
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo, cache, producer)
	productHandler := httpLib.NewProductHandler(productService, validator.New())

	// Init router
	router, err := httpLib.NewRouter(
		&config.Server,
		tokenService,
		l,
		*pingHandler,
		*productHandler,
	)
	if err != nil {
		l.Error("Error initializing router ", zap.Error(err))
		os.Exit(1)
	}

	// Message Queue Consumer 1
	consumer1, err := rabbitmq.New(ctx, &config.Rabbitmq)
	if err != nil {
		l.Error("Error initializing Message Queue consumer", zap.Error(err))
		os.Exit(1)
	}

	l.Info("Successfully connected to the message queue and created consumer")

	// Start consumer
	queue := "user-updates"
	l.Info("Starting consumer on", zap.String("queue", queue))
	go consumer1.Consume(ctx, queue, productService.UpdateProductsFromQueue)

	// Init GRPC server
	grpcListAddr := fmt.Sprintf("%s:%s", config.Server.GrpcUrl, config.Server.GrpcPort)
	list, err := net.Listen("tcp", grpcListAddr)
	l.Info("Starting the GRPC server", zap.String("listen_address", list.Addr().String()))

	server, err := service.NewGRPCServer(&service.Config{}, productRepo, grpc.EmptyServerOption{})
	go func() {
		l.Error("Error starting grpc server", zap.Error(server.Serve(list)))
	}()

	// Start server
	listenAddr := fmt.Sprintf("%s:%s", config.Server.HttpUrl, config.Server.HttpPort)
	l.Info("Starting the HTTP server", zap.String("listen_address", listenAddr))

	err = http.ListenAndServe(listenAddr, router)
	if err != nil {
		l.Error("Error starting the HTTP server", zap.Error(err))
		os.Exit(1)
	}
}
