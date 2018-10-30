package hasaki

import (
	"testing"
	"time"
)

func TestClient_Json(t *testing.T) {
	opt := &RequestOption{
		TimeOut:    5 * time.Second,
		RetryTimes: 1,
		ProxyURL:   "",
	}
	resp, err := Get("https://api.github.com/", opt).GetBody()

	println(string(resp))
}
