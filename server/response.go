package coditra

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func errorResponse(w http.ResponseWriter, code int, err error) {
	jsonResponse(w, code, struct {
		Error string `json:"error"`
	}{err.Error()})
}

func jsonResponse(w http.ResponseWriter, code int, payload any) {
	res, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(res)
	if err != nil {
		fmt.Printf("ERROR: writing HTTP response: %v", err)
	}
}
