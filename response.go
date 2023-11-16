package hasaki

import (
	"context"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
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

func (c *Response) ReadBody() ([]byte, error) {
	if c.err != nil {
		return nil, c.err
	}
	if c.Response == nil || c.Body == nil {
		return nil, errors.WithStack(errEmptyResponse)
	}
	if v, ok := c.Body.(BytesReader); ok {
		return v.Bytes(), nil
	}
	b, err := io.ReadAll(c.Body)
	_ = c.Body.Close()
	return b, errors.WithStack(err)
}

func (c *Response) BindJSON(v any) error { return c.Bind(v, JsonDecode) }

func (c *Response) BindXML(v any) error { return c.Bind(v, XmlDecode) }

func (c *Response) BindForm(v *url.Values) error { return c.Bind(v, FormDecode) }

func (c *Response) Bind(v any, decode func(r io.Reader, ptr any) error) error {
	if c.err != nil {
		return c.err
	}
	if c.Response == nil || c.Body == nil {
		return errors.WithStack(errEmptyResponse)
	}
	err := decode(c.Body, v)
	_ = c.Body.Close()
	return errors.WithStack(err)
}
