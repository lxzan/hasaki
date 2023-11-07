package hasaki

import (
	"context"
	"io"
	"net/http"

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
	return b, errors.WithStack(err)
}

func (c *Response) BindJSON(v any) error { return c.Bind(v, JsonDecode) }

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
