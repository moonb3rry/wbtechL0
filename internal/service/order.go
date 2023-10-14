package service

import (
	"WBTech0/internal/entity"
	"context"
	"errors"
)

type OrderRepo interface {
	AddOrder(ctx context.Context, order entity.Order) error
	GetAllOrders(ctx context.Context) ([]entity.Order, error)
}

type OrderService struct {
	orderRepo OrderRepo
	cacheRepo CacheRepo
}

func NewOrderService(or OrderRepo, cr CacheRepo) *OrderService {
	return &OrderService{
		orderRepo: or,
		cacheRepo: cr,
	}
}

func (s *OrderService) GetOrderWithDetailsById(ctx context.Context, orderUID string) (entity.Order, error) {
	order, ok := s.cacheRepo.Get(orderUID)
	OrderNotFoundInCacheError := errors.New("Order not found in cache")
	if ok != true {
		return entity.Order{}, OrderNotFoundInCacheError
	}
	return *order, nil
}

func (s *OrderService) GetAllOrders(ctx context.Context) ([]entity.Order, error) {
	orders, err := s.orderRepo.GetAllOrders(ctx)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
