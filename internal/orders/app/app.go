package app

import (
	"context"
	"fmt"
	"log"
	"microservices_kafka_project/configs"
	"microservices_kafka_project/internal/common"
	orderApp "microservices_kafka_project/internal/orders/application"
	orderService "microservices_kafka_project/internal/orders/domain/service/orders"
	"microservices_kafka_project/internal/orders/infrastructure/kafka"
	"microservices_kafka_project/internal/orders/infrastructure/orders"
	internalhttp "microservices_kafka_project/internal/orders/ports/internal_http"
	publichttp "microservices_kafka_project/internal/orders/ports/public_http"
	database "microservices_kafka_project/pkg/database/postgres"
	"microservices_kafka_project/pkg/logger"
	"microservices_kafka_project/pkg/trx"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	cfg            *configs.Config
	dbPool         *pgxpool.Pool
	kafkaPub       *kafka.KafkaPublisher
	publicServer   *http.Server
	internalServer *http.Server
}

func NewContainer(cfg *configs.Config) *Container {
	return &Container{cfg: cfg}
}

func (c *Container) Start(ctx context.Context) error {
	var tx trx.TransactionManager
	var logger logger.Logger
	pool, err := database.NewPostgresPool(ctx, c.cfg.GetDBURL())
	if err != nil {
		return fmt.Errorf("подключение к бд: %w", err)
	}
	c.dbPool = pool

	c.kafkaPub = kafka.NewKafkaPublisher(c.cfg.Kafka.Brokers, "orders.v1.events")

	repo := orders.NewRepository(c.dbPool)
	orderService := orderService.NewService(tx, logger, repo, orders.Order{})
	app := orderApp.NewOrdersService(repo, c.kafkaPub, orderService)
	publickHandler := publichttp.NewPublicOrderHandler(app)
	internalHandler := internalhttp.NewPublicOrderHandler(app)

	publicR := gin.Default()
	publicR.Use(common.AuthMiddleware(c.cfg.Jwt.JwtSecret))
	{
		publicR.POST("/orders", publickHandler.CreateOrder)
		publicR.GET("/orders", publickHandler.GetAllUserOrders)
		publicR.GET("/orders/:id", publickHandler.GetOrderById)
	}

	internalR := gin.Default()
	internalR.PATCH("/internal/orders/:id/status", internalHandler.UpdateStatus)

	c.publicServer = &http.Server{Addr: ":" + c.cfg.HTTP.Port, Handler: publicR}
	c.internalServer = &http.Server{Addr: ":" + c.cfg.HTTP.InternalPort, Handler: internalR}

	go func() {
		log.Printf("Публичный Orders API запущен на порту %s", c.cfg.HTTP.Port)
		if err := c.publicServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Public server failed: %v", err)
		}
	}()

	go func() {
		log.Printf("Внутренний Orders API запущен на порту %s", c.cfg.HTTP.InternalPort)
		if err := c.internalServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Internal server failed: %v", err)
		}
	}()

	return nil
}

func (c *Container) Stop(ctx context.Context) error {
	log.Println("Остановка order сервиса")

	if c.dbPool != nil {
		c.dbPool.Close()
	}
	c.internalServer.Shutdown(ctx)
	c.publicServer.Shutdown(ctx)
	return nil
}
