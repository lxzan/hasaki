# Hasaki
http request library for golang

[![OSCS Status](https://www.oscs1024.com/platform/badge/lxzan/hasaki.svg?size=small)](https://www.oscs1024.com/project/lxzan/hasaki?ref=badge_small)
- Package
```bash
github.com/lxzan/hasaki
```

### 基本使用
- 发送GET Query请求
```go
package main

import (
	"github.com/lxzan/hasaki"
	"time"
)

type Result = struct {
	TotalCount        int  `json:"total_count"`
	IncompleteResults bool `json:"incomplete_results"`
	Items             []struct {
		Name             string    `json:"name"`
		DisplayName      string    `json:"display_name"`
		ShortDescription string    `json:"short_description"`
		Description      string    `json:"description"`
		CreatedBy        string    `json:"created_by"`
		Released         string    `json:"released"`
		CreatedAt        time.Time `json:"created_at"`
		UpdatedAt        time.Time `json:"updated_at"`
		Featured         bool      `json:"featured"`
		Curated          bool      `json:"curated"`
		Score            float64   `json:"score"`
	} `json:"items"`
}

type Request struct {
	Q       string `form:"q"`
	PerPage int    `form:"per_page"`
}

func main() {
	{
		// 使用结构体
		var query = &Request{Q: "golang", PerPage: 1}
		var result = &Result{}
		err := hasaki.
			Get("https://api.github.com/search/topics").
			SetQuery(query).
			Send(nil).
			BindJSON(result)
	}

	{
		// 使用hasaki.Any
		var result = &Result{}
		err := hasaki.
			Get("https://api.github.com/search/topics").
			SetQuery(hasaki.Any{
				"q":        "golang",
				"per_page": 1,
			}).
			Send(nil).
			BindJSON(result)
	}
}
```

### 发送POST JSON请求
```go
func main() {
	var body = struct {
		Name string `json:"name"`
	}{Name: "caster"}

	// 默认就是JSON编码, SetEncoder(hasaki.JsonEncoder)可省略
	resp, err := hasaki.
		Post("http://localhost:8080" + "/sendJson").
		SetEncoder(hasaki.JsonEncoder).
		Send(body).
		GetBody()
}
```

### 发送POST Form请求
```go
func main() {
	var body = struct {
		Name string `form:"name"`
	}{Name: "caster"}

    // Content-Type: application/x-www-form-urlencoded
	resp, err := hasaki.
		Post("http://localhost:8080" + "/sendForm").
		SetEncoder(hasaki.FormEncoder).
		Send(body).
		GetBody()
}
```

### 发送字节流
```go
func main() {
	var file io.Reader
	resp, err := hasaki.
		Put("http://localhost:8080" + "/upload").
		SetEncoder(hasaki.StreamEncoder).
		Send(file).
		GetBody()
}
```

- Basic Usage
```go
// POST https://api.github.com/
hasaki.
    Post("https://api.github.com/").
    Send(nil).
    BindJSON()

// POST http://127.0.0.1:8080/server.php
type Request struct {
    Name    string `form:"name"`
    Options []int  `form:"options"`
}

hasaki.
    Post("http://127.0.0.1:8080/server.php").
    SetEncoder(hasaki.FormEncoder).
    Send(Request{
        Name:    "aha",
        Options: []int{1, 2, 3, 4},
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
