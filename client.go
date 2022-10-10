package hasaki

import (
	"crypto/tls"
	"net/http"
	neturl "net/url"
	"time"
)

type Client struct {
	cli *http.Client
}

func NewClient() *Client {
	return &Client{cli: &http.Client{
		Timeout: DefaultTimeout,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
		},
	}}
}

func (c *Client) getTransport() *http.Transport {
	return c.cli.Transport.(*http.Transport)
}

func (c *Client) SetTimeOut(d time.Duration) *Client {
	c.cli.Timeout = d
	return c
}

// SetTransport
// if you used SetTransport after SetTimeOut or SetProxyURL,
// SetTimeOut/SetProxyURL will be invalid
func (c *Client) SetTransport(transport *http.Transport) *Client {
	c.cli.Transport = transport
	return c
}

func (c *Client) SetProxyURL(url string) *Client {
	urlProxy, err := neturl.Parse(url)
	if err != nil {
		panic(err)
	}
	c.getTransport().Proxy = http.ProxyURL(urlProxy)
	return c
}

func (c *Client) SetInsecureSkipVerify(skip bool) *Client {
	c.getTransport().TLSClientConfig = &tls.Config{InsecureSkipVerify: skip}
	return c
}

func (c *Client) Get(url string) *Request {
	return NewRequest(Method_GET, url).SetClient(c.cli)
}

func (c *Client) Post(url string) *Request {
	return NewRequest(Method_POST, url).SetClient(c.cli)
}

func (c *Client) Put(url string) *Request {
	return NewRequest(Method_PUT, url).SetClient(c.cli)
}

func (c *Client) Delete(url string) *Request {
	return NewRequest(Method_DELETE, url).SetClient(c.cli)
}
