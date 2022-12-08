package hasaki

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	neturl "net/url"
)

type Client struct {
	config *Config
	cli    *http.Client
}

// NewClient 新建一个客户端, 支持自定义HttpClient, 错误检查和中间件
func NewClient(options ...Option) (*Client, error) {
	var config = &Config{}
	withInitialize()(config)
	for _, fn := range options {
		fn(config)
	}

	var client = &Client{
		cli: &http.Client{
			Transport: config.Transport,
			Timeout:   config.Timeout,
		},
		config: config,
	}

	if transport, ok := client.cli.Transport.(*http.Transport); ok {
		if config.Proxy != "" {
			urlProxy, err := neturl.Parse(config.Proxy)
			if err != nil {
				return nil, err
			}
			transport.Proxy = http.ProxyURL(urlProxy)
		}

		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		transport.MaxIdleConnsPerHost = config.MaxIdleConnsPerHost
	}

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
		client:  c.cli,
		method:  method,
		url:     url,
		encoder: JsonEncoder,
		before:  c.config.BeforeFunc,
		after:   c.config.AfterFunc,
	}
}
