package hasaki

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"time"
)

var defaultClient = &http.Client{
	Timeout:   10 * time.Second,
	Transport: &http.Transport{},
}

type Request struct {
	//err      error
	client   *http.Client
	encoding Encoding
	method   string
	url      string
	headers  Form
}

func NewRequest(method string, url string) *Request {
	var request = &Request{
		encoding: Encoding_JSON,
		method:   strings.ToUpper(method),
		url:      url,
		headers:  Form{},
	}
	return request
}

func (this *Request) setClient(client *http.Client) *Request {
	this.client = client
	return this
}

func (this *Request) SetEncoding(encoding Encoding) *Request {
	this.encoding = encoding
	switch encoding {
	case Encoding_FORM:
		this.headers["Content-Type"] = ContentType_FORM.String()
	case Encoding_JSON:
		this.headers["Content-Type"] = ContentType_JSON.String()
	}
	return this
}

func (this *Request) SetHeaders(headers Form) *Request {
	for k, v := range headers {
		this.headers[k] = v
	}
	return this
}

// Send support json/form
func (this *Request) Send(param Any) *Response {
	if param == nil {
		param = Any{}
	}
	var r io.Reader

	if this.method == Method_GET {
		URL, err := neturl.Parse(this.url)
		if err != nil {
			return &Response{err: errors.WithStack(err)}
		}

		var query = URL.Query()
		var qs = ""
		if len(query) > 0 || len(param) > 0 {
			for k, item := range query {
				if len(item) > 1 {
					param[k] = item
				} else {
					param[k] = item[0]
				}
			}
			encodingText, err := FormEncode(param)
			if err != nil {
				return &Response{err: errors.WithStack(err)}
			}
			qs = "?" + encodingText
		}
		this.url = fmt.Sprintf("%s://%s%s%s", URL.Scheme, URL.Host, URL.Path, qs)
	} else {
		if this.encoding == Encoding_FORM {
			encodingText, err := FormEncode(param)
			if err != nil {
				return &Response{err: errors.WithStack(err)}
			}
			r = strings.NewReader(encodingText)
		} else if this.encoding == Encoding_JSON {
			b, _ := jsoniter.Marshal(param)
			r = bytes.NewReader(b)
		}
	}
	return this.Raw(r)
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

	if this.client == nil {
		this.client = defaultClient
	}
	resp, err2 := this.client.Do(req)
	if err2 != nil {
		response.err = errors.WithStack(err2)
		return
	}

	response.Response = resp
	return
}
