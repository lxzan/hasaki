package hasaki

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	method      string
	link        string
	url         *url.URL
	contentType string
	header      Form
	client      *http.Client
	option      *RequestOption
	err         error
}

type RequestOption struct {
	TimeOut    time.Duration
	RetryTimes int
	ProxyURL   string
}

var defaultClient = &http.Client{
	Timeout: 10 * time.Second,
}

func NewRequest(method string, link string, opt *RequestOption) *Request {
	if opt == nil {
		opt = &RequestOption{
			TimeOut:    10 * time.Second,
			RetryTimes: 1,
		}
	} else {
		if opt.TimeOut == 0 {
			opt.TimeOut = 10 * time.Second
		}
		if opt.RetryTimes == 0 {
			opt.RetryTimes = 1
		}
	}

	var client *http.Client
	if opt.TimeOut == defaultClient.Timeout && opt.ProxyURL == "" {
		client = defaultClient
	} else if opt.ProxyURL != "" {
		URL := url.URL{}
		urlProxy, _ := URL.Parse(opt.ProxyURL)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlProxy),
			},
			Timeout: opt.TimeOut,
		}
	}

	var req = &Request{
		method:      strings.ToUpper(method),
		link:        link,
		contentType: JsonType,
		client:      client,
		header:      Form{},
	}
	URL, err := url.Parse(link)
	if err != nil {
		req.err = err
	}
	req.url = URL
	return req
}

func Get(link string, opt *RequestOption) *Request {
	return NewRequest("get", link, opt)
}

func Post(link string, opt *RequestOption) *Request {
	return NewRequest("post", link, opt)
}

func Put(link string, opt *RequestOption) *Request {
	return NewRequest("put", link, opt)
}

func Delete(link string, opt *RequestOption) *Request {
	return NewRequest("delete", link, opt)
}

func (c *Request) Type(contentType string) *Request {
	c.contentType = contentType
	return c
}

func (c *Request) Set(header Form) *Request {
	c.header = header
	return c
}

func (c *Request) Send(param Any) (*Response, error) {
	if c.err != nil {
		return nil, c.err
	}

	if param == nil {
		param = Any{}
	}

	var r io.Reader
	if c.method == "GET" {
		var query = c.url.Query()
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
		c.link = fmt.Sprintf("%s://%s%s%s", c.url.Scheme, c.url.Host, c.url.Path, qs)
	} else {
		if c.contentType == FormType {
			r = strings.NewReader(FormEncode(param))
		} else if c.contentType == JsonType {
			b, _ := json.Marshal(param)
			r = bytes.NewReader(b)
		}
	}
	return c.Raw(r)
}

func (c *Request) Raw(r io.Reader) (*Response, error) {
	if c.err != nil {
		return nil, c.err
	}

	var req, err = http.NewRequest(c.method, c.link, r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", c.contentType)
	for k, v := range c.header {
		req.Header.Set(k, v)
	}

	var res, requestError = c.client.Do(req)
	if requestError != nil {
		return nil, requestError
	}

	body, readError := ioutil.ReadAll(res.Body)
	if readError != nil {
		return nil, readError
	}
	return &Response{
		Response:     res,
		responseBody: body,
	}, nil
}
