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
	link         string
	contentType  string
	method       string
	timeout      time.Duration
	headers      Form
	form         JSON
	responseBody []byte
	statusCode   int
	httpClient   *http.Client
}

func Request(method string, url string) *Client {
	method = strings.ToUpper(method)
	var client = new(Client)
	client.headers = Form{}
	client.form = JSON{}
	client.method = method
	client.link = url
	client.timeout = 2 * time.Second
	client.httpClient = http.DefaultClient
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

func (this *Client) SetProxy(proxyURL string) *Client {
	urli := url.URL{}
	urlProxy, _ := urli.Parse(proxyURL)
	this.httpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlProxy),
		},
	}
	return this
}

func (this *Client) Type(contentType string) *Client {
	switch contentType {
	case "json":
		this.contentType = "json"
		this.headers["Content-Type"] = "application/json; charset=utf-8"
		break
	default:
		this.contentType = "json"
		this.headers["Content-Type"] = "application/x-www-form-urlencoded"
		break
	}
	return this
}

func (this *Client) SetTimeout(t time.Duration) *Client {
	this.timeout = t
	return this
}

func (this *Client) Set(headers Form) *Client {
	for k, v := range headers {
		this.headers[k] = v
	}
	return this
}

func (this *Client) Send(data JSON) *Client {
	this.form = data
	return this
}

func (this *Client) GetResponse() (*http.Response, error) {
	var form = url.Values{}
	var err error
	for k, v := range this.form {
		var item = make([]string, 0)
		var varType = reflect.TypeOf(v).String()
		if varType == "[]string" {
			item = v.([]string)
		} else if varType == "[]int" {
			arr := v.([]int)
			for _, num := range arr {
				item = append(item, strconv.Itoa(num))
			}
		} else if varType == "int" {
			item = append(item, strconv.Itoa(int(reflect.ValueOf(v).Int())))
		} else {
			item = append(item, reflect.ValueOf(v).String())
		}

		if varType == "[]string" || varType == "[]int" {
			form[k+"[]"] = item
		} else {
			form[k] = item
		}
	}

	var req = &http.Request{}
	queryString := form.Encode()
	if (this.method == "GET" || this.method == "DELETE") && queryString != "" {
		re, _ := regexp.Compile(`\?.*?=.*?`)
		if re.MatchString(this.link) {
			this.link = this.link + "&" + queryString
		} else {
			this.link = this.link + "?" + queryString
		}
		req, err = http.NewRequest(this.method, this.link, nil)
	} else {
		if this.contentType == "json" {
			bytes, _ := jsoniter.Marshal(this.form)
			req, err = http.NewRequest(this.method, this.link, strings.NewReader(string(bytes)))
		} else {
			req, err = http.NewRequest(this.method, this.link, strings.NewReader(queryString))
		}
		if this.contentType == "" {
			this.headers["Content-Type"] = "application/x-www-form-urlencoded"
		}
	}

	for k, v := range this.headers {
		req.Header.Set(k, v)
	}

	if err != nil {
		return nil, err
	}

	this.httpClient.Timeout = this.timeout
	res, err := this.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(res.Body)
	this.responseBody = body
	this.statusCode = res.StatusCode
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

func (this *Client) GetStatusCode() (int, error) {
	_, err := this.GetResponse()
	return this.statusCode, err
}
