package service

import (
	"WBTech0/internal/entity"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
)

type NatsRepo interface {
	Subscribe(subject, queueGroup string, handler stan.MsgHandler, options ...stan.SubscriptionOption) (stan.Subscription, error)
}

type CacheRepo interface {
	Get(key string) (*entity.Order, bool)
	Set(value *entity.Order)
}

type NatsService struct {
	natsRepo  NatsRepo
	orderRepo OrderRepo
	cacheRepo CacheRepo
}

func NewNatsService(natsRepo NatsRepo, cacheRepo CacheRepo, orderRepo OrderRepo) *NatsService {
	return &NatsService{natsRepo: natsRepo, cacheRepo: cacheRepo, orderRepo: orderRepo}
}

func (s *NatsService) StartListening(subject, queueGroup string) (stan.Subscription, error) {
	return s.natsRepo.Subscribe(subject, queueGroup, func(m *stan.Msg) {
		if err := s.handleMessage(m); err != nil {
			log.Printf("Failed to handle message with sequence number %d: %v", m.Sequence, err)
			log.Printf("Message content: %s", m.Data)
		}
	})
}

func (s *NatsService) handleMessage(m *stan.Msg) error {
	var newOrder entity.Order
	err := json.Unmarshal(m.Data, &newOrder)
	if err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}
	if err := s.orderRepo.AddOrder(context.Background(), newOrder); err != nil {
		return fmt.Errorf("failed to add order to repo: %w", err)
	}
	s.cacheRepo.Set(&newOrder)
	return nil
}
