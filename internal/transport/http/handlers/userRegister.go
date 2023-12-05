package handlers

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
	"strings"

	"net/http"
)

type RegisterHandler struct {
	usecase registerUsecase
}

func NewRegisterHandler(registerUsecase registerUsecase) *RegisterHandler {
	return &RegisterHandler{usecase: registerUsecase}
}

type registerUsecase interface {
	RegisterUser(ctx context.Context, user *users.User) (cookie *http.Cookie, err error)
}

type reqisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	//ConfirmPassword string `json:"confirm"`
	ConfirmPassword string `json:"-"`
}

func (req *reqisterRequest) ToUserWithConfirm() (*users.User, error) {
	if req.Password != req.ConfirmPassword {
		return nil, errors.ErrConfirmPassword
	}
	return users.NewUser(req.Login, req.Password), nil
}

func (req *reqisterRequest) ToUser() (*users.User, error) {
	return users.NewUser(req.Login, req.Password), nil
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Err(errors.ErrMethodNotAllowed).Send()
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("content-type") != "application/json" {
		log.Err(errors.ErrWrongContentType).Send()
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
		log.Err(err).Send()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookie, err := h.usecase.RegisterUser(r.Context(), user)
	if err != nil {
		log.Err(err).Send()
		switch err {
		case bcrypt.ErrPasswordTooLong:
			w.Write([]byte("too long password"))
			w.WriteHeader(http.StatusInternalServerError)
		default:
			if strings.Contains(err.Error(), "duplicate key value") {
				w.Write([]byte("login is already used"))
				w.WriteHeader(http.StatusConflict)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

		}
		return
	}

	w.Header().Set("Set-Cookie", cookie.String())
	w.Write([]byte("OK"))
}
