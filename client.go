package hasaki

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Client HTTP客户端
// HTTP client
type Client struct {
	config *config
}

// NewClient 新建一个客户端
// Create a new client
func NewClient(options ...Option) (*Client, error) {
	var conf = new(config)
	for _, f := range options {
		f(conf)
	}
	withInitialize()(conf)
	var client = &Client{config: conf}
	return client, nil
}

// Get 发送GET请求
// Send GET request
func (c *Client) Get(url string, args ...any) *Request {
	return c.Request(http.MethodGet, url, args...)
}

// Post 发送POST请求
// Send POST request
func (c *Client) Post(url string, args ...any) *Request {
	return c.Request(http.MethodPost, url, args...)
}

// Put 发送PUT请求
// Send PUT request
func (c *Client) Put(url string, args ...any) *Request {
	return c.Request(http.MethodPut, url, args...)
}

// Delete 发送DELETE请求
// Send DELETE request
func (c *Client) Delete(url string, args ...any) *Request {
	return c.Request(http.MethodDelete, url, args...)
}

// Head 发送HEAD请求
// Send HEAD request
func (c *Client) Head(url string, args ...any) *Request {
	return c.Request(http.MethodHead, url, args...)
}

// Options 发送OPTIONS请求
// Send OPTIONS request
func (c *Client) Options(url string, args ...any) *Request {
	return c.Request(http.MethodOptions, url, args...)
}

// Patch 发送PATCH请求
// Send PATCH request
func (c *Client) Patch(url string, args ...any) *Request {
	return c.Request(http.MethodPatch, url, args...)
}

// Request 创建HTTP请求
// Create HTTP request
func (c *Client) Request(method string, url string, args ...any) *Request {
	if len(args) > 0 {
		url = fmt.Sprintf(url, args...)
	}

	// 如果配置了 BaseURL，则合并 URL
	if c.config.BaseURL != "" {
		url = c.config.BaseURL + url
	}

	r := &Request{
		ctx:              context.Background(),
		client:           c.config.HTTPClient,
		method:           strings.ToUpper(method),
		url:              url,
		before:           c.config.BeforeFunc,
		after:            c.config.AfterFunc,
		headers:          http.Header{},
		reuseBodyEnabled: c.config.ReuseBodyEnabled,
	}

	r.SetEncoder(JsonCodec)

	return r
}
