package auth

import (
	"context"
	"log"
	"microservices_kafka_project/configs"
	"microservices_kafka_project/internal/auth/app"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

const confDir = "./configs/main.yaml"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Предупреждение: .env файл не найден")
	}
	cfg, err := configs.NewConfig(confDir)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	cnt := app.New(cfg)
	if err := cnt.Start(context.Background()); err != nil {
		log.Fatalf("Start error: %s", err.Error())
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	s := <-interrupt
	log.Printf("app - Start - signal: " + s.String())

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := cnt.Stop(shutdownCtx); err != nil {
		log.Printf("Stop error: %v", err)
	}
}
