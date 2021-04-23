package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	repo "github.com/dankeka/goWebSocket/pkg/repository"
	"github.com/gorilla/sessions"
)

type RegisterStruct struct {
	Err string
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sessionRedisterChatGo")
	if err != nil {
		http.Error(w, "ERROR", http.StatusBadGateway)
		return
	}

	data := new(RegisterStruct)
	data.Err = fmt.Sprintf("%v", session.Values["errRegister"])

	if data.Err != "" {
		session.Options.MaxAge = -1

		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "ERROR", http.StatusInternalServerError)
			return
		}
	}

	tmpl, errTmpl := template.ParseFiles("web/templates/register.html", "web/templates/base.html")

	if errTmpl != nil {
		http.Error(w, "ERROR 404", 404)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "ERROR 404", 404)
		return
	}
}


type RegistePostStruct struct {
	Username  string `json:"username"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
}

func (h *Handler) RegisterPost(w http.ResponseWriter, r *http.Request) {
	var (
		session *sessions.Session
		db      *sql.DB
		err     error
	)

	registerPostData := new(RegistePostStruct)

	registerPostData.Username = r.FormValue("username") 
	registerPostData.Password1 = r.FormValue("password1")
	registerPostData.Password2 = r.FormValue("password2")

	if registerPostData.Password1 != registerPostData.Password2 {
		session, err = store.Get(r, "sessionRedisterChatGo")
		if err != nil {
			http.Error(w, "ERROR", http.StatusInternalServerError)
			return
		}

		session.Values["errRegister"] = "Пароли не совпадают!"
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "ERROR", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/register/get", http.StatusSeeOther)
		return
	}

	db, err = repo.Conn()
	if err != nil {
		http.Error(w, "ERROR", http.StatusBadGateway)
		return
	}
	defer db.Close()
	
	hashPassword, errHash := repo.HashFunc(registerPostData.Password1)
	if errHash != nil {
		http.Error(w, "ERROR", http.StatusBadGateway)
		return
	}

	row := db.QueryRow(
		"SELECT id FROM User WHERE name=$1",
		registerPostData.Username,
	)
	
	var userId sql.NullInt64
	row.Scan(&userId)

	if userId.Valid {
		session, err = store.Get(r, "sessionRedisterChatGo")
		if err != nil {
			http.Error(w, "ERROR", http.StatusInternalServerError)
			return
		}

		session.Values["errRegister"] = "Такой пользователь уже существует!"
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "ERROR", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/register/get", http.StatusSeeOther)
		return
	}

	_, err = db.Exec(
		"INSERT INTO User(name, password) VALUES ($1, $2)",
		registerPostData.Username,
		hashPassword,
	)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "ERROR", http.StatusBadGateway)
		return
	}

	http.Redirect(w, r, "/auth/get", http.StatusSeeOther)
}
