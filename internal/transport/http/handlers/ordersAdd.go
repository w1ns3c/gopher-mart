package handlers

import (
	"context"
	"gopher-mart/internal/domain/orders"
	"gopher-mart/internal/domain/users"
	"net/http"
)

type OrdersAddHandler struct {
	usecase ordersAddUsecase
}

type ordersAddUsecase interface {
	AddOrder(ctx context.Context, order *orders.Order) error
	ValidateOrder(ctx context.Context, string) (u *users.User, err error)
}

type ordersAddRequest struct {
	orderNumber stringf
}

func (h *OrdersAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req ordersAddRequest
	req.cookie, err := r.Cookie()

}
