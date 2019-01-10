package hasaki

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Param struct {
	URL    string
	Method string
	Header Form
	Option *Option
	Client *http.Client
}

type Option struct {
	TimeOut    time.Duration
	RetryTimes int
	ProxyURL   string
}

type Context struct {
	Param    Param
	Request  *http.Request
	Response *http.Response
	Body     []byte
	Error    error
}

var defaultClient = &http.Client{
	Timeout: 10 * time.Second,
}

func NewRequest(method string, URL string, opt *Option) *Context {
	if opt == nil {
		opt = &Option{
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

	return &Context{
		Param: Param{
			URL:    URL,
			Method: strings.ToUpper(method),
			Header: Form{},
			Option: opt,
			Client: client,
		},
	}
}

func Get(url string, opt *Option) *Context {
	return NewRequest("GET", url, opt)
}

func Put(url string, opt *Option) *Context {
	return NewRequest("PUT", url, opt)
}

func Post(url string, opt *Option) *Context {
	return NewRequest("POST", url, opt)
}

func Delete(url string, opt *Option) *Context {
	return NewRequest("DELETE", url, opt)
}

// set header
func (this *Context) Set(header Form) *Context {
	this.Param.Header = header
	return this
}

// only for get request
func (this *Context) Query(param JSON) *Context {
	if len(param) > 0 {
		this.Param.URL += "?" + querystring(param)
	}
	this.Request, this.Error = http.NewRequest(this.Param.Method, this.Param.URL, nil)
	if this.Error != nil {
		return this
	}

	for k, v := range this.Param.Header {
		this.Request.Header.Set(k, v)
	}
	this.Response, this.Error = this.Param.Client.Do(this.Request)

	if this.Error == nil {
		this.Body, this.Error = ioutil.ReadAll(this.Response.Body)
	}
	return this
}

// send body
func (this *Context) Send(contentType ContentType, body io.Reader) *Context {
	this.Request, this.Error = http.NewRequest(this.Param.Method, this.Param.URL, body)
	if this.Error != nil {
		return this
	}

	this.Param.Header["Content-Type"] = string(contentType)
	for k, v := range this.Param.Header {
		this.Request.Header.Set(k, v)
	}
	this.Response, this.Error = this.Param.Client.Do(this.Request)

	if this.Error == nil {
		this.Body, this.Error = ioutil.ReadAll(this.Response.Body)
	}
	return this
}