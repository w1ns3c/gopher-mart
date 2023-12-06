package handlers

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
	userUsecase "gopher-mart/internal/usecase/users"
	"net/http"
	"time"
)

type withdrawalsHandler struct {
	usecase withdrawalsUsecase
}

func NewWithdrawalsHandler(usecase withdrawalsUsecase) *withdrawalsHandler {
	return &withdrawalsHandler{usecase: usecase}
}

type withdrawalsUsecase interface {
	GetUserWithdrawals(ctx context.Context, user *users.User) (wd []withdraws.Withdraw, err error)
	userUsecase.UserContextUsecase
}
type responseWithdrawls struct {
	OrderID string    `json:"order"`
	Sum     float64   `json:"sum"`
	Date    time.Time `json:"-"`
}

func (r *responseWithdrawls) MarshalJSON() ([]byte, error) {
	type Alias responseWithdrawls
	return json.Marshal(&struct {
		Date string `json:"processed_at"`
		*Alias
	}{
		r.Date.Format(time.RFC3339),
		(*Alias)(r),
	})
}

func (h *withdrawalsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := h.usecase.CheckUserInContext(r.Context())
	if err != nil {
		log.Err(err).Send()
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodGet {
		log.Err(errors.ErrMethodNotAllowed).Send()
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	wd, err := h.usecase.GetUserWithdrawals(r.Context(), user)
	if err != nil {
		log.Err(err).Send()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(wd) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response := make([]responseWithdrawls, len(wd))
	for id, one := range wd {
		response[id].OrderID = one.OrderID
		response[id].Sum = one.Sum
		response[id].Date = one.Date
	}

	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Err(err).Send()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
