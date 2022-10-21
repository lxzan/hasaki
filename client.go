package hasaki

import (
	"context"
	"crypto/tls"
	"github.com/pkg/errors"
	"net/http"
	neturl "net/url"
	"time"
)

type Client struct {
	err   error
	check ErrorChecker
	cli   *http.Client
}

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

func (c *Client) Get(url string) *Request {
	return &Request{
		err:     c.err,
		ctx:     context.Background(),
		client:  defaultHTTPClient,
		check:   defaultErrorChecker,
		method:  http.MethodGet,
		url:     url,
		encoder: JsonEncoder,
	}
}

func (c *Client) Post(url string) *Request {
	return &Request{
		err:     c.err,
		ctx:     context.Background(),
		client:  defaultHTTPClient,
		check:   defaultErrorChecker,
		method:  http.MethodPost,
		url:     url,
		encoder: JsonEncoder,
	}
}

func (c *Client) Put(url string) *Request {
	return &Request{
		err:     c.err,
		ctx:     context.Background(),
		client:  defaultHTTPClient,
		check:   defaultErrorChecker,
		method:  http.MethodPut,
		url:     url,
		encoder: JsonEncoder,
	}
}

func (c *Client) Delete(url string) *Request {
	return &Request{
		err:     c.err,
		ctx:     context.Background(),
		client:  defaultHTTPClient,
		check:   defaultErrorChecker,
		method:  http.MethodDelete,
		url:     url,
		encoder: JsonEncoder,
	}
}

func (c *Client) Request(method string, url string) *Request {
	return &Request{
		err:     c.err,
		ctx:     context.Background(),
		client:  defaultHTTPClient,
		check:   defaultErrorChecker,
		method:  method,
		url:     url,
		encoder: JsonEncoder,
	}
}
