package hasaki

import (
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

func (c *Response) GetBody() ([]byte, error) {
	defer func() {
		if c.Body != nil {
			c.Body.Close()
		}
	}()
	if c.err != nil {
		return nil, c.err
	}
	return ioutil.ReadAll(c.Body)
}

func (c *Response) BindJSON(v interface{}) error {
	defer func() {
		if c.Body != nil {
			c.Body.Close()
		}
	}()

	if c.err != nil {
		return c.err
	}
	body, err := ioutil.ReadAll(c.Body)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(body, v)
}
