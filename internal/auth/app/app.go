package app

import (
	"fmt"
	"log"

	"microservices_kafka_project/configs"
	database "microservices_kafka_project/pkg/database/postgres"
	"microservices_kafka_project/pkg/logger"
	"microservices_kafka_project/pkg/trx"

	app "microservices_kafka_project/internal/auth/application"
	"microservices_kafka_project/internal/auth/domain/services/token"
	auth "microservices_kafka_project/internal/auth/domain/services/user"
	tokenRepo "microservices_kafka_project/internal/auth/infrastructure/token"
	"microservices_kafka_project/internal/auth/infrastructure/user"
	authHttp "microservices_kafka_project/internal/auth/ports/http"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
)

type Container struct {
	cfg    *configs.Config
	server *http.Server
	dbPool *pgxpool.Pool
}

func New(cfg *configs.Config) *Container {
	return &Container{
		cfg: cfg,
	}
}

func (c *Container) Start(ctx context.Context) error {
	var tx trx.TransactionManager
	var logger logger.Logger
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.cfg.Postgres.Username, c.cfg.Postgres.Password, c.cfg.Postgres.Host, c.cfg.Postgres.Port, c.cfg.Postgres.DBName)

	pool, err := database.NewPostgresPool(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("подключение к бд: %w", err)
	}
	c.dbPool = pool

	uRepo := user.NewRepository(c.dbPool)
	tRepo := tokenRepo.NewRepository(c.dbPool)
	uService := auth.NewService(tx, logger, uRepo)
	tService := token.NewService(c.cfg.Jwt.RefreshTTL, c.cfg.Jwt.AccessTTL, c.cfg.Jwt.JwtSecret, tRepo)
	app := app.NewAuthService(uRepo, tx, logger, uService, tService)
	handler := authHttp.NewAuthHandler(app)

	r := gin.Default()
	r.POST("/register", handler.SignUp)
	r.POST("/login", handler.SignIn)

	c.server = &http.Server{
		Addr:    ":" + c.cfg.HTTP.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Auth сервис запущен на порту %s", c.cfg.HTTP.Port)
		if err := c.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	return nil
}

func (c *Container) Stop(ctx context.Context) error {
	log.Println("Остановка auth сервиса")

	if c.dbPool != nil {
		c.dbPool.Close()
	}

	return c.server.Shutdown(ctx)
}
