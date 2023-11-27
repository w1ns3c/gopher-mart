package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/orders"
	ordersUsecase "gopher-mart/internal/usecase/orders"
	"net/http"
)

type orderStatusHandler struct {
	usecase getOrderStatusUsecase
}

func NewOrderStatusHandler(usecase getOrderStatusUsecase) *orderStatusHandler {
	return &orderStatusHandler{usecase: usecase}
}

type getOrderStatusUsecase interface {
	ordersUsecase.OrderValidator
	CheckOrderStatus(ctx context.Context, orderNumber string) (order *orders.Order, err error)
}

type orderResponse struct {
	ID      string             `json:"order"`
	Accrual uint64             `json:"accrual"`
	Status  orders.OrderStatus `json:"status"`
}

func (h *orderStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	orderNum := chi.URLParam(r, "number")
	if !h.usecase.ValidateOrderFormat(r.Context(), orderNum) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	order, err := h.usecase.CheckOrderStatus(r.Context(), orderNum)
	if err != nil {
		switch err {
		case errors.ErrOrderNotFound:
			w.WriteHeader(http.StatusNoContent)
		}
		return
	}

	var resp orderResponse
	resp.ID = order.ID
	resp.Accrual = order.Cashback
	resp.Status = order.Status

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
