package hasaki

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"
)

const (
	// defaultTimeout 默认请求超时时间
	// Default request timeout
	defaultTimeout = 30 * time.Second
	// defaultMaxIdleConnsPerHost 每个主机的默认最大空闲连接数
	// Default maximum idle connections per host
	defaultMaxIdleConnsPerHost = 128
	// defaultMaxConnsPerHost 每个主机的默认最大连接数
	// Default maximum connections per host
	defaultMaxConnsPerHost = 128
	// defaultMaxIdleConns 全局最大空闲连接数
	// Global maximum idle connections
	defaultMaxIdleConns = 100
	// defaultIdleConnTimeout 空闲连接超时时间
	// Idle connection timeout
	defaultIdleConnTimeout = 90 * time.Second
	// defaultDialTimeout 连接超时时间
	// Dial timeout
	defaultDialTimeout = 30 * time.Second
	// defaultTLSHandshakeTimeout TLS握手超时时间
	// TLS handshake timeout
	defaultTLSHandshakeTimeout = 10 * time.Second
	// defaultResponseHeaderTimeout 响应头超时时间
	// Response header timeout
	defaultResponseHeaderTimeout = 10 * time.Second
	// defaultExpectContinueTimeout Expect: 100-continue 超时时间
	// Expect: 100-continue timeout
	defaultExpectContinueTimeout = 1 * time.Second
)

type (
	// BeforeFunc 请求前中间件函数类型
	// Pre-request middleware function type
	BeforeFunc func(ctx context.Context, request *http.Request) (context.Context, error)
	// AfterFunc 请求后中间件函数类型
	// Post-request middleware function type
	AfterFunc func(ctx context.Context, response *http.Response) (context.Context, error)
)

var (
	// defaultHttpClient 默认HTTP客户端，包含超时和传输配置
	// Default HTTP client with timeout and transport configuration
	defaultHttpClient = &http.Client{
		Timeout: defaultTimeout,
		Transport: &http.Transport{
			// 连接池配置
			// Connection pool configuration
			MaxIdleConns:        defaultMaxIdleConns,
			MaxIdleConnsPerHost: defaultMaxIdleConnsPerHost,
			MaxConnsPerHost:     defaultMaxConnsPerHost,
			IdleConnTimeout:     defaultIdleConnTimeout,

			// 超时配置
			// Timeout configuration
			DialContext: (&net.Dialer{
				Timeout: defaultDialTimeout,
				// KeepAlive 默认开启，显式设置以保持连接活跃
				// KeepAlive is enabled by default, explicitly set to keep connections alive
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   defaultTLSHandshakeTimeout,
			ResponseHeaderTimeout: defaultResponseHeaderTimeout,
			ExpectContinueTimeout: defaultExpectContinueTimeout,

			// 启用 HTTP/2 和 keepalive
			// Enable HTTP/2 and keepalive
			DisableKeepAlives: false,
			ForceAttemptHTTP2: true,
		},
	}

	defaultClient, _ = NewClient(WithHTTPClient(defaultHttpClient))

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
	defaultClient = c
}

type (
	config struct {
		BeforeFunc       BeforeFunc   // 请求前中间件
		AfterFunc        AfterFunc    // 请求后中间件
		HTTPClient       *http.Client // HTTP客户端
		ReuseBodyEnabled bool         // 是否复用body
		BaseURL          string       // 基础URL
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

// WithReuseBody 开启Body可重复读; Response.Body可以被断言为BytesReadCloser, 调用Bytes()方法重复读取.
// Turn on Body repeatable read; Response.Body can be asserted as BytesReadCloser, call Bytes() method to repeat reads.
func WithReuseBody() Option {
	return func(c *config) {
		c.ReuseBodyEnabled = true
	}
}

// WithBaseURL 设置基础URL; 所有请求的URL都会与BaseURL拼接
// Setting the base URL; all request URLs will be concatenated with the BaseURL
func WithBaseURL(baseURL string) Option {
	return func(c *config) {
		c.BaseURL = baseURL
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
			c.HTTPClient = defaultHttpClient
		}
	}
}

type BytesReadCloser interface {
	io.ReadCloser
	Bytes() []byte
}
