package main

import (
	"github.com/lxzan/hasaki"
	"time"
)

func main() {
	var options = new(hasaki.Options).
		SetTimeOut(10 * time.Second).
		SetProxyURL("http://127.0.0.1:10809")
	var client = hasaki.NewClient(options)

	content, err := client.Get("https://google.com/").Send(nil).GetBody()
	println(&content, err)
}
