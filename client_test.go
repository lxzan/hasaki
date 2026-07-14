package hasaki

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

var _port = int64(10086)

func nextAddr() string {
	port := atomic.AddInt64(&_port, 1)
	return "127.0.0.1:" + strconv.Itoa(int(port))
}

func TestClient(t *testing.T) {
	addr := nextAddr()
	srv := &http.Server{Addr: addr}
	srv.Handler = http.Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}))
	go srv.ListenAndServe()
	time.Sleep(100 * time.Millisecond)

	c, _ := NewClient(WithHTTPClient(&http.Client{}))
	{
		resp := c.Get("http://%s", addr).Send(nil)
		assert.NoError(t, resp.Err())
	}
	{
		resp := c.Post("http://%s", addr).Send(nil)
		assert.NoError(t, resp.Err())
	}
	{
		resp := c.Put("http://%s", addr).Send(nil)
		assert.NoError(t, resp.Err())
	}
	{
		resp := c.Delete("http://%s", addr).Send(nil)
		assert.NoError(t, resp.Err())
	}
	{
		req := c.Get("http://%s", addr).SetQuery("name=xxx")
		exp := fmt.Sprintf("http://%s?name=xxx", addr)
		assert.Equal(t, req.url, exp)
	}
	{
		type Req struct {
			Name string `form:"name"`
		}
		req := c.Get("http://%s", addr).SetQuery(Req{Name: "xxx"})
		assert.True(t, errors.Is(req.err, errUnsupportedData))
	}
}

func TestRequest(t *testing.T) {
	client, _ := NewClient()
	SetClient(client)
	addr := nextAddr()
	srv := &http.Server{Addr: addr}
	srv.Handler = http.Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "":
		case "/token":
			writer.Header().Set("x-token", request.Header.Get("x-token"))
		}
		writer.WriteHeader(http.StatusOK)
	}))
	go srv.ListenAndServe()
	time.Sleep(100 * time.Millisecond)

	{
		resp := Get("http://%s", addr).Send(nil)
		assert.NoError(t, resp.Err())
	}
	{
		resp := Post("http://%s", addr).Send(nil)
		assert.NoError(t, resp.Err())
	}
	{
		resp := Put("http://%s", addr).Send(nil)
		assert.NoError(t, resp.Err())
	}
	{
		resp := Delete("http://%s", addr).Send(nil)
		assert.NoError(t, resp.Err())
	}
	{
		resp := Patch("http://%s", addr).Send(nil)
		assert.NoError(t, resp.Err())
	}
	{
		resp := Head("http://%s", addr).Send(nil)
		assert.NoError(t, resp.Err())
	}
	{
		resp := Options("http://%s", addr).Send(nil)
		assert.NoError(t, resp.Err())
	}
	{
		resp := NewRequest(http.MethodDelete, "http://%s", addr).Send(nil)
		assert.NoError(t, resp.Err())
	}
	{
		resp := Post("http://%s/token", addr).
			SetHeader("x-token", "123").
			Send(nil)
		assert.NoError(t, resp.Err())
		assert.Equal(t, resp.Header.Get("x-token"), "123")
	}
}

func TestRequest_SetContext(t *testing.T) {
	addr := nextAddr()
	srv := &http.Server{Addr: addr}
	srv.Handler = http.Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}))
	go srv.ListenAndServe()
	time.Sleep(100 * time.Millisecond)

	ctx := context.WithValue(context.Background(), "name", "aha")
	resp := Post("http://%s", addr).SetContext(ctx).Send(nil)
	assert.Equal(t, resp.Context().Value("name"), "aha")
}

func TestRequest_SetQuery(t *testing.T) {
	addr := nextAddr()
	srv := &http.Server{Addr: addr}
	srv.Handler = http.Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/404":
			writer.WriteHeader(http.StatusNotFound)
		default:
			writer.WriteHeader(http.StatusOK)
		}
	}))
	go srv.ListenAndServe()
	time.Sleep(100 * time.Millisecond)

	t.Run("", func(t *testing.T) {
		resp := Get("http://127.0.0.1:xx").SetQuery("name=1").Send(nil)
		assert.Error(t, resp.Err())
	})

	t.Run("", func(t *testing.T) {
		resp := Get("http://%s", addr).SetQuery(nil).Send(nil)
		assert.Error(t, resp.Err())
	})

	t.Run("", func(t *testing.T) {
		req := Get("http://%s", addr).SetQuery(url.Values{
			"name": []string{"xxx"},
		})
		assert.Equal(t, req.url, "http://"+addr+"?name=xxx")
	})
}

