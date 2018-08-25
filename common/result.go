package common

import (
	"net/http"
	"encoding/json"
)

type Result struct{
	Status int `json:"status"`
	Msg interface{} `json:"msg"`
}
func OutputJson(w http.ResponseWriter, ret int, reason interface{}) {
	out := Result{ret, reason}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}
	w.Write(b)
}
