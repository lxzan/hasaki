package hasaki

import (
	"testing"
)

func TestGet(t *testing.T) {
	res, _ := Post("https://api.github.com/", nil).Send(nil)
	println(string(res.Body))
}
