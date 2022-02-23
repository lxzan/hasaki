package hasaki

import (
	"github.com/lxzan/hasaki/options"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	var client = NewClient(options.New().SetTimeOut(time.Millisecond))
	body, err := client.Post("https://api.github.com/").Send(nil).GetBody()
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Logf("Response: %s\n", string(body))
	}
}
