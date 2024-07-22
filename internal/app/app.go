package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kuromii5/messagio/internal/app/server"
	"github.com/kuromii5/messagio/internal/config"
	"github.com/kuromii5/messagio/internal/db"
	"github.com/kuromii5/messagio/internal/kafka"
	"github.com/kuromii5/messagio/internal/service"
	"github.com/kuromii5/messagio/pkg/logger"
)

type App struct {
	server  *server.Server
	gateway *server.Gateway
}

func NewApp() *App {
	// Init config
	config := config.Load()

	// Init logger
	logger := logger.New(config.Env, config.LogLevel)

	// Init database
	db := db.NewDB(config.PGConfig)

	// Init Kafka producer and consumer
	producer := kafka.NewProducer(config.KafkaBrokers)
	consumer := kafka.NewConsumer(config.KafkaBrokers)

	// Init service
	msgService := service.NewService(logger, db, db, producer, consumer, config.KafkaTopic)

	// Init gateway
	gateway := server.NewGateway(
		config.HttpPort,
		config.GrpcPort,
		logger,
	)

	// Init server
	server := server.NewServer(
		logger,
		config.GrpcPort,
		msgService,
	)

	logger.Debug("",
		slog.Group("Settings",
			slog.Any("Kafka brokers", config.KafkaBrokers),
			slog.Any("Postgres", config.PGConfig),
			slog.String("Environment", config.Env),
			slog.Int("GRPC port", config.GrpcPort),
			slog.Int("HTTP port", config.HttpPort),
		),
	)

	return &App{server: server, gateway: gateway}
}

func (a *App) Run() {
	ctx := context.Background()

	// run grpc server
	go func() {
		a.server.Run()
	}()

	// run grpc gateway to receive http requests
	go func() {
		a.gateway.Run(ctx)
	}()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	a.Shutdown()
}

func (a *App) Shutdown() {
	a.server.Shutdown()
}
