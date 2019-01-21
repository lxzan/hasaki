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
	Method      string
	Link        string
	URL         *url.URL
	ContentType string
	Header      Form
	Client      *http.Client
	Option      *RequestOption
	Body        io.Reader
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
			RetryTimes: 3,
		}
	} else {
		if opt.TimeOut == 0 {
			opt.TimeOut = 10 * time.Second
		}
		if opt.RetryTimes == 0 {
			opt.RetryTimes = 3
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
		Method:      strings.ToUpper(method),
		Link:        link,
		ContentType: FormType,
		Client:      client,
		Header:      Form{},
	}
	URL, err := url.Parse(link)
	if err != nil {
		req.err = err
	}
	req.URL = URL
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
	this.ContentType = contentType
	return this
}

func (this *Request) Set(header Form) *Request {
	this.Header = header
	return this
}

func (this *Request) Send(param JSON) (*Response, error) {
	if this.err != nil {
		return nil, this.err
	}

	if param == nil {
		param = JSON{}
	}

	if this.Method == "GET" {
		var query = this.URL.Query()
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
		this.Link = fmt.Sprintf("%s://%s%s%s", this.URL.Scheme, this.URL.Host, this.URL.Path, qs)
	} else {
		if this.ContentType == FormType {
			this.Body = strings.NewReader(FormEncode(param))
		} else if this.ContentType == JsonType {
			b, _ := json.Marshal(param)
			this.Body = bytes.NewReader(b)
		}
	}

	var req, err = http.NewRequest(this.Method, this.Link, this.Body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", this.ContentType)
	for k, v := range this.Header {
		req.Header.Set(k, v)
	}

	var res, resError = this.Client.Do(req)
	if resError != nil {
		return nil, resError
	}

	body, readError := ioutil.ReadAll(res.Body)
	if resError != nil {
		return nil, readError
	}

	return &Response{
		HttpResponse: res,
		Body:         body,
	}, nil
}
