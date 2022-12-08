package main

import (
	"bytes"
	"context"
	_ "embed"
	jsoniter "github.com/json-iterator/go"
	"github.com/lxzan/hasaki"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"mime"
	"net/http"
	"testing"
)

//go:embed logo.png
var pic []byte

type BaseResult struct {
	Code    *int   `json:"code"`
	Message string `json:"message"`
}

var afterFunc = func(ctx context.Context, resp *http.Response) (context.Context, error) {
	if resp.StatusCode != http.StatusOK {
		return ctx, hasaki.ErrUnexpectedStatusCode
	}
	typ, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if typ != "application/json" {
		return ctx, nil
	}

	rc, err := hasaki.NewReadCloser(resp.Body)
	if err != nil {
		return ctx, err
	}
	var result = BaseResult{}
	if err := jsoniter.Unmarshal(rc.Bytes(), &result); err != nil {
		return ctx, err
	}
	if result.Code == nil || *result.Code != 0 {
		return ctx, errors.New(result.Message)
	}
	resp.Body = rc
	return ctx, nil
}

func TestRequest(t *testing.T) {
	as := assert.New(t)
	hasaki.SetDefaultAfterFunc(afterFunc)

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
		}{
			Name: "caster",
			Age:  []int{1, 3, 5},
		}

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

	t.Run("500", func(t *testing.T) {
		err := hasaki.Post(baseURL + "/500").Send(nil).Err()
		as.Error(err)
	})
}
