package http

import (
	"WBTech0/internal/entity"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type OrderService interface {
	GetOrderWithDetailsById(ctx context.Context, orderUID string) (entity.Order, error)
}

type controller struct {
	orderService OrderService
}

func newController(orderService OrderService) *controller {
	return &controller{
		orderService: orderService,
	}
}

func (c *controller) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("zawel")
	orderID := r.URL.Query().Get("id")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	order, err := c.orderService.GetOrderWithDetailsById(r.Context(), orderID)

	if err != nil {
		http.Error(w, `{"data":"Order not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
