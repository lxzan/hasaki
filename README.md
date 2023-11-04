# Hasaki

http request library for golang

[![Build Status][1]][2] [![codecov][3]][4]

[1]: https://github.com/lxzan/hasaki/workflows/Go%20Test/badge.svg?branch=master

[2]: https://github.com/lxzan/hasaki/actions?query=branch%3Amaster

[3]: https://codecov.io/gh/lxzan/hasaki/graph/badge.svg?token=0VY55RLS3G

[4]: https://codecov.io/gh/lxzan/hasaki

- [Hasaki](#hasaki)
    - [Features](#features)
    - [Install](#install)
    - [Usage](#usage)
      - [GET](#get)
      - [POST](#post)
    - [Middleware](#middleware)

### Features

- [x] Buffer Pool
- [x] Trace the Error Stack
- [x] JSON / WWWForm Encoder
- [x] Request Before and After Middleware

### Install

```bash
go get -v github.com/lxzan/hasaki
```

### Usage

##### GET

```go
package main

import (
	"log"

	"github.com/lxzan/hasaki"
)

func main() {
	type Query struct {
		Q     string `form:"q"`
		Page  int    `form:"page"`
		Order string `form:"-"`
	}
	var out = make(map[string]any)
	var err = hasaki.
		Get("https://api.github.com/search/repositories").
		SetQuery(Query{
			Q:    "go-ws",
			Page: 1,
		}).
		Send(nil).
		BindJSON(out)
	if err != nil {
		log.Printf("%+v", err)
	}
}

```

##### POST

```go
package main

import (
	"log"

	"github.com/lxzan/hasaki"
)

func main() {
	type Query struct {
		Q     string `form:"q"`
		Page  int    `form:"page"`
		Order string `form:"-"`
	}
	var out = make(map[string]any)
	var err = hasaki.
		Post("https://api.github.com/search/repositories").
		SetEncoder(hasaki.FormEncoder).
		Send(Query{
			Q:    "go-ws",
			Page: 1,
		}).
		BindJSON(&out)
	if err != nil {
		log.Printf("%+v", err)
	}
}
```

### Middleware

```go
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/lxzan/hasaki"
)

func main() {
	before := hasaki.WithBefore(func(ctx context.Context, request *http.Request) (context.Context, error) {
		return context.WithValue(ctx, "t0", time.Now()), nil
	})

	after := hasaki.WithAfter(func(ctx context.Context, response *http.Response) (context.Context, error) {
		t0 := ctx.Value("t0").(time.Time)
		log.Printf("latency=%s", time.Since(t0).String())
		return ctx, nil
	})

	var url = "https://api.github.com/search/repositories"
	cli, _ := hasaki.NewClient(before, after)
	cli.Get(url).Send(nil)
}
```