package hasaki

import (
	"bytes"
	"crypto/tls"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"time"
)

type Request struct {
	method      string
	url         string
	neturl      *neturl.URL
	contentType ContentTypeInterface
	header      Form
	client      *http.Client
	option      *RequestOption
	err         error
}

type RequestOption struct {
	TimeOut            time.Duration
	RetryTimes         int
	ProxyURL           string
	InsecureSkipVerify bool // skip verify certificate
}

var defaultClient = &http.Client{
	Timeout: 10 * time.Second,
}

func NewRequest(method string, url string, opt ...*RequestOption) *Request {
	var option *RequestOption
	if len(opt) == 0 {
		option = &RequestOption{}
	} else {
		option = opt[0]
	}

	if option.TimeOut == 0 {
		option.TimeOut = 10 * time.Second
	}
	if option.RetryTimes == 0 {
		option.RetryTimes = 1
	}

	var client *http.Client
	if opt == nil {
		client = defaultClient
	} else {
		var transport = &http.Transport{}
		if option.ProxyURL != "" {
			URL := neturl.URL{}
			urlProxy, _ := URL.Parse(option.ProxyURL)
			transport.Proxy = http.ProxyURL(urlProxy)
		}
		if option.InsecureSkipVerify {
			transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}

		client = &http.Client{
			Transport: transport,
			Timeout:   option.TimeOut,
		}
	}

	var req = &Request{
		method:      strings.ToUpper(method),
		url:         url,
		contentType: JSON,
		option:      option,
		client:      client,
		header:      Form{},
	}
	URL, err := neturl.Parse(url)
	if err != nil {
		req.err = err
	}
	req.neturl = URL
	return req
}

func Get(url string, opt ...*RequestOption) *Request {
	return NewRequest("get", url, opt...)
}

func Post(url string, opt ...*RequestOption) *Request {
	return NewRequest("post", url, opt...)
}

func Put(url string, opt ...*RequestOption) *Request {
	return NewRequest("put", url, opt...)
}

func Delete(url string, opt ...*RequestOption) *Request {
	return NewRequest("delete", url, opt...)
}

func (c *Request) Type(contentType ContentTypeInterface) *Request {
	c.contentType = contentType
	return c
}

func (c *Request) Set(header Form) *Request {
	c.header = header
	return c
}

func (c *Request) Send(param Any) (response *Response) {
	response = &Response{}
	if c.err != nil {
		return &Response{err: c.err}
	}

	if param == nil {
		param = Any{}
	}

	var r io.Reader
	if c.method == "GET" {
		var query = c.neturl.Query()
		var qs = ""
		if len(query) > 0 || len(param) > 0 {
			for k, item := range query {
				if len(item) > 1 {
					param[k] = item
				} else {
					param[k] = item[0]
				}
			}
			qs = "?" + FormEncode(param)
		}
		c.url = fmt.Sprintf("%s://%s%s%s", c.neturl.Scheme, c.neturl.Host, c.neturl.Path, qs)
	} else {
		if c.contentType == FORM {
			r = strings.NewReader(FormEncode(param))
		} else if c.contentType == JSON {
			b, _ := jsoniter.Marshal(param)
			r = bytes.NewReader(b)
		}
	}
	return c.Raw(r)
}

func (c *Request) Raw(r io.Reader) (response *Response) {
	response = &Response{}
	if c.err != nil {
		return &Response{err: c.err}
	}

	for i := 1; i <= c.option.RetryTimes; i++ {
		req, err1 := http.NewRequest(c.method, c.url, r)
		if err1 != nil {
			if i == c.option.RetryTimes {
				return &Response{err: err1}
			}
			continue
		}

		req.Header.Set("Content-Type", c.contentType.String())
		for k, v := range c.header {
			req.Header.Set(k, v)
		}

		res, err2 := c.client.Do(req)
		if err2 != nil {
			if i == c.option.RetryTimes {
				return &Response{err: err2}
			}
			continue
		}

		response.Response = res
	}

	return
}
