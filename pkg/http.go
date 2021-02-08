package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// seriable the response with JSON
// @param w      response writer
// @param data   data to respond
func WriteJsonOfResponce(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json;charset=utf-8")

	b, err := json.Marshal(data)
	if err != nil {
		Loge(err.Error())
	}
	fmt.Fprint(w, string(b))
}

// set 400 for the response
// @param msg    message of 400 error
func BadRequest(w http.ResponseWriter, msg string) {
	Logd(msg)
	w.WriteHeader(400)
	fmt.Fprint(w, msg)
}

// 404
func NotFound(w http.ResponseWriter, msg string) {
	Logd(msg)
	w.WriteHeader(404)
	fmt.Fprint(w, msg)
}

// CORS
func Cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("content-type", "application/json;charset=UTF-8")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		f(w, r)
	}
}
