package hasaki

import (
	"crypto/tls"
	"net/http"
	neturl "net/url"
)

type Client struct {
	cli *http.Client
}

func NewClient(options ...*Options) *Client {
	var client = &http.Client{}
	if len(options) == 0 {
		options = []*Options{new(Options)}
	}
	if options[0].timeOut == 0 {
		options[0].timeOut = DefaultTimeout
	}

	var option = options[0]
	client.Timeout = option.timeOut
	var transport = &http.Transport{}
	if option.proxyURL != "" {
		URL := neturl.URL{}
		urlProxy, err := URL.Parse(option.proxyURL)
		if err != nil {
			panic(err)
		}
		transport.Proxy = http.ProxyURL(urlProxy)
	}
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: option.insecureSkipVerify}
	client.Transport = transport

	return &Client{cli: client}
}

func (this *Client) Get(url string) *Request {
	return NewRequest(Method_GET, url).setClient(this.cli)
}

func (this *Client) Post(url string) *Request {
	return NewRequest(Method_POST, url).setClient(this.cli)
}

func (this *Client) Put(url string) *Request {
	return NewRequest(Method_PUT, url).setClient(this.cli)
}

func (this *Client) Delete(url string) *Request {
	return NewRequest(Method_DELETE, url).setClient(this.cli)
}
