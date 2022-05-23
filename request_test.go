package hasaki

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type responseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (this responseBody) Error() string {
	return this.Message
}

func TestNewRequest(t *testing.T) {
	const baseurl = "http://localhost:9000"

	var checker ErrorChecker = func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			return errors.New("StatusCodeError")
		}

		var buf = bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, resp.Body); err != nil {
			return err
		}
		resp.Body.Close()
		var body = &responseBody{}
		jsoniter.Unmarshal(buf.Bytes(), body)
		if body.Code != http.StatusOK {
			return body
		}

		resp.Body = io.NopCloser(buf)
		return nil
	}

	var as = assert.New(t)

	t.Run("StatusOK", func(t *testing.T) {
		body, err := Post(baseurl + "/200").
			SetErrorChecker(checker).
			Send(nil).
			GetBody()
		as.NoError(err)
		t.Logf("body: %s", string(body))
	})

	t.Run("StatusBadRequest", func(t *testing.T) {
		_, err := Post(baseurl + "/400").
			SetErrorChecker(checker).
			Send(nil).
			GetBody()
		as.Error(err)
		t.Logf("err: %v", err)
	})
}
