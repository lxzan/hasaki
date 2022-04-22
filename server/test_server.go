package main

import (
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

var response = map[string]interface{}{
	"success": true,
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
	http.HandleFunc("/p1", func(writer http.ResponseWriter, request *http.Request) {
		WriteJson(writer, 200, response)
	})

	http.HandleFunc("/p2", func(writer http.ResponseWriter, request *http.Request) {
		WriteJson(writer, 400, response)
	})

	http.HandleFunc("/p4", func(writer http.ResponseWriter, request *http.Request) {
		WriteJson(writer, 500, response)
	})

	http.ListenAndServe(":9000", nil)
}
