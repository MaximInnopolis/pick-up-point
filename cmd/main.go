package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	service "route/internal/app/api"
	"route/internal/app/cache"
	"route/internal/app/cli"
	"route/internal/app/config"
	"route/internal/app/kafka"
	"route/internal/app/models"
	"route/internal/app/module"
	"route/internal/app/repository/database"
	"route/internal/app/repository/postgresql"
	order "route/pkg/api/proto/order/v1/order/v1"
)

// main is entry point of the program
func main() {

	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize KafkaProducer and KafkaConsumer
	producer, err := kafka.NewKafkaProducer(cfg.KafkaConfig)
	if err != nil {
		fmt.Println("Failed to create Kafka producer:", err)
		os.Exit(1)
	}

	consumer, err := kafka.NewKafkaConsumer(cfg.KafkaConfig)
	if err != nil {
		fmt.Println("Failed to create Kafka consumer:", err)
		os.Exit(1)
	}

	// Create a new connection pool to database
	pool, err := database.NewPool(cfg.DbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer pool.Close()

	// Create a new Database with connection pool
	db := database.NewDatabase(pool)

	// Create a new repo with Database
	repo := postgresql.New(*db)

	// Create in-memory cache
	imCache := cache.NewIMCache[int, models.Order](cfg.CacheTTL)

	// Create a new module with explicit type parameters
	mod := module.New[int, models.Order](repo, imCache)

	// Create a map of commands
	commands := cli.NewCommands(mod)

	// Create a new CLI
	cliCommands := cli.New(commands, cfg.OutputMode, consumer, producer)

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	orderService := service.New(*mod)

	// Register the service with the server
	order.RegisterOrderServiceServer(grpcServer, orderService)

	listener, err := net.Listen("tcp", cfg.ServerConfig.GrpcPort)
	log.Println("Start")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("gRPC server listening on", cfg.ServerConfig.GrpcPort)

	// Start cache invalidation routine
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			imCache.InvalidateExpired()
		}
	}()

	go func() {
		if err = grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":9090", nil))
	}()

	if err = cliCommands.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
