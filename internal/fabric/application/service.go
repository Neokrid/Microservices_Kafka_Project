package fabric

import (
	"context"
	"log"
	"math/rand"
	"microservices_kafka_project/internal/fabric/infrastructure/http"
	"microservices_kafka_project/internal/fabric/infrastructure/kafka"
	"time"
)

type FabricService struct {
	ordersClient *http.OrderClient
}

func NewFabricService(client *http.OrderClient) *FabricService {
	return &FabricService{
		ordersClient: client,
	}
}

func (s *FabricService) HandleOrder(ctx context.Context, event kafka.OrderEvent) {

	if event.EventType != "OrderCreated" {
		return
	}

	log.Printf("fabric начал сборку заказа %s (Items: %v)", event.Payload.ID, event.Payload.Items)

	duration := time.Duration(rand.Intn(11)+5) * time.Second

	select {
	case <-time.After(duration):
		log.Printf("заказ %s готов! (заняло %v)", event.Payload.ID, duration)

		err := s.ordersClient.UpdateStatus(ctx, event.Payload.ID, "completed")
		if err != nil {
			log.Printf("ошибка обратного вызова для %s: %v", event.Payload.ID, err)
		} else {
			log.Printf("статус заказа %s успешно обновлен на 'completed'", event.Payload.ID)
		}
	case <-ctx.Done():
		log.Printf("работа над заказом %s прервана", event.Payload.ID)
	}
}