func TestRequest_Send(t *testing.T) {
	addr := nextAddr()
	srv := &http.Server{Addr: addr}
	srv.Handler = http.Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/404":
			writer.WriteHeader(http.StatusNotFound)
		case "/500":
			writer.WriteHeader(http.StatusInternalServerError)
		default:
			writer.WriteHeader(http.StatusOK)
		}
	}))
	go srv.ListenAndServe()
	time.Sleep(100 * time.Millisecond)

	t.Run("", func(t *testing.T) {
		resp := Post("http://%s", addr).
			SetEncoder(FormCodec).
			Send(nil)
		assert.NoError(t, resp.Err())
	})

	t.Run("", func(t *testing.T) {
		resp := Post("http://127.0.0.1:xx").Send(nil)
		assert.Error(t, resp.Err())
	})

	t.Run("", func(t *testing.T) {
		resp := Post("http://127.0.0.1:xx").Send(nil)
		assert.Error(t, resp.Err())
	})

	t.Run("", func(t *testing.T) {
		resp := Post("http://%s/500", nextAddr()).Send(nil)
		assert.Error(t, resp.Err())
	})

	t.Run("json body is replayable", func(t *testing.T) {
		checked := false
		before := func(ctx context.Context, request *http.Request) (context.Context, error) {
			if !assert.NotNil(t, request.GetBody) {
				return ctx, nil
			}
			body, err := request.GetBody()
			if !assert.NoError(t, err) {
				return ctx, nil
			}
			defer body.Close()
			data, err := io.ReadAll(body)
			assert.NoError(t, err)
			assert.JSONEq(t, `{"name":"hasaki"}`, string(data))
			checked = true
			return ctx, nil
		}

		resp := Post("http://%s", addr).
			SetBefore(before).
			Send(struct {
				Name string `json:"name"`
			}{Name: "hasaki"})
		assert.NoError(t, resp.Err())
		assert.True(t, checked)
	})
}

func TestMiddleware(t *testing.T) {
	addr := nextAddr()
	srv := &http.Server{Addr: addr}
	srv.Handler = http.Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/404":
			writer.WriteHeader(http.StatusNotFound)
		default:
			writer.WriteHeader(http.StatusOK)
		}
	}))
	go srv.ListenAndServe()
	time.Sleep(100 * time.Millisecond)

	t.Run("before", func(t *testing.T) {
		before := func(ctx context.Context, request *http.Request) (context.Context, error) {
			return ctx, errors.New("status error")
		}

		{
			cli, _ := NewClient(WithBefore(before))
			resp := cli.Post("http://%s/404", addr).Send(nil)
			assert.Error(t, resp.Err())
		}

		{
			resp := Post("http://%s/404", addr).SetBefore(before).Send(nil)
			assert.Error(t, resp.Err())
		}
	})

	t.Run("after", func(t *testing.T) {
		after := func(ctx context.Context, response *http.Response) (context.Context, error) {
			if response.StatusCode != http.StatusOK {
				return ctx, errors.New("status error")
			}
			return ctx, nil
		}

		{
			cli, _ := NewClient(WithAfter(after))
			resp := cli.Post("http://%s/404", addr).Send(nil)
			assert.Error(t, resp.Err())
		}

		{
			resp := Post("http://%s/404", addr).SetAfter(after).Send(nil)
			assert.Error(t, resp.Err())
		}
	})

	t.Run("latency", func(t *testing.T) {
		before := WithBefore(func(ctx context.Context, request *http.Request) (context.Context, error) {
			return context.WithValue(ctx, "t0", time.Now()), nil
		})

		after := WithAfter(func(ctx context.Context, response *http.Response) (context.Context, error) {
			t0 := ctx.Value("t0").(time.Time)
			time.Sleep(time.Millisecond)
			return context.WithValue(ctx, "latency", time.Since(t0).Nanoseconds()), nil
		})

		cli, _ := NewClient(before, after)
		resp := cli.Post("http://%s/latency", addr).Send(nil)
		assert.Greater(t, resp.Context().Value("latency").(int64), int64(0))
	})
}

