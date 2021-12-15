package handle

import (
	"encoding/json"
	"errors"
	"github.com/mahjadan/login/pkg/repository"
	"github.com/mahjadan/login/pkg/service"
	"github.com/mahjadan/login/pkg/token"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func New(service service.UserService, maker token.Maker) Handler {
	return Handler{
		srv:   service,
		token: maker,
	}
}

type Handler struct {
	srv   service.UserService
	token token.Maker
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user UserRequest
	all, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErr(w, NewBadRequestResponse(err.Error()))
		return
	}
	err = json.Unmarshal(all, &user)
	if err != nil {
		writeErr(w, NewBadRequestResponse(err.Error()))
		return
	}

	err = h.srv.Login(r.Context(), user.Username, user.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserOrPassword) {
			writeErr(w, UnauthorizedLoginResponse)
			return
		}
		log.Println("login, ", err)
		writeErr(w, InternalServerErrorResponse)
		return
	}

	t, err := h.token.Create(token.UserToken{Username: user.Username}, 10*time.Minute)
	if err != nil {
		log.Println("create token, ", err)
		writeErr(w, InternalServerErrorResponse)
		return
	}
	resp := UserResponse{Token: t}
	bytes, err := json.Marshal(resp)
	if err != nil {
		log.Println("marshal token, ", err)
		writeErr(w, InternalServerErrorResponse)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	log.Println("successful logged in, user: ", user.Username)
}

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user UserRequest
	all, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErr(w, NewBadRequestResponse(err.Error()))
		return
	}
	err = json.Unmarshal(all, &user)
	if err != nil {
		writeErr(w, NewBadRequestResponse(err.Error()))
		return
	}

	err = h.srv.Register(r.Context(), user.Username, user.Password)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			writeErr(w, ConflictRequestErrorResponse)
			return
		}
		log.Println("register, ", err)
		writeErr(w, InternalServerErrorResponse)
		return
	}
	t, err := h.token.Create(token.UserToken{Username: user.Username}, 10*time.Minute)
	if err != nil {
		log.Println("create token, ", err)
		writeErr(w, InternalServerErrorResponse)
		return
	}

	resp := UserResponse{Token: t}
	bytes, err := json.Marshal(resp)
	if err != nil {
		log.Println("marshal token, ", err)
		writeErr(w, InternalServerErrorResponse)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	log.Println("successful registration, user: ", user.Username)
}

func writeErr(w http.ResponseWriter, errorResponse HTTPErrorResponse) {
	log.Println("error :", errorResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorResponse.HTTPStatusCode)
	w.Write(errorResponse.ToJSON())
}
