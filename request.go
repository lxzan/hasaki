package hasaki

import (
	"context"
	"fmt"
	"github.com/go-playground/form/v4"
	"io"
	"net/http"
	neturl "net/url"

	"github.com/pkg/errors"
)

var (
	ErrDataNotSupported     = errors.New("data type is not supported")
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)

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

func NewRequest(method string, url string, args ...interface{}) *Request {
	if len(args) > 0 {
		url = fmt.Sprintf(url, args...)
	}
	var request = &Request{
		ctx:     context.Background(),
		client:  defaultHTTPClient,
		method:  method,
		url:     url,
		before:  defaultBeforeFunc,
		after:   defaultAfterFunc,
		headers: http.Header{},
	}
	request.SetEncoder(JsonEncoder)
	return request
}

func Get(url string, args ...interface{}) *Request {
	return NewRequest(http.MethodGet, url, args...)
}

func Post(url string, args ...interface{}) *Request {
	return NewRequest(http.MethodPost, url, args...)
}

func Put(url string, args ...interface{}) *Request {
	return NewRequest(http.MethodPut, url, args...)
}

func Delete(url string, args ...interface{}) *Request {
	return NewRequest(http.MethodDelete, url, args...)
}

// SetEncoder 设置编码器
func (c *Request) SetEncoder(encoder Encoder) *Request {
	c.encoder = encoder
	c.headers.Set("Content-Type", encoder.GetContentType())
	return c
}

// SetHeader 设置请求头
func (c *Request) SetHeader(k, v string) *Request {
	c.headers.Add(k, v)
	return c
}

// SetContext 设置请求上下文, 可用于单个请求级别的超时控制
func (c *Request) SetContext(ctx context.Context) *Request {
	c.ctx = ctx
	return c
}

// SetQuery 设置URL Query String
// 支持hasaki.Any | struct(使用form标签) 等数据类型
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
// 支持 hasaki.Any, struct, io.Reader 等数据类型
func (c *Request) Send(body interface{}) *Response {
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
