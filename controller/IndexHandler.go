package controller

import (
	"net/http"
	"html/template"
	"github.com/cihub/seelog"
)

func IndexHandler(w http.ResponseWriter,r *http.Request)  {
	t, err := template.ParseFiles("template/html/index.html")
	if (err != nil) {
		seelog.Infof(err.Error())
	}
	t.Execute(w, nil)
}