# Hasaki
http request library for golang

- Package
```bash
github.com/lxzan/hasaki
```

- Basic Usage
```go
// POST https://api.github.com/
hasaki.
	Post("https://api.github.com/").
	Send(nil).
	BindJSON()

// POST http://127.0.0.1:8080/server.php
hasaki.
	Post("http://127.0.0.1:8080/server.php").
	SetEncoder(hasaki.FormEncoder).
	Send(hasaki.Any{
		"hello": []string{"world", "连续性"},
		"me":    "lxzan",
	}).
	GetBody()

// POST
hasaki.
	Post("http://127.0.0.1:9999/").
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
