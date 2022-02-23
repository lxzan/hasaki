package hasaki

import (
	"crypto/tls"
	"net/http"
	neturl "net/url"
)

type Client struct {
	opt *Options
	cli *http.Client
}

func NewClient(opt ...*Options) *Client {
	var client = &Client{
		cli: &http.Client{},
	}

	if len(opt) == 0 {
		client.opt = new(Options).SetTimeOut(DefaultTimeout)
	} else {
		client.opt = opt[0]
		if client.opt.TimeOut == 0 {
			client.opt.TimeOut = DefaultTimeout
		}
	}

	client.cli.Timeout = client.opt.TimeOut
	var transport = &http.Transport{}
	if client.opt.ProxyURL != "" {
		URL := neturl.URL{}
		urlProxy, err := URL.Parse(client.opt.ProxyURL)
		if err != nil {
			panic(err)
		}
		transport.Proxy = http.ProxyURL(urlProxy)
	}
	if client.opt.InsecureSkipVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client.cli.Transport = transport

	return client
}

func (this *Client) Get(url string) *Request {
	return NewRequest(this, Method_GET, url)
}

func (this *Client) Post(url string) *Request {
	return NewRequest(this, Method_POST, url)
}

func (this *Client) Put(url string) *Request {
	return NewRequest(this, Method_PUT, url)
}

func (this *Client) Delete(url string) *Request {
	return NewRequest(this, Method_DELETE, url)
}
