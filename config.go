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
		Timeout             time.Duration     // HttpClient超时
		Proxy               string            // 代理, 支持socks5, http, https等协议
		InsecureSkipVerify  bool              // 是否跳过安全检查
		MaxIdleConnsPerHost int               // 每个主机地址的最大空闲连接数
		BeforeFunc          BeforeFunc        // 请求前中间件
		AfterFunc           AfterFunc         // 请求后中间件
		Transport           http.RoundTripper // HTTP传输层
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

// WithTimeout 设置HttpClient超时
func WithTimeout(d time.Duration) Option {
	return func(c *Config) {
		c.Timeout = d
	}
}

// WithMaxIdleConnsPerHost 设置每个主机地址的最大空闲连接数
func WithMaxIdleConnsPerHost(v int) Option {
	return func(c *Config) {
		c.MaxIdleConnsPerHost = v
	}
}

// WithProxy 设置代理, 支持socks5, http, https等协议
func WithProxy(p string) Option {
	return func(c *Config) {
		c.Proxy = p
	}
}

// WithInsecureSkipVerify 设置是否跳过安全检查
func WithInsecureSkipVerify(skip bool) Option {
	return func(c *Config) {
		c.InsecureSkipVerify = true
	}
}

// WithTransport 设置HTTP传输层
// 部分选项会被WithProxy, WithInsecureSkipVerify覆盖, 使用时需注意
func WithTransport(t http.RoundTripper) Option {
	return func(c *Config) {
		c.Transport = t
	}
}

func withInitialize() Option {
	return func(c *Config) {
		if c.Timeout <= 0 {
			c.Timeout = DefaultTimeout
		}

		if c.MaxIdleConnsPerHost <= 0 {
			c.MaxIdleConnsPerHost = DefaultMaxIdleConnsPerHost
		}

		if c.BeforeFunc == nil {
			c.BeforeFunc = defaultBeforeFunc
		}

		if c.AfterFunc == nil {
			c.AfterFunc = defaultAfterFunc
		}

		if c.Transport == nil {
			c.Transport = &http.Transport{
				MaxIdleConnsPerHost: c.MaxIdleConnsPerHost,
			}
		}
	}
}
