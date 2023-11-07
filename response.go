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
	latency atomic.Int64
	ctx     context.Context
	err     error
}

// 返回 Response 错误
func (c *Response) Error() error {
	return c.err
}

// 返回 Response 上下文
func (c *Response) Context() context.Context {
	return c.ctx
}

// 读取 Response Body 中内容
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

// Latency 获取请求耗时
// Get request latency
func (c *Response) Latency() int64 {
	return c.latency.Load()
}

// BindJSON 绑定JSON解码模块
func (c *Response) BindJSON(v any) error { return c.Bind(v, JSONDecode) }

// BindXML 绑定XML解码模块
func (c *Response) BindXML(v any) error { return c.Bind(v, XMLDecode) }

// Bind 绑定解码模块
// Binding decoding module
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
