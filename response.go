package hasaki

import (
	"net/http"
)

type Response struct {
	*http.Response
	responseBody []byte
}

func (c *Response) GetBody() []byte {
	return c.responseBody
}
