package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func WrapError(w http.ResponseWriter, err error) {
	WrapErrorWithStatus(w, err, http.StatusBadRequest)
}

func WrapErrorWithStatus(w http.ResponseWriter, err error, httpStatus int) {
	var m = map[string]string{
		"result": "error",
		"data":   err.Error(),
	}

	res, _ := json.Marshal(m)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.Header().Set("X-Content-Type-Options", "nosniff") //даем понять что ответ приходит в формате json
	w.WriteHeader(httpStatus)
	fmt.Fprintln(w, string(res))
}

func WrapOK(w http.ResponseWriter, m map[string]interface{}) {
	res, _ := json.Marshal(m)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(res))
}
func WrapOKImage(w http.ResponseWriter, m string) {
	fileBytes, err := os.ReadFile(m)
	if err != nil {
		WrapError(w, fmt.Errorf("файл не найден"))
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}
