package hasaki

import (
	"testing"
)

func TestGet(t *testing.T) {
	res, _ := Get("https://api3.feng.com/v1/notice/chatList", nil).
		Set(Form{
			"X-Access-Token": "c17a34fc-655e-42f1-8c8f-781de494e729",
		}).
		Send(JSON{
			"account":  "FatalErrorX",
			"password": "Wei123456",
		})
	println(string(res.Body))
}
