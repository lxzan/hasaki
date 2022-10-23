package main

import (
	"crypto/md5"
	"encoding/hex"
	jsoniter "github.com/json-iterator/go"
	"github.com/lxzan/hasaki"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const addr = "127.0.0.1:9200"

var baseURL = "http://" + addr

type ResponseBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func WriteJson(writer http.ResponseWriter, code int, v interface{}) error {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	content, encodeErr := jsoniter.Marshal(v)
	if encodeErr != nil {
		return encodeErr
	}
	_, err := writer.Write(content)
	return err
}

func main() {
	http.HandleFunc("/json", func(writer http.ResponseWriter, request *http.Request) {
		var req = struct {
			Name string `json:"name"`
		}{}
		if err := jsoniter.NewDecoder(request.Body).Decode(&req); err != nil {
			WriteJson(writer, http.StatusOK, ResponseBody{Code: 1})
			return
		}
		WriteJson(writer, http.StatusOK, ResponseBody{Code: 0, Data: req})
	})

	http.HandleFunc("/form", func(writer http.ResponseWriter, request *http.Request) {
		b, _ := ioutil.ReadAll(request.Body)
		v, err := url.ParseQuery(string(b))
		if err != nil {
			WriteJson(writer, http.StatusOK, ResponseBody{Code: 1})
			return
		}

		WriteJson(writer, http.StatusOK, ResponseBody{Code: 0, Data: hasaki.Any{
			"name": v.Get("name"),
			"age":  v["age"],
		}})
	})

	http.HandleFunc("/upload", func(writer http.ResponseWriter, request *http.Request) {
		b, _ := ioutil.ReadAll(request.Body)
		hash := md5.New()
		hash.Write(b)
		code := hex.EncodeToString(hash.Sum(nil))
		WriteJson(writer, http.StatusOK, ResponseBody{
			Code: 0,
			Data: strings.ToUpper(code),
		})
	})

	http.HandleFunc("/500", func(writer http.ResponseWriter, request *http.Request) {
		WriteJson(writer, http.StatusBadGateway, ResponseBody{
			Code: 500,
		})
	})

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
