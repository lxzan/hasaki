package hasaki

import "net/http"

type Response struct {
	HttpResponse *http.Response
	Body         []byte
}
