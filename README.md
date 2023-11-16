# Hasaki

HTTP Request Library for Go

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

### Features

-   [x] Buffer Pool
-   [x] Trace the Error Stack
-   [x] Build-In JSON / XML / WWWForm / Protobuf/YAML Codec 
-   [x] Request Before and After Middleware
-   [x] Export cURL Command in Debug Mode

### Install

```bash
go get -v github.com/lxzan/hasaki
```

### Usage

#### Get

```go
// GET https://api.example.com/search
// Send get request with path parameters. Turn on data compression.

resp := hasaki.
    Get("https://api.example.com/%s", "search").
    SetHeader("Accept-Encoding", "gzip, deflate").
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

resp := hasaki.
    Post("https://api.example.com/search").
    SetEncoder(hasaki.FormEncoder).
    Send(url.Values{
        "q":    []string{"hasaki"},
        "page": []string{"1"},
    })
```

#### Stream

```go
// POST https://api.example.com/upload
// Send a put request, using a byte stream

var reader io.Reader
encoder := hasaki.NewStreamEncoder(hasaki.MimeStream)
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

Very useful middleware, you can use it to do something before and after the request is sent.

The middleware is a function, it receives a context and a request or response object, and returns a context and an error.

Under code is a simple middleware example , record the request latency.

```go
// You can use the before and after middleware to do something before and after the request is sent

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
```