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

func (this *Client) getTransport() *http.Transport {
	return this.cli.Transport.(*http.Transport)
}

func (this *Client) SetTimeOut(d time.Duration) *Client {
	this.cli.Timeout = d
	return this
}

func (this *Client) SetTransport(transport *http.Transport) *Client {
	this.cli.Transport = transport
	return this
}

func (this *Client) SetProxyURL(url string) *Client {
	urlProxy, err := neturl.Parse(url)
	if err != nil {
		panic(err)
	}
	this.getTransport().Proxy = http.ProxyURL(urlProxy)
	return this
}

func (this *Client) SetInsecureSkipVerify(skip bool) *Client {
	this.getTransport().TLSClientConfig = &tls.Config{InsecureSkipVerify: skip}
	return this
}

func (this *Client) Get(url string) *Request {
	return NewRequest(Method_GET, url).SetClient(this.cli)
}

func (this *Client) Post(url string) *Request {
	return NewRequest(Method_POST, url).SetClient(this.cli)
}

func (this *Client) Put(url string) *Request {
	return NewRequest(Method_PUT, url).SetClient(this.cli)
}

func (this *Client) Delete(url string) *Request {
	return NewRequest(Method_DELETE, url).SetClient(this.cli)
}
