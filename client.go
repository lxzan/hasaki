package hasaki

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	config *config
}

// NewClient 新建一个客户端
// Create a new client
func NewClient(options ...Option) (*Client, error) {
	var c = new(config)
	options = append(options, withInitialize())
	for _, f := range options {
		f(c)
	}
	var client = &Client{config: c}
	return client, nil
}

func (c *Client) Get(url string, args ...any) *Request {
	return c.Request(http.MethodGet, url, args...)
}

func (c *Client) Post(url string, args ...any) *Request {
	return c.Request(http.MethodPost, url, args...)
}

func (c *Client) Put(url string, args ...any) *Request {
	return c.Request(http.MethodPut, url, args...)
}

func (c *Client) Delete(url string, args ...any) *Request {
	return c.Request(http.MethodDelete, url, args...)
}

func (c *Client) Head(url string, args ...any) *Request {
	return c.Request(http.MethodHead, url, args...)
}

func (c *Client) Options(url string, args ...any) *Request {
	return c.Request(http.MethodOptions, url, args...)
}

func (c *Client) Patch(url string, args ...any) *Request {
	return c.Request(http.MethodPatch, url, args...)
}

func (c *Client) Request(method string, url string, args ...any) *Request {
	if len(args) > 0 {
		url = fmt.Sprintf(url, args...)
	}

	r := &Request{
		ctx:     context.Background(),
		client:  c.config.HTTPClient,
		method:  strings.ToUpper(method),
		url:     url,
		before:  c.config.BeforeFunc,
		after:   c.config.AfterFunc,
		headers: http.Header{},
	}

	r.SetEncoder(c.config.Encoder)

	return r
}
