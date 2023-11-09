package hasaki

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
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
		exp := fmt.Sprintf("http://%s?name=xxx", addr)
		assert.Equal(t, req.url, exp)
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

func TestRequest_Header(t *testing.T) {
	{
		var req = Post("http://%s", nextAddr()).SetEncoder(FormEncoder)
		var typ = req.Header().Get("Content-Type")
		assert.Equal(t, typ, MimeForm)
	}
	{
		var req = Get("http://%s", nextAddr())
		var typ = req.Header().Get("Content-Type")
		assert.Equal(t, typ, MimeJson)
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

	t.Run("", func(t *testing.T) {
		type Req struct {
			Name string `form:"name"`
		}
		req := Get("http://%s", addr).SetQuery(Req{Name: "xxx"})
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
			SetEncoder(FormEncoder).
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

	t.Run("bind yaml ok", func(t *testing.T) {
		resp := Post("http://%s/yaml", addr).Send(nil)
		input := struct{ Name string }{}
		err := resp.BindYAML(&input)
		assert.NoError(t, err)
		assert.Equal(t, input.Name, "caster")
	})

	t.Run("bind yaml error 1", func(t *testing.T) {
		resp := Post("http://%s/yaml", nextAddr()).Send(nil)
		inputs := struct{ Name string }{}
		err := resp.BindYAML(&inputs)
		assert.Error(t, err)
	})

	t.Run("bind yaml error 2", func(t *testing.T) {
		resp := Post("http://%s/yaml", addr).Send(map[string]any{
			"name": "xxx",
		})
		resp.Body = nil
		inputs := struct{ Name string }{}
		err := resp.BindYAML(&inputs)
		assert.Error(t, err)
	})

	t.Run("bind xml ok", func(t *testing.T) {
		resp := Post("http://%s/xml", addr).Send(nil)
		input := struct {
			Name string `xml:"name"`
		}{}
		err := resp.BindXML(&input)
		assert.NoError(t, err)
		assert.Equal(t, input.Name, "caster")
	})

	t.Run("bind xml error 1", func(t *testing.T) {
		resp := Post("http://%s/xml", nextAddr()).Send(nil)
		inputs := struct {
			Name string `xml:"name"`
		}{}
		err := resp.BindXML(&inputs)
		assert.Error(t, err)
	})

	t.Run("bind xml error 2", func(t *testing.T) {
		resp := Post("http://%s/xml", addr).Send(map[string]any{
			"name": "xxx",
		})
		resp.Body = nil
		inputs := struct {
			Name string `xml:"name"`
		}{}
		err := resp.BindXML(&inputs)
		assert.Error(t, err)
	})

	t.Run("bind proto ok", func(t *testing.T) {
		resp := Post("http://%s/proto", addr).Send(nil)
		input := Test{}
		err := resp.BindPROTO(&input)
		assert.NoError(t, err)
		assert.Equal(t, input.Name, "caster")
	})

	t.Run("bind proto error 1", func(t *testing.T) {
		resp := Post("http://%s/proto", nextAddr()).Send(nil)
		inputs := Test{}
		err := resp.BindPROTO(&inputs)
		assert.Error(t, err)
	})

	t.Run("bind proto error 2", func(t *testing.T) {
		resp := Post("http://%s/proto", addr).Send(map[string]any{
			"name": "xxx",
		})
		resp.Body = nil
		inputs := Test{}
		err := resp.BindPROTO(&inputs)
		assert.Error(t, err)
	})

	t.Run("204 response", func(t *testing.T) {
		resp1 := Post("http://%s/204", addr).
			Debug().
			SetEncoder(JsonEncoder).
			Send(map[string]any{
				"name": "xxx",
			})
		assert.NotNil(t, resp1.Body)

		resp2 := Post("http://%s/204", addr).
			Debug().
			SetEncoder(FormEncoder).
			Send(nil)
		assert.NotNil(t, resp2.Body)
	})

	t.Run("post form error", func(t *testing.T) {
		var netConn *net.TCPConn
		resp := Post("http://%s/204", addr).
			Debug().
			SetEncoder(FormEncoder).
			Send(net.Conn(netConn))
		assert.Error(t, resp.Err())
	})
}

func TestResponse_Latency(t *testing.T) {
	addr := nextAddr()
	srv := &http.Server{Addr: addr}
	srv.Handler = http.Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(100 * time.Millisecond)
		writer.WriteHeader(http.StatusOK)
	}))
	go srv.ListenAndServe()
	time.Sleep(100 * time.Millisecond)

	resp := Get("http://%s", addr).Send(nil)
	assert.Greater(t, resp.Latency(), int64(0))
}
