package response

import (
	"encoding/json"
	"net/http"
)

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, status int, body Body) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, Body{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, Body{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, Body{
		Code:    status,
		Message: message,
	})
}
