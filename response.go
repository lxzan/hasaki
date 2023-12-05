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
	if v, ok := c.Body.(BytesReadCloser); ok {
		return v.Bytes(), nil
	}
	b, err := io.ReadAll(c.Body)
	_ = c.Body.Close()
	return b, errors.WithStack(err)
}

func (c *Response) BindJSON(v any) error { return c.Bind(v, JsonCodec) }

func (c *Response) BindXML(v any) error { return c.Bind(v, XmlCodec) }

func (c *Response) BindForm(v *url.Values) error { return c.Bind(v, FormCodec) }

func (c *Response) Bind(v any, decoder Decoder) error {
	if c.err != nil {
		return c.err
	}
	if c.Response == nil || c.Body == nil {
		return errors.WithStack(errEmptyResponse)
	}
	err := decoder.Decode(c.Body, v)
	_ = c.Body.Close()
	return errors.WithStack(err)
}
