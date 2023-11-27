package handlers

import (
	"context"
	"encoding/json"
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
	Sum     uint64    `json:"sum"`
	Date    time.Time `json:"-"`
}

func (r *responseWithdrawls) MarshalJSON() ([]byte, error) {
	type Alias responseWithdrawls
	return json.Marshal(&struct {
		Date string `json:"processed_at"`
		*Alias
	}{
		r.Date.Format("2006-01-02T15:04:05-07:00"),
		(*Alias)(r),
	})
}

func (h *withdrawalsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := h.usecase.CheckUserInContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	wd, err := h.usecase.GetUserWithdrawals(r.Context(), user)
	if err != nil {
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

	body, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(body)

}
