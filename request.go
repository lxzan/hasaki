package hasaki

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	neturl "net/url"
)

type Request struct {
	err     error
	ctx     context.Context
	client  *http.Client
	check   ErrorChecker
	method  string
	url     string
	headers H
	encoder Encoder
}

func NewRequest(method string, url string) *Request {
	var request = &Request{
		ctx:     context.Background(),
		client:  defaultHTTPClient,
		check:   defaultErrorChecker,
		method:  method,
		url:     url,
		encoder: JsonEncoder,
		headers: H{
			"Content-Type": ContentType_JSON.String(),
		},
	}
	return request
}

func Get(url string) *Request {
	return NewRequest(http.MethodGet, url)
}

func Post(url string) *Request {
	return NewRequest(http.MethodPost, url)
}

func Put(url string) *Request {
	return NewRequest(http.MethodPut, url)
}

func Delete(url string) *Request {
	return NewRequest(http.MethodDelete, url)
}

func (c *Request) setClient(client *http.Client) *Request {
	c.client = client
	return c
}

func (c *Request) setError(err error) *Request {
	c.err = err
	return c
}

// SetErrorChecker check response error
func (c *Request) setErrorChecker(checker ErrorChecker) *Request {
	c.check = checker
	return c
}

func (c *Request) SetEncoder(encoder Encoder) *Request {
	c.encoder = encoder
	c.headers["Content-Type"] = encoder.GetContentType()
	return c
}

func (c *Request) SetHeader(headers H) *Request {
	for k, v := range headers {
		c.headers[k] = v
	}
	return c
}

func (c *Request) SetContext(ctx context.Context) *Request {
	c.ctx = ctx
	return c
}

func (c *Request) SetQuery(query interface{}) *Request {
	URL, err := neturl.Parse(c.url)
	if err != nil {
		c.err = errors.WithStack(err)
		return c
	}
	encodeBytes, err := FormEncoder.Encode(query)
	if err != nil {
		c.err = errors.WithStack(err)
		return c
	}
	c.url = fmt.Sprintf("%s://%s%s?%s", URL.Scheme, URL.Host, URL.Path, string(encodeBytes))
	return c
}

func (c *Request) Send(v interface{}) *Response {
	reader, ok := v.(io.Reader)
	if !ok {
		encodeBytes, err := c.encoder.Encode(v)
		if err != nil {
			return &Response{err: errors.WithStack(err)}
		}
		reader = bytes.NewReader(encodeBytes)
	}

	response := &Response{}
	if c.err != nil {
		response.err = c.err
		return response
	}

	req, err1 := http.NewRequestWithContext(c.ctx, c.method, c.url, reader)
	if err1 != nil {
		response.err = errors.WithStack(err1)
		return response
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	resp, err2 := c.client.Do(req)
	if err2 != nil {
		response.err = errors.WithStack(err2)
		return response
	}

	response.Response = resp
	response.err = c.check(resp)
	return response
}
