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

func (this *Request) Type(contentType string) *Request {
	this.contentType = contentType
	return this
}

func (this *Request) Set(header Form) *Request {
	this.header = header
	return this
}

func (this *Request) Send(param JSON) (*Response, error) {
	if this.err != nil {
		return nil, this.err
	}

	if param == nil {
		param = JSON{}
	}

	var r io.Reader
	if this.method == "GET" {
		var query = this.url.Query()
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
		this.link = fmt.Sprintf("%s://%s%s%s", this.url.Scheme, this.url.Host, this.url.Path, qs)
	} else {
		if this.contentType == FormType {
			r = strings.NewReader(FormEncode(param))
		} else if this.contentType == JsonType {
			b, _ := json.Marshal(param)
			r = bytes.NewReader(b)
		}
	}
	return this.Raw(r)
}

func (this *Request) Raw(r io.Reader) (*Response, error) {
	if this.err != nil {
		return nil, this.err
	}

	var req, err = http.NewRequest(this.method, this.link, r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", this.contentType)
	for k, v := range this.header {
		req.Header.Set(k, v)
	}

	var res, requestError = this.client.Do(req)
	if requestError != nil {
		return nil, requestError
	}

	body, readError := ioutil.ReadAll(res.Body)
	if readError != nil {
		return nil, readError
	}
	return &Response{
		HttpResponse: res,
		Body:         body,
	}, nil
}
