package hasaki

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	neturl "net/url"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	err    error
	check  ErrorChecker
	cli    *http.Client
	before func(ctx context.Context, request *http.Request) (context.Context, error)
	after  func(ctx context.Context, request *http.Response) (context.Context, error)
}

// NewClient 新建一个客户端, 支持自定义HttpClient, 错误检查和中间件
func NewClient() *Client {
	return &Client{
		check: defaultErrorChecker,
		cli:   defaultHTTPClient,
	}
}

func (c *Client) SetTimeOut(d time.Duration) *Client {
	c.cli.Timeout = d
	return c
}

func (c *Client) SetProxyURL(url string) *Client {
	urlProxy, err := neturl.Parse(url)
	if err != nil {
		c.err = errors.WithStack(err)
		return c
	}

	c.getTransport().Proxy = http.ProxyURL(urlProxy)
	return c
}

func (c *Client) SetInsecureSkipVerify(skip bool) *Client {
	c.getTransport().TLSClientConfig = &tls.Config{InsecureSkipVerify: skip}
	return c
}

func (c *Client) getTransport() *http.Transport {
	return c.cli.Transport.(*http.Transport)
}

// SetTransport
// if you used SetTransport after SetTimeOut or SetProxyURL,
// SetTimeOut/SetProxyURL will be invalid
func (c *Client) SetTransport(transport *http.Transport) *Client {
	c.cli.Transport = transport
	return c
}

// SetErrorChecker check response error
func (c *Client) SetErrorChecker(checker ErrorChecker) *Client {
	c.check = checker
	return c
}

func (c *Client) SetBefore(fn func(ctx context.Context, request *http.Request) (context.Context, error)) *Client {
	c.before = fn
	return c
}

func (c *Client) SetAfter(fn func(ctx context.Context, request *http.Response) (context.Context, error)) *Client {
	c.after = fn
	return c
}

func (c *Client) Get(url string, args ...interface{}) *Request {
	return &Request{
		err:     c.err,
		ctx:     context.Background(),
		client:  defaultHTTPClient,
		check:   defaultErrorChecker,
		method:  http.MethodGet,
		url:     fmt.Sprintf(url, args...),
		encoder: JsonEncoder,
		before:  c.before,
		after:   c.after,
	}
}

func (c *Client) Post(url string, args ...interface{}) *Request {
	return &Request{
		err:     c.err,
		ctx:     context.Background(),
		client:  defaultHTTPClient,
		check:   defaultErrorChecker,
		method:  http.MethodPost,
		url:     fmt.Sprintf(url, args...),
		encoder: JsonEncoder,
		before:  c.before,
		after:   c.after,
	}
}

func (c *Client) Put(url string, args ...interface{}) *Request {
	return &Request{
		err:     c.err,
		ctx:     context.Background(),
		client:  defaultHTTPClient,
		check:   defaultErrorChecker,
		method:  http.MethodPut,
		url:     fmt.Sprintf(url, args...),
		encoder: JsonEncoder,
		before:  c.before,
		after:   c.after,
	}
}

func (c *Client) Delete(url string, args ...interface{}) *Request {
	return &Request{
		err:     c.err,
		ctx:     context.Background(),
		client:  defaultHTTPClient,
		check:   defaultErrorChecker,
		method:  http.MethodDelete,
		url:     fmt.Sprintf(url, args...),
		encoder: JsonEncoder,
		before:  c.before,
		after:   c.after,
	}
}

func (c *Client) Request(method string, url string, args ...interface{}) *Request {
	return &Request{
		err:     c.err,
		ctx:     context.Background(),
		client:  defaultHTTPClient,
		check:   defaultErrorChecker,
		method:  method,
		url:     fmt.Sprintf(url, args...),
		encoder: JsonEncoder,
		before:  c.before,
		after:   c.after,
	}
}
