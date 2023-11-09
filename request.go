package hasaki

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strings"

	"github.com/go-playground/form/v4"
	"github.com/pkg/errors"
)

var (
	errEmptyResponse   = errors.New("unexpected empty response")
	errUnsupportedData = errors.New("unsupported data type")
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
	debug   bool
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

// Debug 开启调试模式, 打印CURL命令
// Enable debug mode, print CURL commands
func (c *Request) Debug() *Request {
	c.debug = true
	return c
}

// SetBefore 设置请求前中间件
// Setting up pre-request middleware
func (c *Request) SetBefore(f BeforeFunc) *Request {
	c.before = f
	return c
}

// SetAfter 设置请求后中间件
// Setting up post-request middleware
func (c *Request) SetAfter(f AfterFunc) *Request {
	c.after = f
	return c
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

	switch v := query.(type) {
	case string:
		if len(v) > 0 {
			URL.RawQuery = v
		}
	case neturl.Values:
		URL.RawQuery = v.Encode()
	default:
		values, err := form.NewEncoder().Encode(query)
		if err != nil {
			c.err = errors.WithStack(err)
			return c
		}
		URL.RawQuery = values.Encode()
	}
	c.url = URL.String()
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

	reader, err := c.encoder.Encode(body)
	if err != nil {
		response.err = err
		return response
	}

	req, err1 := http.NewRequestWithContext(c.ctx, c.method, c.url, reader)
	if err1 != nil {
		response.err = errors.WithStack(err1)
		return response
	}

	if c.method == http.MethodGet && body == nil {
		c.headers.Del("Content-Type")
	}
	req.Header = c.headers

	// 执行请求前中间件
	response.ctx, response.err = c.before(c.ctx, req)
	if response.err != nil {
		return response
	}

	// 打印CURL命令
	c.printCURL(req)

	// 发起请求
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

func (c *Request) printCURL(req *http.Request) {
	if !c.debug {
		return
	}

	var body = bytes.NewBufferString("")
	if req.Body != nil {
		_, _ = io.Copy(body, req.Body)
	}

	var builder = strings.Builder{}
	var line = fmt.Sprintf("curl -X %s '%s' \\\n", req.Method, req.URL.String())
	builder.WriteString(line)
	for k, arr := range req.Header {
		for _, v := range arr {
			line = fmt.Sprintf("    --header '%s: %s'", k, v)
			builder.WriteString(line)
			if req.Body != nil {
				builder.WriteString(" \\\n")
			} else {
				builder.WriteString(" \n")
			}
		}
	}

	builder.WriteString("    --data-raw ")
	builder.WriteString("'")
	if body.Len() < 128*1024 {
		s := strings.TrimSuffix(body.String(), "\n")
		s = strings.Replace(s, "'", "\\'", -1)
		builder.WriteString(s)
	}
	builder.WriteString("'")
	println(builder.String())
	req.Body = io.NopCloser(body)
}
