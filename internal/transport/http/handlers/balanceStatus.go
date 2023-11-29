package handlers

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
	usecaseUsers "gopher-mart/internal/usecase/users"
	"net/http"
)

type BalanceStatusHandler struct {
	usecase balanceUsecase
}

func NewBalanceHandler(usecase balanceUsecase) *BalanceStatusHandler {
	return &BalanceStatusHandler{usecase: usecase}
}

type balanceUsecase interface {
	CheckBalance(ctx context.Context, user *users.User) (curBalance, withDrawn int, err error)
	usecaseUsers.UserContextUsecase
}
type responseBalance struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}

func (h *BalanceStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	curBalance, withDrawn, err := h.usecase.CheckBalance(r.Context(), user)
	if err != nil {
		log.Err(err).Send()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &responseBalance{
		Current:   curBalance,
		Withdrawn: withDrawn,
	}

	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Err(err).Send()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
