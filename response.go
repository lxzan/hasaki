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
	return ioutil.ReadAll(c.Body)
}

func (c *Response) BindJSON(v interface{}) error {
	body, err := ioutil.ReadAll(c.Body)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(body, v)
}
