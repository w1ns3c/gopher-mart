package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"gopher-mart/internal/domain/accruals"
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

type OrderResponse struct {
	ID      string                           `json:"order"`
	Accrual float64                          `json:"accrual,omitempty"`
	Status  accruals.AccrualSystemRegistered `json:"status"`
}

func (h *orderStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	orderNum := chi.URLParam(r, "number")
	if !h.usecase.ValidateOrderFormat(orderNum) {
		w.WriteHeader(http.StatusNoContent)
		log.Error().Msg("number not in accrual system")
		return
	}

	order, err := h.usecase.CheckOrderStatus(r.Context(), orderNum)
	if err != nil {
		log.Err(err).Send()
		switch err {
		case errors.ErrOrderNotFound:
			w.WriteHeader(http.StatusNoContent)
		}
		return
	}

	var resp OrderResponse
	resp.ID = order.ID
	resp.Accrual = order.Cashback
	resp.Status = accruals.AccrualSystemRegistered(order.Status)

	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Err(err).Send()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