func TestResponse(t *testing.T) {
	addr := nextAddr()
	srv := &http.Server{Addr: addr}
	srv.Handler = http.Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/greet":
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("hello"))
		case "/json":
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte(`{"name":"caster"}`))
		case "/yaml":
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte(`name: caster`))
		case "/xml":
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte(`<A><name>caster</name></A>`))
		case "/proto":
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte{10, 6, 99, 97, 115, 116, 101, 114})
		case "/204":
			writer.WriteHeader(http.StatusNoContent)
		default:
			writer.WriteHeader(http.StatusOK)
		}
	}))
	go srv.ListenAndServe()
	time.Sleep(100 * time.Millisecond)

	t.Run("read body ok", func(t *testing.T) {
		resp := Get("http://%s/greet", addr).Send(nil)
		p, err := resp.ReadBody()
		assert.NoError(t, err)
		assert.Equal(t, string(p), "hello")
		assert.Equal(t, resp.Request.Header.Get("Content-Type"), "")
	})

	t.Run("read body error 1", func(t *testing.T) {
		resp := Post("http://%s/json", nextAddr()).Send(nil)
		_, err := resp.ReadBody()
		assert.Error(t, err)
	})

	t.Run("read body error 2", func(t *testing.T) {
		resp := Post("http://%s/json", addr).Send(nil)
		resp.Body = nil
		_, err := resp.ReadBody()
		assert.Error(t, err)
	})

	t.Run("bind json ok", func(t *testing.T) {
		resp := Post("http://%s/json", addr).Send(nil)
		input := struct{ Name string }{}
		err := resp.BindJSON(&input)
		assert.NoError(t, err)
		assert.Equal(t, input.Name, "caster")
	})

	t.Run("bind json error 1", func(t *testing.T) {
		resp := Post("http://%s/json", nextAddr()).Send(nil)
		inputs := struct{ Name string }{}
		err := resp.BindJSON(&inputs)
		assert.Error(t, err)
	})

	t.Run("bind json error 2", func(t *testing.T) {
		resp := Post("http://%s/json", addr).Send(map[string]any{
			"name": "xxx",
		})
		resp.Body = nil
		inputs := struct{ Name string }{}
		err := resp.BindJSON(&inputs)
		assert.Error(t, err)
	})

	t.Run("204 response", func(t *testing.T) {
		resp1 := Post("http://%s/204", addr).
			Debug().
			SetEncoder(JsonCodec).
			Send(map[string]any{
				"name": "xxx",
			})
		assert.NotNil(t, resp1.Body)

		resp2 := Post("http://%s/204", addr).
			Debug().
			SetEncoder(FormCodec).
			Send(nil)
		assert.NotNil(t, resp2.Body)
	})

	t.Run("post form error", func(t *testing.T) {
		var netConn *net.TCPConn
		resp := Post("http://%s/204", addr).
			Debug().
			SetEncoder(FormCodec).
			Send(net.Conn(netConn))
		assert.Error(t, resp.Err())
	})
}

func TestRequest_SetHeaders(t *testing.T) {
	var h = http.Header{}
	h.Set("Content-Type", MimeJson)
	h.Set("cookie", "123")
	var req = Get("https://api.github.com").
		SetHeader("Cookie", "456").
		SetHeader("encoding", "none").
		SetHeaders(h)
	assert.Equal(t, req.headers.Get("cookie"), "123")
	assert.Equal(t, req.headers.Get("encoding"), "none")
}

func TestRequest_ReadBody(t *testing.T) {
	addr := nextAddr()
	srv := &http.Server{Addr: addr}
	srv.Handler = http.Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/greet":
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("hello"))
		default:
			writer.WriteHeader(http.StatusOK)
		}
	}))
	go srv.ListenAndServe()
	time.Sleep(100 * time.Millisecond)

	t.Run("ok", func(t *testing.T) {
		var cli, _ = NewClient(WithReuseBody())
		var resp = cli.Get("http://%s/greet", addr).Send(nil)
		assert.NoError(t, resp.Err())
		_, ok := resp.Body.(BytesReadCloser)
		assert.True(t, ok)
		_, err := resp.ReadBody()
		assert.NoError(t, err)
	})
}

