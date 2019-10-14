package server

import (
	"log"
	"encoding/json"
	"net/http"
)

func WriteJSON(v interface{}, s int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(v)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(s)
	w.Write(b)
}
