package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"gopher-mart/internal/domain/accruals"
	"gopher-mart/internal/domain/orders"
	"gopher-mart/internal/transport/http/handlers"
	"gopher-mart/internal/transport/http/middlewares"
	"gopher-mart/internal/utils"
	"net/http"
)

type withdrawUsecae struct {
	orders []handlers.OrderResponse
}

func (w withdrawUsecae) ValidateOrderFormat(orderNumber string) bool {
	return utils.LuhnValidator(orderNumber)
}

func (w withdrawUsecae) CheckOrderStatus(ctx context.Context,
	orderNumber string) (order *orders.Order, err error) {

	for _, order := range w.orders {
		if orderNumber == order.ID {
			return &orders.Order{
				ID:       order.ID,
				Cashback: order.Accrual,
				Status:   orders.OrderStatus(order.Status),
			}, nil
		}
	}

	return nil, fmt.Errorf("no such order number")
}

func main() {
	addr := `localhost:9000`

	orders := []handlers.OrderResponse{
		{
			ID:     "18",
			Status: accruals.PROCESSING,
		},
		{
			ID:      "26",
			Status:  accruals.PROCESSED,
			Accrual: 100,
		},
		{
			ID:     "109",
			Status: accruals.INVALID,
		},
		{
			ID:      "901",
			Status:  accruals.PROCESSED,
			Accrual: 200,
		},
		{
			ID:     "9084",
			Status: accruals.PROCESSING,
		},
		{
			ID:     "234567",
			Status: accruals.PROCESSING,
		},
		{
			ID:     "123455",
			Status: accruals.REGISTERED,
		},
		{
			ID:     "1234566",
			Status: accruals.INVALID,
		},
	}
	usecase := &withdrawUsecae{orders: orders}

	router := chi.NewRouter()
	handler := handlers.NewOrderStatusHandler(usecase)
	router.Use(middlewares.LoggingMiddleware)
	ddos := middlewares.NewDDOSMiddleware(10)
	router.Use(ddos.DDOSMiddleware)
	router.Get("/api/orders/{number}", handler.ServeHTTP)

	fmt.Println("[i] Server listen on:", addr)
	http.ListenAndServe(addr, router)

}
