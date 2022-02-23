package hasaki

import (
	"errors"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type Response struct {
	*http.Response
	err error
}

func (c *Response) Err() error {
	return c.err
}

func (c *Response) GetBody() (content []byte, err error) {
	if c.err != nil {
		return nil, c.err
	}
	if c.Response == nil {
		return nil, errors.New("response is nil")
	}
	content, err = ioutil.ReadAll(c.Body)
	c.Body.Close()
	return
}

func (c *Response) BindJSON(v interface{}) error {
	if c.err != nil {
		return c.err
	}
	if c.Response == nil {
		return errors.New("response is nil")
	}

	content, err := ioutil.ReadAll(c.Body)
	if err != nil {
		return err
	}
	c.Body.Close()
	return jsoniter.Unmarshal(content, v)
}
