package hasaki

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	defaultTimeout             = 30 * time.Second
	defaultMaxIdleConnsPerHost = 128
	defaultMaxConnsPerHost     = 128
)

type (
	BeforeFunc func(ctx context.Context, request *http.Request) (context.Context, error)
	AfterFunc  func(ctx context.Context, response *http.Response) (context.Context, error)
)

var (
	once_initer = sync.Once{}

	defaultClient, _ = NewClient(WithHTTPClient(&http.Client{
		Timeout: defaultTimeout,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: defaultMaxIdleConnsPerHost,
			MaxConnsPerHost:     defaultMaxConnsPerHost,
		},
	}))

	defaultBeforeFunc BeforeFunc = func(ctx context.Context, request *http.Request) (context.Context, error) {
		return ctx, nil
	}

	defaultAfterFunc AfterFunc = func(ctx context.Context, response *http.Response) (context.Context, error) {
		return ctx, nil
	}
)

// SetClient 设置全局客户端
// Setting the global client
func SetClient(c *Client) {
	once_initer.Do(func() { defaultClient = c })
}

type (
	config struct {
		BeforeFunc       BeforeFunc   // 请求前中间件
		AfterFunc        AfterFunc    // 请求后中间件
		HTTPClient       *http.Client // HTTP客户端
		ReuseBodyEnabled bool         // 是否复用body
	}

	Option func(c *config)
)

// WithBefore 设置请求前中间件
// Setting up pre-request middleware
func WithBefore(fn BeforeFunc) Option {
	return func(c *config) {
		c.BeforeFunc = fn
	}
}

// WithAfter 设置请求后中间件
// Setting up post-request middleware
func WithAfter(fn AfterFunc) Option {
	return func(c *config) {
		c.AfterFunc = fn
	}
}

// WithHTTPClient 设置HTTP客户端
// Setting the HTTP client
func WithHTTPClient(client *http.Client) Option {
	return func(c *config) {
		c.HTTPClient = client
	}
}

// WithReuseBody 开启body复用; Response.Body可以被断言为BytesReader, 重复调用Bytes()方法.
// Enable body reuse; Response.Body can be asserted as a BytesReader, calling the Bytes() method repeatedly.
func WithReuseBody() Option {
	return func(c *config) {
		c.ReuseBodyEnabled = true
	}
}

func withInitialize() Option {
	return func(c *config) {

		if c.BeforeFunc == nil {
			c.BeforeFunc = defaultBeforeFunc
		}

		if c.AfterFunc == nil {
			c.AfterFunc = defaultAfterFunc
		}

		if c.HTTPClient == nil {
			c.HTTPClient = &http.Client{
				Timeout: defaultTimeout,
				Transport: &http.Transport{
					MaxIdleConnsPerHost: defaultMaxIdleConnsPerHost,
					MaxConnsPerHost:     defaultMaxConnsPerHost,
				},
			}
		}
	}
}

type BytesReader interface {
	io.Reader
	Bytes() []byte
}
