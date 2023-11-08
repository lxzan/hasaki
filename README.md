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
      - [Get](#get)
      - [Post](#post)
      - [Stream](#stream)
      - [Error Stack](#error-stack)
      - [Middleware](#middleware)
      - [Debug](#debug)

### Features

- [x] Buffer Pool
- [x] Trace the Error Stack
- [x] JSON / WWWForm Encoder
- [x] Request Before and After Middleware
- [x] Export CURL Command

### Install

```bash
go get -v github.com/lxzan/hasaki
```

### Usage

#### Get

```go
// GET https://api.example.com/search
// Send get request with path parameters

resp := hasaki.
    Get("https://api.example.com/%s", "search").
    Send(nil)
```

```go
// GET https://api.example.com/search?q=hasaki&page=1
// Send get request, with Query parameter, encoded with url.Values

resp := hasaki.
    Get("https://api.example.com/search").
    SetQuery(url.Values{
      "q":    []string{"hasaki"},
      "page": []string{"1"},
    }).
    Send(nil)
```

```go
// GET https://api.example.com/search?q=hasaki&page=1
// Send get request, with Query parameter, encoded with struct

type Req struct {
    Q    string `form:"q"`
    Page int    `form:"page"`
}
resp := hasaki.
    Get("https://api.example.com/search").
    SetQuery(Req{
        Q:    "hasaki",
        Page: 1,
    }).
    Send(nil)
```

#### Post

```go
// POST https://api.example.com/search
// Send post request, encoded with json

type Req struct {
    Q    string `json:"q"`
    Page int    `json:"page"`
}
resp := hasaki.
    Post("https://api.example.com/search").
    Send(Req{
        Q:    "hasaki",
        Page: 1,
    })
```

```go
// POST https://api.example.com/search
// Send post request, encoded with www-form

type Req struct {
    Q    string `form:"q"`
    Page int    `form:"page"`
}
resp := hasaki.
    Post("https://api.example.com/search").
    SetEncoder(hasaki.FormEncoder).
    Send(Req{
        Q:    "hasaki",
        Page: 1,
    })
```

#### Stream

```go
// POST https://api.example.com/upload
// Send a put request, using a byte stream

var reader io.Reader
encoder := hasaki.NewStreamEncoder(hasaki.MimeSTREAM)
resp := hasaki.
    Put("https://api.example.com/upload").
    SetEncoder(encoder).
    Send(reader)
```

#### Error Stack

```go
// Print the error stack
data := make(map[string]any)
err := hasaki.
    Post("https://api.example.com/upload").
    Send(nil).
    BindJSON(&data)
if err != nil {
    log.Printf("%+v", err)
}
```

#### Middleware

```go
// Statistics on time spent on requests

before := hasaki.WithBefore(func (ctx context.Context, request *http.Request) (context.Context, error) {
    return context.WithValue(ctx, "t0", time.Now()), nil
})

after := hasaki.WithAfter(func (ctx context.Context, response *http.Response) (context.Context, error) {
    t0 := ctx.Value("t0").(time.Time)
    log.Printf("latency=%s", time.Since(t0).String())
    return ctx, nil
})

var url = "https://api.github.com/search/repositories"
cli, _ := hasaki.NewClient(before, after)
cli.Get(url).Send(nil)
```

#### Debug

```go
hasaki.
    Post("http://127.0.0.1:3000").
    Debug().
    SetHeader("x-username", "123").
    SetHeader("user-agent", "hasaki").
    SetEncoder(hasaki.FormEncoder).
    Send(url.Values{
        "name": []string{"洪荒"},
        "age":  []string{"456"},
    })
```

```bash
curl -X POST 'http://127.0.0.1:3000' \
    --header 'User-Agent: hasaki' \
    --header 'Content-Type: application/x-www-form-urlencoded' \
    --header 'X-Username: 123' \
    --data-raw 'age=456&name=%E6%B4%AA%E8%8D%92'
```