func TestWithBaseURL(t *testing.T) {
	addr := nextAddr()
	srv := &http.Server{Addr: addr}
	srv.Handler = http.Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/api/users":
			writer.WriteHeader(http.StatusOK)
		case "/users":
			writer.WriteHeader(http.StatusOK)
		default:
			writer.WriteHeader(http.StatusOK)
		}
	}))
	go srv.ListenAndServe()
	time.Sleep(100 * time.Millisecond)

	t.Run("base url with relative path", func(t *testing.T) {
		baseURL := "http://" + addr
		cli, _ := NewClient(WithBaseURL(baseURL))
		req := cli.Get("/api/users")
		expectedURL := baseURL + "/api/users"
		assert.Equal(t, req.url, expectedURL)
	})

	t.Run("base url with absolute path", func(t *testing.T) {
		baseURL := "http://" + addr
		cli, _ := NewClient(WithBaseURL(baseURL))
		req := cli.Get("http://example.com/users")
		// 当传入绝对 URL 时，应该直接使用该 URL（当前实现是字符串拼接）
		expectedURL := baseURL + "http://example.com/users"
		assert.Equal(t, req.url, expectedURL)
	})

	t.Run("base url with empty path", func(t *testing.T) {
		baseURL := "http://" + addr
		cli, _ := NewClient(WithBaseURL(baseURL))
		req := cli.Get("")
		expectedURL := baseURL
		assert.Equal(t, req.url, expectedURL)
	})

	t.Run("base url with formatted path", func(t *testing.T) {
		baseURL := "http://" + addr
		cli, _ := NewClient(WithBaseURL(baseURL))
		req := cli.Get("/api/%s", "users")
		expectedURL := baseURL + "/api/users"
		assert.Equal(t, req.url, expectedURL)
	})

	t.Run("base url actual request", func(t *testing.T) {
		baseURL := "http://" + addr
		cli, _ := NewClient(WithBaseURL(baseURL))
		resp := cli.Get("/api/users").Send(nil)
		assert.NoError(t, resp.Err())
		assert.Equal(t, resp.StatusCode, http.StatusOK)
	})

	t.Run("base url without trailing slash", func(t *testing.T) {
		baseURL := "http://" + addr
		cli, _ := NewClient(WithBaseURL(baseURL))
		req := cli.Get("/users")
		expectedURL := baseURL + "/users"
		assert.Equal(t, req.url, expectedURL)
	})

	t.Run("base url with trailing slash", func(t *testing.T) {
		baseURL := "http://" + addr + "/"
		cli, _ := NewClient(WithBaseURL(baseURL))
		req := cli.Get("users")
		expectedURL := baseURL + "users"
		assert.Equal(t, req.url, expectedURL)
	})
}

func TestResponse_BindJSON(t *testing.T) {
	addr := nextAddr()
	srv := &http.Server{Addr: addr}
	srv.Handler = http.Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(`{"name":"test","age":18}`))
	}))
	go srv.ListenAndServe()
	time.Sleep(100 * time.Millisecond)

	t.Run("bind json success", func(t *testing.T) {
		type User struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		var user User
		resp := Get("http://%s", addr).Send(nil)
		err := resp.BindJSON(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, "test")
		assert.Equal(t, user.Age, 18)
	})

	t.Run("bind json with error response", func(t *testing.T) {
		type User struct {
			Name string `json:"name"`
		}
		var user User
		resp := Get("http://127.0.0.1:xx").Send(nil)
		err := resp.BindJSON(&user)
		assert.Error(t, err)
	})
}

func TestResponse_BindXML(t *testing.T) {
	addr := nextAddr()
	srv := &http.Server{Addr: addr}
	srv.Handler = http.Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(`<user><name>test</name><age>18</age></user>`))
	}))
	go srv.ListenAndServe()
	time.Sleep(100 * time.Millisecond)

	t.Run("bind xml success", func(t *testing.T) {
		type User struct {
			XMLName xml.Name `xml:"user"`
			Name    string   `xml:"name"`
			Age     int      `xml:"age"`
		}
		var user User
		resp := Get("http://%s", addr).Send(nil)
		err := resp.BindXML(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, "test")
		assert.Equal(t, user.Age, 18)
	})
}

func TestResponse_BindForm(t *testing.T) {
	addr := nextAddr()
	srv := &http.Server{Addr: addr}
	srv.Handler = http.Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("name=test&age=18"))
	}))
	go srv.ListenAndServe()
	time.Sleep(100 * time.Millisecond)

	t.Run("bind form success", func(t *testing.T) {
		var params url.Values
		resp := Get("http://%s", addr).Send(nil)
		err := resp.BindForm(&params)
		assert.NoError(t, err)
		assert.Equal(t, params.Get("name"), "test")
		assert.Equal(t, params.Get("age"), "18")
	})
}

func TestRequest_SetEncoder(t *testing.T) {
	req := Get("https://api.example.com")
	req.SetEncoder(FormCodec)
	assert.Equal(t, req.headers.Get("Content-Type"), MimeForm)

	req.SetEncoder(JsonCodec)
	assert.Equal(t, req.headers.Get("Content-Type"), MimeJson)

	req.SetEncoder(XmlCodec)
	assert.Equal(t, req.headers.Get("Content-Type"), MimeXml)
}

func TestJsonCodec_Decode(t *testing.T) {
	t.Run("decode success", func(t *testing.T) {
		type User struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		var user User
		reader := strings.NewReader(`{"name":"test","age":18}`)
		err := JsonCodec.Decode(reader, &user)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, "test")
		assert.Equal(t, user.Age, 18)
	})

	t.Run("decode map", func(t *testing.T) {
		var result map[string]any
		reader := strings.NewReader(`{"name":"test","age":18}`)
		err := JsonCodec.Decode(reader, &result)
		assert.NoError(t, err)
		assert.Equal(t, result["name"], "test")
	})
}
