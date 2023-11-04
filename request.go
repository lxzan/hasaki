package hasaki

import (
	"context"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"

	"github.com/go-playground/form/v4"
	"github.com/pkg/errors"
)

var errEmptyResponse = errors.New("unexpected empty response")

type Request struct {
	err     error
	ctx     context.Context
	client  *http.Client
	method  string
	url     string
	headers http.Header
	encoder Encoder
	before  BeforeFunc
	after   AfterFunc
}

// NewRequest 新建一个请求
// Create a new request
func NewRequest(method string, url string, args ...any) *Request {
	return defaultClient.Request(method, url, args...)
}

func Get(url string, args ...any) *Request {
	return defaultClient.Get(url, args...)
}

func Post(url string, args ...any) *Request {
	return defaultClient.Post(url, args...)
}

func Put(url string, args ...any) *Request {
	return defaultClient.Put(url, args...)
}

func Delete(url string, args ...any) *Request {
	return defaultClient.Delete(url, args...)
}

// SetEncoder 设置编码器
// Set request body encoder
func (c *Request) SetEncoder(encoder Encoder) *Request {
	c.encoder = encoder
	c.headers.Set("Content-Type", encoder.ContentType())
	return c
}

// SetHeader 设置请求头
// Set Request Header
func (c *Request) SetHeader(k, v string) *Request {
	c.headers.Set(k, v)
	return c
}

// Header 获取请求头
// Get request header
func (c *Request) Header() http.Header {
	return c.headers
}

// SetContext 设置请求上下文
// Set Request context
func (c *Request) SetContext(ctx context.Context) *Request {
	c.ctx = ctx
	return c
}

// SetQuery 设置查询参数, 详情请参考 https://github.com/go-playground/form
// To set the query parameters, please refer to https://github.com/go-playground/form for details.
func (c *Request) SetQuery(query any) *Request {
	URL, err := neturl.Parse(c.url)
	if err != nil {
		c.err = errors.WithStack(err)
		return c
	}

	str := fmt.Sprintf("%s://%s%s", URL.Scheme, URL.Host, URL.Path)
	switch v := query.(type) {
	case string:
		if len(v) > 0 {
			str += "?" + v
		}
	default:
		str += "?"
		values, err := form.NewEncoder().Encode(query)
		if err != nil {
			c.err = errors.WithStack(err)
			return c
		}
		str += values.Encode()
	}
	c.url = str
	return c
}

// Send 发送请求
// Send http request
func (c *Request) Send(body any) *Response {
	response := &Response{ctx: c.ctx}
	if c.err != nil {
		response.err = c.err
		return response
	}

	reader, ok := body.(io.Reader)
	if !ok {
		reader, c.err = c.encoder.Encode(body)
		if c.err != nil {
			response.err = c.err
			return response
		}
	}

	req, err1 := http.NewRequestWithContext(c.ctx, c.method, c.url, reader)
	if err1 != nil {
		response.err = errors.WithStack(err1)
		return response
	}

	// 执行请求前中间件
	response.ctx, response.err = c.before(c.ctx, req)
	if response.err != nil {
		return response
	}

	if c.method == http.MethodGet && body == nil {
		c.headers.Del("Content-Type")
	}
	req.Header = c.headers
	resp, err2 := c.client.Do(req)
	if err2 != nil {
		response.err = errors.WithStack(err2)
		return response
	}

	// 执行请求后中间件
	response.ctx, response.err = c.after(response.ctx, resp)
	if response.err != nil {
		return response
	}

	response.Response = resp
	return response
}
