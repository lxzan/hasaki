package hasaki

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var DefaultClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
	},
}

type Request struct {
	client  *http.Client
	checker ErrorChecker
	method  string
	url     string
	headers Form
	encoder Encoder
}

func NewRequest(method string, url string) *Request {
	var request = &Request{
		client:  DefaultClient,
		checker: DefaultErrorChecker,
		method:  strings.ToUpper(method),
		url:     url,
		encoder: JsonEncoder,
		headers: Form{
			"Content-Type": ContentType_JSON.String(),
		},
	}
	return request
}

func Get(url string) *Request {
	return NewRequest(Method_GET, url)
}

func Post(url string) *Request {
	return NewRequest(Method_POST, url)
}

func Put(url string) *Request {
	return NewRequest(Method_PUT, url)
}

func Delete(url string) *Request {
	return NewRequest(Method_DELETE, url)
}

func (this *Request) SetClient(client *http.Client) *Request {
	this.client = client
	return this
}

func (this *Request) SetEncoder(encoder Encoder) {
	this.encoder = encoder
	this.headers["Content-Type"] = encoder.GetContentType()
}

// SetErrorChecker check response error
func (this *Request) SetErrorChecker(checker ErrorChecker) *Request {
	this.checker = checker
	return this
}

func (this *Request) SetHeaders(headers Form) *Request {
	for k, v := range headers {
		this.headers[k] = v
	}
	return this
}

// Send only support json and form
func (this *Request) Send(param interface{}) *Response {
	if param == nil {
		param = Any{}
	}

	var reader io.Reader
	if this.method == Method_GET {
		URL, err := neturl.Parse(this.url)
		if err != nil {
			return &Response{err: errors.WithStack(err)}
		}
		encodeBytes, err := FormEncoder.Encode(param)
		if err != nil {
			return &Response{err: errors.WithStack(err)}
		}
		this.url = fmt.Sprintf("%s://%s%s?%s", URL.Scheme, URL.Host, URL.Path, string(encodeBytes))
		return this.Raw(reader)
	}

	encodeBytes, err := this.encoder.Encode(param)
	if err != nil {
		return &Response{err: errors.WithStack(err)}
	}
	reader = bytes.NewReader(encodeBytes)
	return this.Raw(reader)
}

func (this *Request) Raw(r io.Reader) (response *Response) {
	response = &Response{}

	req, err1 := http.NewRequest(this.method, this.url, r)
	if err1 != nil {
		response.err = errors.WithStack(err1)
		return
	}

	for k, v := range this.headers {
		req.Header.Set(k, v)
	}

	resp, err2 := this.client.Do(req)
	if err2 != nil {
		response.err = errors.WithStack(err2)
		return
	}

	response.Response = resp
	response.err = this.checker(resp)
	return
}
