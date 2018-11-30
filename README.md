# Hasaki
http request library for golang

- install
```bash
go get github.com/lxzan/hasaki
```

- usage
```go
// GET https://api.github.com/
hasaki.
	Get("https://api.github.com/", nil).
	Json()

// GET http://127.0.0.1:8080/server.php?hello%5B%5D=world&hello%5B%5D=%E8%BF%9E%E7%BB%AD%E6%80%A7&me=lxzan
hasaki.
	Get("http://127.0.0.1:8080/server.php", nil).
	Send(hasaki.JSON{
		"hello": []string{"world", "连续性"},
		"me":    "lxzan",
	}).
	Json()

// POST
hasaki.
	POST("http://127.0.0.1:9999/", nil).
	Set(hasaki.Form{
		"X-Access-Token": token,
		"X-Running-Env":  env,
	}).
	Json()

// Only view whether request success or not
code,_ := Get("https://api.github.com/", nil).GetStatusCode()
println(code == 200)

// Advanced
opt := &RequestOption{
		TimeOut:    5 * time.Second, // default 2s
		RetryTimes: 1,				 // default 3 times
		ProxyURL:   "",
	}
	resp, err := Get("https://api.github.com/", opt).GetBody()
```
