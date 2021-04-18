package handler

import (
	"net/http"
	"text/template"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	tmpl, errTmpl := template.ParseFiles("web/templates/register.html", "web/templates/base.html")

	if errTmpl != nil {
		http.Error(w, "ERROR 404", 404)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "ERROR 404", 404)
		return
	}
}