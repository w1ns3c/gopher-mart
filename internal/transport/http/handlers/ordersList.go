package handlers

import (
	"context"
	"gopher-mart/internal/domain/orders"
	"net/http"
	"time"
)

type ordersGetHandler struct {
	usecase ordersUsecase
}

type ordersUsecase interface {
	// orders MUST be sorted by date
	ListOrders(ctx context.Context) (orders []orders.Order, err error)
}

type ordersRequest struct {
	ID       string             `json:"number"`
	Status   orders.OrderStatus `json:"status"`
	Cashbash uint64             `json:"accrual,omitempty"` // accrual
	Date     time.Time          `json:"uploaded_at"`
}

func (h *ordersGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	orders, err := h.usecase.ListOrders(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// orders not found
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// TODO orders or ordersRequest Here
}
