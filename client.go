package hasaki

import (
	"bytes"
	"encoding/json"
	"github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ContentType string

const (
	FORM_TYPE ContentType = "application/x-www-form-urlencoded"
	JSON_TYPE             = "application/json; charset=utf-8"
)

type Client struct {
	requestParam *RequestParam
	requestBody  io.Reader
	link         string
	contentType  ContentType
	method       string
	headers      Form
	response     *http.Response
	httpClient   *http.Client
}

type RequestParam struct {
	ContentType
	TimeOut    time.Duration
	RetryTimes int
	ProxyURL   string
}

func Request(method string, URL string, opt *RequestParam) *Client {
	if opt == nil {
		opt = &RequestParam{
			TimeOut:     2 * time.Second,
			RetryTimes:  3,
			ContentType: FORM_TYPE,
		}
	} else {
		if opt.TimeOut == 0 {
			opt.TimeOut = 2 * time.Second
		}
		if opt.RetryTimes == 0 {
			opt.RetryTimes = 3
		}
		if opt.ContentType == "" {
			opt.ContentType = FORM_TYPE
		}
	}

	var client = &Client{
		requestBody:  strings.NewReader(""),
		headers:      Form{},
		requestParam: opt,
		method:       strings.ToUpper(method),
		link:         URL,
		httpClient:   http.DefaultClient,
	}
	client.httpClient.Timeout = opt.TimeOut

	if opt.ProxyURL != "" {
		urli := url.URL{}
		urlProxy, _ := urli.Parse(opt.ProxyURL)
		client.httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlProxy),
			},
			Timeout: opt.TimeOut,
		}
	}
	return client
}

func Get(url string, opt *RequestParam) *Client {
	return Request("GET", url, opt)
}

func Put(url string, opt *RequestParam) *Client {
	return Request("PUT", url, opt)
}

func Post(url string, opt *RequestParam) *Client {
	return Request("POST", url, opt)
}

func Delete(url string, opt *RequestParam) *Client {
	return Request("DELETE", url, opt)
}

func (this *Client) Set(headers Form) *Client {
	for k, v := range headers {
		this.headers[k] = v
	}
	return this
}

func (this *Client) SetBody(r io.Reader) *Client {
	this.requestBody = r
	return this
}

func (this *Client) Send(data JSON) *Client {
	if this.requestParam.ContentType == JSON_TYPE {
		body, _ := jsoniter.Marshal(data)
		this.requestBody = bytes.NewReader(body)
		return this
	}

	var form = url.Values{}
	for k, v := range data {
		var item = make([]string, 0)
		var varType = reflect.TypeOf(v).String()

		switch varType {
		case "[]string":
			item = v.([]string)
			form[k+"[]"] = item
			break
		case "[]int":
			arr := v.([]int)
			for _, num := range arr {
				item = append(item, strconv.Itoa(num))
			}
			form[k+"[]"] = item
			break
		case "int":
			tmp, _ := v.(int64)
			form.Set(k, strconv.Itoa(int(tmp)))
			break
		default:
			tmp, _ := v.(string)
			form.Set(k, tmp)
			break
		}
	}

	qs := form.Encode()
	if this.method == "GET" || this.method == "DELETE" {
		u, _ := url.Parse(this.link)
		if u.RawQuery != "" {
			qs += "&" + u.RawQuery
		}
		this.link = u.Scheme + "://" + u.Host + u.Path + "?" + qs
	} else {
		this.requestBody = strings.NewReader(qs)
	}
	return this
}

func (this *Client) GetResponse() (resp *http.Response, err error) {
	for i := 0; i < this.requestParam.RetryTimes; i++ {
		resp, err = this.query()
		if err == nil {
			break
		}
	}
	return resp, err
}

func (this *Client) query() (*http.Response, error) {
	var req, err = http.NewRequest(this.method, this.link, this.requestBody)
	if err != nil {
		return nil, err
	}

	for k, v := range this.headers {
		req.Header.Set(k, v)
	}
	if this.requestParam.ContentType != "" {
		req.Header.Set("Content-Type", string(this.requestParam.ContentType))
	}

	res, err := this.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	this.response = res
	return res, err
}

func (this *Client) GetBody() ([]byte, error) {
	res, err := this.query()
	if err != nil {
		return []byte{}, err
	}
	return ioutil.ReadAll(res.Body)
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
	res, err := this.query()
	if err != nil {
		return 0, err
	}
	return res.StatusCode, nil
}
