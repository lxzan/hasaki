package hasaki

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

type (
	BeforeFunc func(ctx context.Context, request *http.Request) (context.Context, error)
	AfterFunc  func(ctx context.Context, response *http.Response) (context.Context, error)
)

const (
	DefaultTimeout             = 30 * time.Second
	DefaultMaxIdleConnsPerHost = 32
)

var (
	defaultHTTPClient = &http.Client{
		Timeout: DefaultTimeout,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
		},
	}

	defaultBeforeFunc BeforeFunc = func(ctx context.Context, request *http.Request) (context.Context, error) {
		return ctx, nil
	}

	defaultAfterFunc AfterFunc = func(ctx context.Context, response *http.Response) (context.Context, error) {
		if response.StatusCode != http.StatusOK {
			return ctx, ErrUnexpectedStatusCode
		}
		return ctx, nil
	}
)

// SetDefaultHttpClient 设置默认HTTP客户端
func SetDefaultHttpClient(client *http.Client) {
	defaultHTTPClient = client
}

// SetDefaultBeforeFunc 设置默认请求前中间件
func SetDefaultBeforeFunc(fn func(ctx context.Context, request *http.Request) (context.Context, error)) {
	defaultBeforeFunc = fn
}

// SetDefaultAfterFunc 设置默认请求后中间件
func SetDefaultAfterFunc(fn func(ctx context.Context, response *http.Response) (context.Context, error)) {
	defaultAfterFunc = fn
}

func NewReadCloser(body io.ReadCloser) (*ReadCloser, error) {
	var rc = &ReadCloser{Buffer: bytes.NewBuffer(nil)}
	_, err := io.Copy(rc, body)
	_ = body.Close()
	return rc, err
}

type ReadCloser struct {
	*bytes.Buffer
}

func (r *ReadCloser) Read(p []byte) (n int, err error) {
	return r.Buffer.Read(p)
}

func (r *ReadCloser) Close() error {
	return nil
}

type (
	Config struct {
		//Timeout             time.Duration     // HttpClient超时
		//Proxy               string            // 代理, 支持socks5, http, https等协议
		//InsecureSkipVerify  bool              // 是否跳过安全检查
		//MaxIdleConnsPerHost int               // 每个主机地址的最大空闲连接数
		//Transport           http.RoundTripper // HTTP传输层

		BeforeFunc BeforeFunc // 请求前中间件
		AfterFunc  AfterFunc  // 请求后中间件
		HTTPClient *http.Client
	}

	Option func(c *Config)
)

// WithBefore 设置请求前中间件
func WithBefore(fn BeforeFunc) Option {
	return func(c *Config) {
		c.BeforeFunc = fn
	}
}

// WithAfter 设置请求后中间件
func WithAfter(fn AfterFunc) Option {
	return func(c *Config) {
		c.AfterFunc = fn
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) {
		c.HTTPClient = client
	}
}

func withInitialize() Option {
	return func(c *Config) {

		if c.BeforeFunc == nil {
			c.BeforeFunc = defaultBeforeFunc
		}

		if c.AfterFunc == nil {
			c.AfterFunc = defaultAfterFunc
		}

		if c.HTTPClient == nil {
			c.HTTPClient = &http.Client{
				Timeout: 30 * time.Second,
			}
		}
	}
}
