package handler

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/dankeka/goWebSocket/pkg/handler/auth"
	"github.com/dgrijalva/jwt-go"
)

type MainStruct struct {
	Auth      bool
	UserData jwt.MapClaims
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	data := new(MainStruct)

	session, err := store.Get(r, "sessionAuthChatGo")
	if err != nil {
		data.Auth = false
		data.UserData = jwt.MapClaims{}
	} else {
		jwtTokenSession := session.Values["jwtToken"]

		var parseJwt jwt.MapClaims
		parseJwt, err = auth.ParseJWT(fmt.Sprintf("%v", jwtTokenSession))
		if err != nil {
			data.Auth = false
			data.UserData = jwt.MapClaims{}
		}

		data.Auth = true
		data.UserData = parseJwt
	}

	tmpl, errTmpl := template.ParseFiles("web/templates/index.html", "web/templates/base.html")

	if errTmpl != nil {
		http.Error(w, "ERROR 404", 404)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "ERROR 404", 404)
		return
	}
}