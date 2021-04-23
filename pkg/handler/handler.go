package handler

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Handler struct {
	
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func (h *Handler) InitRouters() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", h.Index)
	r.HandleFunc("/register/get", h.Register)
	r.HandleFunc("/register/post", h.RegisterPost)
	r.HandleFunc("/auth/get", h.AuthGet)
	r.HandleFunc("/auth/post", h.AuthPost)

	fileServer := http.FileServer(http.Dir("web/static"))
	r.PathPrefix("/").Handler(http.StripPrefix("/static/", fileServer))

	return r
}