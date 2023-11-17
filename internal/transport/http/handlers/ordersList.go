package handlers

import (
	"context"
	"encoding/json"
	"gopher-mart/internal/domain/orders"
	"net/http"
)

type ordersGetHandler struct {
	usecase ordersUsecase
}

type ordersUsecase interface {
	// TODO orders MUST be sorted by date, should it be on repo level?
	ListOrders(ctx context.Context) (orders []orders.Order, err error)
}

type ordersResponse struct {
	ID      string             `json:"number"`
	Status  orders.OrderStatus `json:"status"`
	Accrual uint64             `json:"accrual,omitempty"` // accrual
	Date    string             `json:"uploaded_at"`
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
	l := len(orders)
	if l == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := make([]ordersResponse, l)
	for id, order := range orders {
		resp[id].ID = order.ID
		resp[id].Status = order.Status
		resp[id].Accrual = order.Cashback
		resp[id].Date = order.Date.Format("2006-01-02T15:04:05-07:00")
	}

	body, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
