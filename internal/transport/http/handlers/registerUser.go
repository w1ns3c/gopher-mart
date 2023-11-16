package handlers

import (
	"context"
	"encoding/json"
	"gopher-mart/internal/domain/localerrors"
	"gopher-mart/internal/domain/users"

	"net/http"
)

type RegisterHandler struct {
	registerUsecase registerUsecase
}

type registerUsecase interface {
	RegisterUser(ctx context.Context, user *users.User) error
}

type reqisterRequest struct {
	Login           string `json:"login"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm"`
}

func (req *reqisterRequest) ToUser() (*users.User, error) {
	if req.Password != req.ConfirmPassword {
		return nil, localerrors.ErrConfirmPassword
	}
	return users.NewUser(req.Login, req.Password), nil
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req reqisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := req.ToUser()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.registerUsecase.RegisterUser(r.Context(), user)
	if err != nil {
		w.Write([]byte("login is already used"))
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.Write([]byte("OK"))
}
