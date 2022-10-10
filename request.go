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

var defaultHTTPClient = &http.Client{
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
		client:  defaultHTTPClient,
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

func (c *Request) SetClient(client *http.Client) *Request {
	c.client = client
	return c
}

func (c *Request) SetEncoder(encoder Encoder) *Request {
	c.encoder = encoder
	c.headers["Content-Type"] = encoder.GetContentType()
	return c
}

// SetErrorChecker check response error
func (c *Request) SetErrorChecker(checker ErrorChecker) *Request {
	c.checker = checker
	return c
}

func (c *Request) SetHeaders(headers Form) *Request {
	for k, v := range headers {
		c.headers[k] = v
	}
	return c
}

// Send only support json and form
func (c *Request) Send(param interface{}) *Response {
	if param == nil {
		param = Any{}
	}

	var reader io.Reader
	if c.method == Method_GET {
		URL, err := neturl.Parse(c.url)
		if err != nil {
			return &Response{err: errors.WithStack(err)}
		}
		encodeBytes, err := FormEncoder.Encode(param)
		if err != nil {
			return &Response{err: errors.WithStack(err)}
		}
		c.url = fmt.Sprintf("%s://%s%s?%s", URL.Scheme, URL.Host, URL.Path, string(encodeBytes))
		return c.Raw(reader)
	}

	encodeBytes, err := c.encoder.Encode(param)
	if err != nil {
		return &Response{err: errors.WithStack(err)}
	}
	reader = bytes.NewReader(encodeBytes)
	return c.Raw(reader)
}

func (c *Request) Raw(r io.Reader) (response *Response) {
	response = &Response{}

	req, err1 := http.NewRequest(c.method, c.url, r)
	if err1 != nil {
		response.err = errors.WithStack(err1)
		return
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	resp, err2 := c.client.Do(req)
	if err2 != nil {
		response.err = errors.WithStack(err2)
		return
	}

	response.Response = resp
	response.err = c.checker(resp)
	return
}
