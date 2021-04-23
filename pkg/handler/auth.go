package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"golang.org/x/crypto/bcrypt"

	"github.com/dankeka/goWebSocket/pkg/handler/auth"
	repo "github.com/dankeka/goWebSocket/pkg/repository"
	"github.com/dankeka/goWebSocket/pkg/types"
	"github.com/gorilla/sessions"
)

type AuthStruct struct {
	Err string
}

func (h *Handler) AuthGet(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sessionAuthChatGo")
	if err != nil {
		http.Error(w, "ERROR", http.StatusBadGateway)
		return
	}

	data := new(AuthStruct)
	data.Err = fmt.Sprintf("%v", session.Values["authErr"])

	if data.Err != "" {
		session.Options.MaxAge = -1

		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "ERROR", http.StatusInternalServerError)
			return
		}
	}

	tmpl, errTmpl := template.ParseFiles("web/templates/auth.html", "web/templates/base.html")

	if errTmpl != nil {
		http.Error(w, "ERROR 404", 404)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "ERROR 404", 404)
		return
	}
}

func (h *Handler) AuthPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	db, err := repo.Conn()
	if err != nil {
		http.Error(w, "ERROR", http.StatusBadGateway)
		return
	}
	defer db.Close()

	row := db.QueryRow(
		"SELECT password FROM User WHERE name=$1",
		username,
	)

	var passwordDB []byte
	err = row.Scan(&passwordDB)
	
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "ERROR", http.StatusBadGateway)
		return
	} else if err == sql.ErrNoRows {
		var session *sessions.Session
		session, err = store.Get(r, "sessionAuthChatGo")
		if err != nil {
			http.Error(w, "ERROR", http.StatusBadGateway)
			return
		}

		session.Values["authErr"] = "Такого пользователя не существует!"

		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "ERROR", http.StatusBadGateway)
			return
		}
	}
	
	err = bcrypt.CompareHashAndPassword(passwordDB, []byte(password+os.Getenv("SALT")))
	
	if err != nil {
		var session *sessions.Session
		session, err = store.Get(r, "sessionAuthChatGo")
		if err != nil {
			http.Error(w, "ERROR", http.StatusBadGateway)
			return
		}

		session.Values["authErr"] = "Пароли не совпадают!"

		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "ERROR", http.StatusBadGateway)
			return
		}
	}

	row = db.QueryRow(
		"SELECT id FROM User WHERE name=$1 AND password=$2",
		username,
		passwordDB,
	)

	var userId sql.NullInt64
	err = row.Scan(&userId)
	if err != nil {
		http.Error(w, "ERROR", http.StatusBadGateway)
		return
	}

	if userId.Int64 == 0 || !userId.Valid {
		var session *sessions.Session
		session, err = store.Get(r, "sessionAuthChatGo")
		if err != nil {
			http.Error(w, "ERROR", http.StatusBadGateway)
			return
		}

		session.Values["authErr"] = "Такого пользователя не существует!"

		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "ERROR", http.StatusBadGateway)
			return
		}
	}

	authUser := types.User {
		ID: uint(userId.Int64),
		Name: username,
	}

	var newJwtToken string
	newJwtToken, err = auth.GenerateJWT(authUser)
	if err != nil {
		http.Error(w, "ERROR", http.StatusBadGateway)
		return
	}

	var session *sessions.Session
	session, err = store.Get(r, "sessionAuthChatGo")
	if err != nil {
		http.Error(w, "ERROR", http.StatusBadGateway)
		return
	}

	session.Values["jwtToken"] = newJwtToken

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "ERROR", http.StatusBadGateway)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}