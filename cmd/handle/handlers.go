package handle

import (
	"encoding/json"
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
		writeErr(w, err)
		return
	}
	err = json.Unmarshal(all, &user)
	if err != nil {
		writeErr(w, err)
		return
	}

	err = h.srv.Login(r.Context(), user.Username, user.Password)
	if err != nil {
		writeErr(w, err)
		return
	}

	t, err := h.token.Create(token.UserToken{Username: user.Username}, 10*time.Minute)
	if err != nil {
		writeErr(w, err)
		return
	}
	resp := UserResponse{Token: t}
	bytes, err := json.Marshal(resp)
	if err != nil {
		writeErr(w, err)
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
		writeErr(w, err)
		return
	}
	err = json.Unmarshal(all, &user)
	if err != nil {
		writeErr(w, err)
		return
	}

	err = h.srv.Register(r.Context(), user.Username, user.Password)
	if err != nil {
		writeErr(w, err)
		return
	}
	t, err := h.token.Create(token.UserToken{Username: user.Username}, 10*time.Minute)
	if err != nil {
		writeErr(w, err)
		return
	}

	resp := UserResponse{Token: t}
	bytes, err := json.Marshal(resp)
	if err != nil {
		writeErr(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	log.Println("successful registration, user: ", user.Username)
}

func writeErr(w http.ResponseWriter, err error) {
	log.Println("error :", err)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}
