package hasaki

import (
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

type Response struct {
	*http.Response
	err error
}

func (c *Response) Err() error {
	return c.err
}

func (c *Response) GetBody() ([]byte, error) {
	if c.err != nil {
		return nil, c.err
	}
	if c.Response == nil || c.Body == nil {
		return nil, errors.New("response is nil")
	}
	defer c.Body.Close()
	return ioutil.ReadAll(c.Body)
}

func (c *Response) BindJSON(v interface{}) error {
	if c.err != nil {
		return c.err
	}
	if c.Response == nil || c.Body == nil {
		return errors.New("response is nil")
	}
	defer c.Body.Close()
	return errors.WithStack(jsoniter.NewDecoder(c.Body).Decode(v))
}
