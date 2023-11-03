package hasaki

import (
	"context"
	"fmt"
	"net/http"
)

type Client struct {
	config *Config
}

// NewClient 新建一个客户端, 支持自定义HttpClient, 错误检查和中间件
func NewClient(options ...Option) (*Client, error) {
	var c = new(Config)
	options = append(options, withInitialize())
	for _, f := range options {
		f(c)
	}
	var client = &Client{config: c}
	return client, nil
}

func (c *Client) Get(url string, args ...interface{}) *Request {
	return c.Request(http.MethodGet, url, args...)
}

func (c *Client) Post(url string, args ...interface{}) *Request {
	return c.Request(http.MethodPost, url, args...)
}

func (c *Client) Put(url string, args ...interface{}) *Request {
	return c.Request(http.MethodPut, url, args...)
}

func (c *Client) Delete(url string, args ...interface{}) *Request {
	return c.Request(http.MethodDelete, url, args...)
}

func (c *Client) Request(method string, url string, args ...interface{}) *Request {
	if len(args) > 0 {
		url = fmt.Sprintf(url, args...)
	}
	return &Request{
		ctx:     context.Background(),
		client:  c.config.HTTPClient,
		method:  method,
		url:     url,
		encoder: JsonEncoder,
		before:  c.config.BeforeFunc,
		after:   c.config.AfterFunc,
		headers: http.Header{},
	}
}
