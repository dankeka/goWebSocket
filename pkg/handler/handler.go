package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	
}

func (h *Handler) InitRouters() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", h.Index)
	r.HandleFunc("/register/get", h.Register)
	r.HandleFunc("/register/post", nil)

	fileServer := http.FileServer(http.Dir("web/static"))
	r.PathPrefix("/").Handler(http.StripPrefix("/static/", fileServer))

	return r
}