package cache

import (
	"WBTech0/internal/entity"
	"context"
	"sync"
)

type OrderRepo interface {
	GetAllOrders(ctx context.Context) ([]entity.Order, error)
}

type CacheRepo struct {
	mu    sync.RWMutex
	cache map[string]*entity.Order
}

func NewOrderCache(orderRepo OrderRepo) *CacheRepo {
	orderCache := &CacheRepo{cache: make(map[string]*entity.Order)}
	orders, _ := orderRepo.GetAllOrders(context.Background())
	for _, v := range orders {
		orderCache.Set(&v)
	}
	return orderCache
}

func (oc *CacheRepo) Get(key string) (*entity.Order, bool) {
	oc.mu.RLock()
	defer oc.mu.RUnlock()
	value, ok := oc.cache[key]
	return value, ok
}

func (oc *CacheRepo) Set(order *entity.Order) {
	oc.mu.Lock()
	defer oc.mu.Unlock()
	oc.cache[order.OrderUID] = order
}
