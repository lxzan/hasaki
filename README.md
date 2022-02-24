# Hasaki
http request library for golang

- Package
```bash
github.com/lxzan/hasaki
```

- Basic Usage
```go
// GET https://api.github.com/
hasaki.
	Get("https://api.github.com/").
	Send(nil).
	GetBody()

// GET http://127.0.0.1:8080/server.php?hello%5B%5D=world&hello%5B%5D=%E8%BF%9E%E7%BB%AD%E6%80%A7&me=lxzan
hasaki.
	Get("http://127.0.0.1:8080/server.php").
	Send(hasaki.Any{
		"hello": []string{"world", "连续性"},
		"me":    "lxzan",
	}).
	GetBody()

// POST
hasaki.
	POST("http://127.0.0.1:9999/").
        SetHeaders(hasaki.Form{
		"X-Access-Token": token,
		"X-Running-Env":  env,
	}).
	Send(nil).
	GetBody()
```

- Advanced

```go
var cli = hasaki.NewClient().SetProxyURL("socks5://127.0.0.1:10808")
body, err := cli.Get("https://google.com/").Send(nil).GetBody()
```
