package hasaki

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
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
	b, err := io.ReadAll(c.Body)
	_ = c.Body.Close()
	c.Body = io.NopCloser(bytes.NewReader(b)) // reset body, so that it can be read again
	return b, errors.WithStack(err)
}

func (c *Response) BindJSON(v any) error { return c.Bind(v, JsonDecode) }

func (c *Response) BindXML(v any) error { return c.Bind(v, XmlDecode) }

func (c *Response) BindForm(v *url.Values) error { return c.Bind(v, FormDecode) }

func (c *Response) Bind(v any, decode func(r io.Reader, ptr any) error) error {
	_, err := c.ReadBody()
	if err != nil {
		return errors.WithStack(err)
	}
	err = decode(c.Body, v)
	return errors.WithStack(err)
}
