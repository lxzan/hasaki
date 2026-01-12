package hasaki

import (
	"context"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
)

// Response HTTP响应结构体
// HTTP response struct
type Response struct {
	*http.Response
	ctx context.Context
	err error
}

// Err 返回响应错误
// Return response error
func (c *Response) Err() error {
	return c.err
}

// Context 返回响应上下文
// Return response context
func (c *Response) Context() context.Context {
	return c.ctx
}

// ReadBody 读取响应体
// Read response body
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

// BindJSON 将响应体绑定为JSON格式
// Bind response body as JSON
func (c *Response) BindJSON(v any) error { return c.Bind(v, JsonCodec) }

// BindXML 将响应体绑定为XML格式
// Bind response body as XML
func (c *Response) BindXML(v any) error { return c.Bind(v, XmlCodec) }

// BindForm 将响应体绑定为表单格式
// Bind response body as form
func (c *Response) BindForm(v *url.Values) error { return c.Bind(v, FormCodec) }

// Bind 使用指定的解码器将响应体绑定到目标对象
// Bind response body to target object using specified decoder
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
