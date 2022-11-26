package hasaki

import (
	"context"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type Response struct {
	*http.Response
	ctx context.Context
	err error
}

func (c *Response) Err() error {
	return c.err
}

func (c *Response) Context() context.Context {
	return c.ctx
}

func (c *Response) GetBody() ([]byte, error) {
	if c.err != nil {
		return nil, c.err
	}
	if c.Response == nil || c.Body == nil {
		return []byte{}, nil
	}
	if rc, ok := c.Body.(*ReadCloser); ok {
		return rc.Bytes(), nil
	}
	b, err := ioutil.ReadAll(c.Body)
	_ = c.Body.Close()
	return b, err
}

func (c *Response) BindJSON(v interface{}) error {
	if c.err != nil {
		return c.err
	}
	if c.Response == nil || c.Body == nil {
		return nil
	}
	if rc, ok := c.Body.(*ReadCloser); ok {
		return jsoniter.Unmarshal(rc.Bytes(), v)
	}
	err := jsoniter.NewDecoder(c.Body).Decode(v)
	_ = c.Body.Close()
	return err
}
