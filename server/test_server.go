package main

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type ResponseBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func WriteJson(writer http.ResponseWriter, code int, v interface{}) error {
	writer.WriteHeader(code)
	writer.Header().Set("Content-Type", "application/json")
	content, encodeErr := jsoniter.Marshal(v)
	if encodeErr != nil {
		return encodeErr
	}
	_, err := writer.Write(content)
	return err
}

func main() {
	http.HandleFunc("/200", func(writer http.ResponseWriter, request *http.Request) {
		WriteJson(writer, http.StatusOK, ResponseBody{
			Code:    http.StatusOK,
			Message: "success",
		})
	})

	http.HandleFunc("/400", func(writer http.ResponseWriter, request *http.Request) {
		WriteJson(writer, http.StatusOK, ResponseBody{
			Code:    http.StatusBadRequest,
			Message: "StatusBadRequest",
		})
	})

	http.ListenAndServe(":9000", nil)
}
