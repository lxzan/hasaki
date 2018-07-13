package hasaki

import (
	"encoding/json"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	URL          string
	ContentType  string
	Method       string
	Timeout      time.Duration
	Headers      Form
	Form         JSON
	responseBody []byte
}

func Request(method string, url string) *Client {
	method = strings.ToUpper(method)
	var client = new(Client)
	client.Headers = Form{}
	client.Form = JSON{}
	client.Method = method
	client.URL = url
	client.Timeout = 2 * time.Second
	return client
}

func Get(url string) *Client {
	return Request("GET", url)
}

func Put(url string) *Client {
	return Request("PUT", url)
}

func Post(url string) *Client {
	return Request("POST", url)
}

func Delete(url string) *Client {
	return Request("DELETE", url)
}

func (this *Client) Type(contentType string) *Client {
	switch contentType {
	case "json":
		this.ContentType = "json"
		this.Headers["Content-Type"] = "application/json; charset=utf-8"
		break
	default:
		this.ContentType = "json"
		this.Headers["Content-Type"] = "application/x-www-form-urlencoded"
		break
	}
	return this
}

func (this *Client) SetTimeout(t time.Duration) *Client {
	this.Timeout = t
	return this
}

func (this *Client) Set(headers Form) *Client {
	for k, v := range headers {
		this.Headers[k] = v
	}
	return this
}

func (this *Client) Send(data JSON) *Client {
	this.Form = data
	return this
}

func (this *Client) GetResponse() (*http.Response, error) {
	var form = url.Values{}
	var err error
	for k, v := range this.Form {
		var item = make([]string, 0)
		var varType = reflect.TypeOf(v).String()
		if varType == "[]string" {
			item = v.([]string)
		} else if varType == "int" {
			item = append(item, strconv.Itoa(int(reflect.ValueOf(v).Int())))
		} else {
			item = append(item, reflect.ValueOf(v).String())
		}
		form[k] = item
	}

	var req = &http.Request{}
	queryString := form.Encode()
	if this.Method == "GET" || this.Method == "DELETE" {
		re, _ := regexp.Compile(`\?.*?=.*?`)
		if re.MatchString(this.URL) {
			this.URL = this.URL + "&" + queryString
		} else {
			this.URL = this.URL + "?" + queryString
		}
		req, err = http.NewRequest(this.Method, this.URL, nil)
	} else {
		if this.ContentType == "json" {
			bytes, _ := jsoniter.Marshal(this.Form)
			req, err = http.NewRequest(this.Method, this.URL, strings.NewReader(string(bytes)))
		} else {
			req, err = http.NewRequest(this.Method, this.URL, strings.NewReader(queryString))
		}
		if this.ContentType == "" {
			this.Headers["Content-Type"] = "application/x-www-form-urlencoded"
		}
	}

	for k, v := range this.Headers {
		req.Header.Set(k, v)
	}

	if err != nil {
		return nil, err
	}

	httpClient := http.DefaultClient
	httpClient.Timeout = this.Timeout
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(res.Body)
	this.responseBody = body
	defer res.Body.Close()
	return res, err
}

func (this *Client) GetBody() ([]byte, error) {
	_, err := this.GetResponse()
	return this.responseBody, err
}

func (this *Client) Json() (JSON, error) {
	bytes, err := this.GetBody()
	var res = JSON{}
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}
