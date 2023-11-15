package handlers

import (
	"encoding/json"
	dusers "gopher-mart/internal/domain/users"
	"gopher-mart/internal/usecase/users"
	"net/http"
)

type RegisterHandler struct {
	registerUsecase users.UserUsecase
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

	var user dusers.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.registerUsecase.RegisterUser(user)
	if err != nil {
		w.Write([]byte("login is already used"))
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.Write([]byte("OK"))
}
