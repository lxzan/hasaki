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

	defaultBeforeFunc = func(ctx context.Context, request *http.Request) (context.Context, error) {
		return ctx, nil
	}

	defaultAfterFunc = func(ctx context.Context, response *http.Response) (context.Context, error) {
		if response.StatusCode != http.StatusOK {
			return ctx, ErrUnexpectedStatusCode
		}
		return ctx, nil
	}
)

func SetGlobalClient(client *http.Client) {
	defaultHTTPClient = client
}

func SetBefore(fn func(ctx context.Context, request *http.Request) (context.Context, error)) {
	defaultBeforeFunc = fn
}

func SetAfter(fn func(ctx context.Context, response *http.Response) (context.Context, error)) {
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
		Timeout             time.Duration
		Proxy               string
		InsecureSkipVerify  bool
		MaxIdleConnsPerHost int
		BeforeFunc          BeforeFunc
		AfterFunc           AfterFunc
		Transport           http.RoundTripper
	}

	Option func(c *Config)
)

func WithBefore(fn BeforeFunc) Option {
	return func(c *Config) {
		c.BeforeFunc = fn
	}
}

func WithAfter(fn AfterFunc) Option {
	return func(c *Config) {
		c.AfterFunc = fn
	}
}

func WithTimeout(d time.Duration) Option {
	return func(c *Config) {
		c.Timeout = d
	}
}

func WithMaxIdleConnsPerHost(v int) Option {
	return func(c *Config) {
		c.MaxIdleConnsPerHost = v
	}
}

func WithProxy(p string) Option {
	return func(c *Config) {
		c.Proxy = p
	}
}

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
