package controllers

import (
	"html/template"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/home.html", "views/common/head.tpl")
	if pageNotFound(w, err) {
		return
	}

	data := make(map[string]interface{})
	_, err = r.Cookie("Uname")
	if err != nil {
		data["isLogin"] = false
	} else {
		data["isLogin"] = true
	}

	t.Execute(w, data)
}
