package main

import (
	"bytes"
	_ "embed"
	jsoniter "github.com/json-iterator/go"
	"github.com/lxzan/hasaki"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"testing"
)

//go:embed logo.png
var pic []byte

type BaseResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var checker hasaki.ErrorChecker = func(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected status code")
	}
	typ, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if typ != "application/json" {
		return nil
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result = BaseResult{}
	if err := jsoniter.Unmarshal(b, &result); err != nil {
		return err
	}
	if result.Code != 0 {
		return errors.New(result.Message)
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(b))
	return nil
}

func TestRequest(t *testing.T) {
	as := assert.New(t)
	hasaki.SetGlobalErrorChecker(checker)

	t.Run("json", func(t *testing.T) {
		var result = struct {
			BaseResult
			Data struct {
				Name string `json:"name"`
			} `json:"data"`
		}{}

		err := hasaki.
			Post(baseURL + "/json").
			SetEncoder(hasaki.JsonEncoder).
			Send(hasaki.Any{
				"name": "caster",
			}).
			BindJSON(&result)
		as.NoError(err)
		as.Equal("caster", result.Data.Name)
	})

	t.Run("form", func(t *testing.T) {
		var result = struct {
			BaseResult
			Data struct {
				Name string   `json:"name"`
				Age  []string `json:"age"`
			} `json:"data"`
		}{}

		var req = struct {
			Name string `form:"name"`
			Age  []int  `form:"age"`
		}{}

		err := hasaki.
			Post(baseURL + "/form").
			SetEncoder(hasaki.FormEncoder).
			Send(&req).
			BindJSON(&result)
		as.NoError(err)
		as.Equal("caster", result.Data.Name)
		as.ElementsMatch([]string{"1", "3", "5"}, result.Data.Age)
	})

	t.Run("error", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"name":1}`))
		err := hasaki.
			Post(baseURL + "/json").
			Send(body).
			Err()
		as.Error(err)
	})

	t.Run("file", func(t *testing.T) {
		var result = struct {
			BaseResult
			Data string `json:"data"`
		}{}

		err := hasaki.
			Put(baseURL + "/upload").
			SetHeader(hasaki.H{"Content-Type": hasaki.ContentType_STREAM.String()}).
			Send(bytes.NewBuffer(pic)).
			BindJSON(&result)
		as.NoError(err)
		as.Equal("87D649E1E74817F8A302355EA11191F4", result.Data)
	})
}
