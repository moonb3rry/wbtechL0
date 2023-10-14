package http

import "github.com/gorilla/mux"

func NewOrderController(s OrderService) *mux.Router {
	controller := newController(s)
	router := mux.NewRouter()
	router.HandleFunc("/order", controller.GetOrderHandler).Methods("GET")
	return router
}
