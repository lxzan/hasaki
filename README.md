# Hasaki
http request library for golang

[![Build Status](https://github.com/lxzan/hasaki/workflows/Go%20Test/badge.svg?branch=master)](https://github.com/lxzan/hasaki/actions?query=branch%3Amaster)
[![OSCS Status](https://www.oscs1024.com/platform/badge/lxzan/hasaki.svg?size=small)](https://www.oscs1024.com/project/lxzan/hasaki?ref=badge_small)

- Package
```bash
github.com/lxzan/hasaki
```

### 基本使用
#### 发送GET请求
```go
package main

import (
	"github.com/lxzan/hasaki"
)

func main() {
	result := make([]map[string]interface{}, 0)
	err := hasaki.
		Get("https://api.github.com/users/%s/repos", "lxzan").
		SetQuery(hasaki.Any{
			"type":     "go",
			"per_page": 1,
		}).
		Send(nil).
		BindJSON(&result)
}
```

#### 发送POST JSON请求
```go
func main() {
	var body = struct {
		Name string `json:"name"`
	}{Name: "caster"}

	// 默认就是JSON编码, SetEncoder(hasaki.JsonEncoder)可省略
	err := hasaki.
		Post("http://localhost:8080" + "/sendJson").
		SetEncoder(hasaki.JsonEncoder).
		Send(body).
		Err()
}
```

#### 发送POST Form请求
```go
func main() {
	var body = struct {
		Name string `form:"name"`
	}{Name: "caster"}

    // Content-Type: application/x-www-form-urlencoded
	err := hasaki.
		Post("http://localhost:8080" + "/sendForm").
		SetEncoder(hasaki.FormEncoder).
		Send(body).
		Err()
}
```

#### 发送字节流
```go
package main

import (
	"github.com/lxzan/hasaki"
	"os"
)

func main() {
	file, _ := os.Open("")
	err := hasaki.
		Put("http://localhost:8080" + "/upload").
		SetHeader(hasaki.H{"Content-Type": hasaki.ContentTypeSTREAM}).
		Send(file).
		Err()
}
```

### 高级

#### 设置代理
```go
cli := hasaki.NewClient().SetProxyURL("socks5://127.0.0.1:1080")
```

#### 统一的错误处理
```go
package main

import (
	"bytes"
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/lxzan/hasaki"
	"github.com/pkg/errors"
	"io"
	"mime"
	"net/http"
)

type BaseResult struct {
	Code    *int   `json:"code"`
	Message string `json:"message"`
}

func After(ctx context.Context, response *http.Response) (context.Context, error) {
	if response.StatusCode != http.StatusOK {
		return ctx, hasaki.ErrUnexpectedStatusCode
	}

	typ, _, _ := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if typ != "application/json" {
		return ctx, nil
	}

	var buf = bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, response.Body); err != nil {
		return ctx, err
	}
	defer response.Body.Close()
	var result = BaseResult{}
	if err := jsoniter.Unmarshal(buf.Bytes(), &result); err != nil {
		return ctx, err
	}
	if result.Code == nil || *result.Code != 0 {
		return ctx, errors.New(result.Message)
	}
	response.Body = io.NopCloser(buf)
	return ctx, nil
}

func main() {
	cli := hasaki.NewClient().SetAfter(After)
}
```

#### 统计请求耗时
```go
package main

import (
	"context"
	"fmt"
	"github.com/lxzan/hasaki"
	"net/http"
	"time"
)

func main() {
	cli := hasaki.
		NewClient().
		SetBefore(func(ctx context.Context, request *http.Request) (context.Context, error) {
			ctx = context.WithValue(ctx, "t0", time.Now())
			return ctx, nil
		}).
		SetAfter(func(ctx context.Context, response *http.Response) (context.Context, error) {
			t0 := ctx.Value("t0").(time.Time)
			fmt.Printf("cost = %dms\n", time.Since(t0).Milliseconds())
			return ctx, nil
		})
	cli.Get("https://api.github.com/users/%s", "lxzan").Send(nil)
}

```