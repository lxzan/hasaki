package hasaki

import (
	"context"
	"io"
	"net/http"
	"sync/atomic"

	"github.com/pkg/errors"
)

type Response struct {
	*http.Response
	ctx     context.Context
	latency int64
	err     error
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

func (c *Response) BindJSON(v any) error  { return c.BindCodec(v, JSONCodec) }
func (c *Response) BindYAML(v any) error  { return c.BindCodec(v, YAMLCodec) }
func (c *Response) BindXML(v any) error   { return c.BindCodec(v, XMLCodec) }
func (c *Response) BindPROTO(v any) error { return c.BindCodec(v, PROTOCodec) }

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

func (c *Response) BindCodec(v any, codec BytesCodec) error {
	if c.err != nil {
		return c.err
	}
	if c.Response == nil || c.Body == nil {
		return errors.WithStack(errEmptyResponse)
	}
	body, err := io.ReadAll(c.Body)
	defer c.Body.Close()
	if err != nil {
		return errors.WithStack(err)
	}
	err = codec.Unmarshal(body, v)
	return errors.WithStack(err)
}

// 访问延迟，数值单位 ms (毫秒)
// Latency returns the latency for the request/response
func (c *Response) Latency() int64 {
	return atomic.LoadInt64(&c.latency)
}
