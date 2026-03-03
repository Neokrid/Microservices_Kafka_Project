package main

import (
	"context"
	"log"
	"microservices_kafka_project/configs"
	fabric "microservices_kafka_project/internal/fabric/application"
	"microservices_kafka_project/internal/fabric/infrastructure/http"
	"microservices_kafka_project/internal/fabric/infrastructure/kafka"
	"os"
	"os/signal"
	"syscall"
)

const confDir = "./configs/main.yaml"

func main() {
	cfg, err := configs.NewConfig(confDir, ".env.fabric")
	if err != nil {
		log.Fatalf("ошибка конфига: %s", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	brokers := cfg.Kafka.Brokers
	topic := cfg.Kafka.Topic
	groupID := "fabric-service-group"
	ordersBaseURL := cfg.Kafka.OrderPath

	ordersClient := http.NewOrderClient(ordersBaseURL)
	fabricService := fabric.NewFabricService(ordersClient)
	consumer := kafka.NewKafkaConsumer(brokers, topic, groupID)

	go func() {

		consumer.Listen(ctx, fabricService.HandleOrder)
	}()

	<-ctx.Done()
	log.Println(" Сервис Fabric завершает работу...")
	consumer.Close()
}
