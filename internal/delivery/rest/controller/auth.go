package controller

import (
	"encoding/json"
	"github.com/minishop/internal/domain"
	"net/http"
	"time"
)

type AuthController struct {
	authUsecase domain.AuthUsecase
}

func NewAuthController(authUsecase domain.AuthUsecase) *AuthController {
	return &AuthController{authUsecase: authUsecase}
}

func (a *AuthController) login(w http.ResponseWriter, r *http.Request) {
	var loginRequest domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	auth, err := a.authUsecase.Login(r.Context(), &loginRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	bt, err := json.Marshal(auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *AuthController) logout(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 8)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("ok!"))
}

func (a *AuthController) Handle() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", a.login)
	mux.HandleFunc("/logout", a.logout)

	return mux
}
