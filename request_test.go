package hasaki

import (
	"testing"
)

func TestGet(t *testing.T) {
	body, err := Post("https://api.github.com/", &RequestOption{
		InsecureSkipVerify: false,
	}).Send(nil).GetBody()
	if err != nil {
		t.Error(err.Error())
	} else {
		println(string(body))
	}
}
