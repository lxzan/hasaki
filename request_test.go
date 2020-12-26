package hasaki

import (
	"testing"
)

func TestGet(t *testing.T) {
	body, err := Post("https://api.github.com/rate_limit", &RequestOption{
		InsecureSkipVerify: true,
	}).Send(nil).GetBody()
	if err != nil {
		println(string(body))
	}
}
