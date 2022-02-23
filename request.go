package hasaki

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
)

type Request struct {
	err         error
	client      *Client
	contentType string
	method      string
	url         string
	header      Form
}

func NewRequest(client *Client, method string, url string) *Request {
	var request = &Request{
		client:      client,
		contentType: ContentType_JSON,
		method:      strings.ToUpper(method),
		url:         url,
		header:      Form{},
	}
	return request
}

func (this *Request) Type(contentType string) *Request {
	this.contentType = contentType
	return this
}

// Send support json/form
func (this *Request) Send(param Any) *Response {
	if this.err != nil {
		return &Response{err: this.err}
	}

	if param == nil {
		param = Any{}
	}

	var r io.Reader
	if this.method == Method_GET {
		URL, err := neturl.Parse(this.url)
		if err != nil {
			return &Response{err: err}
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
			qs = "?" + FormEncode(param)
		}
		this.url = fmt.Sprintf("%s://%s%s%s", URL.Scheme, URL.Host, URL.Path, qs)
	} else {
		if this.contentType == ContentType_FORM {
			r = strings.NewReader(FormEncode(param))
		} else if this.contentType == ContentType_JSON {
			b, _ := jsoniter.Marshal(param)
			r = bytes.NewReader(b)
		}
	}
	return this.Raw(r)
}

func (this *Request) Raw(r io.Reader) (response *Response) {
	response = &Response{}
	if this.err != nil {
		return &Response{err: this.err}
	}

	req, err1 := http.NewRequest(this.method, this.url, r)
	if err1 != nil {
		return &Response{err: err1}
	}

	req.Header.Set("Content-Type", this.contentType)
	for k, v := range this.header {
		req.Header.Set(k, v)
	}

	resp, err2 := this.client.cli.Do(req)
	if err2 != nil {
		return &Response{err: err2}
	}

	response.Response = resp
	return
}
