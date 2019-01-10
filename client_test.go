package hasaki

import (
	"testing"
)

func TestPost(t *testing.T) {
	var ctx = Get("https://api.github.com/", nil).Query(nil)
	println(string(ctx.Body))
}
