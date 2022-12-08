package hasaki

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

const testURL = "https://api.github.com"

func TestWithTimeout(t *testing.T) {
	var value = 5 * time.Second
	c, _ := NewClient(WithTimeout(value))
	assert.Equal(t, c.cli.Timeout, value)
}

func TestWithMaxIdleConnsPerHost(t *testing.T) {
	var value = 128
	c, _ := NewClient(WithMaxIdleConnsPerHost(value))
	assert.Equal(t, c.cli.Transport.(*http.Transport).MaxIdleConnsPerHost, value)
}

func TestWithInsecureSkipVerify(t *testing.T) {
	var value = true
	c, _ := NewClient(WithInsecureSkipVerify(value))
	assert.Equal(t, c.cli.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify, value)
}

func TestWithProxy(t *testing.T) {
	var value = "socks5://127.0.0.1:1080"
	c, _ := NewClient(WithProxy(value))
	assert.NotNil(t, c.cli.Transport.(*http.Transport).Proxy)
}

func TestWithTransport(t *testing.T) {
	{
		var value = 128
		var transport = &http.Transport{MaxIdleConns: value}
		c, _ := NewClient(WithTransport(transport), WithInsecureSkipVerify(true))
		assert.Equal(t, c.cli.Transport.(*http.Transport).MaxIdleConns, 128)
	}

	{
		var value = 128
		var transport = &http.Transport{MaxIdleConns: value}
		c, _ := NewClient(WithInsecureSkipVerify(true), WithTransport(transport))
		assert.Equal(t, c.cli.Transport.(*http.Transport).MaxIdleConns, 128)
	}
}

func TestNewClient(t *testing.T) {
	cli, _ := NewClient()
	resp := cli.Get(testURL).Send(nil)
	assert.NoError(t, resp.Err())
}

func TestWithBefore(t *testing.T) {
	cli, _ := NewClient(WithBefore(func(ctx context.Context, request *http.Request) (context.Context, error) {
		ctx = context.WithValue(ctx, "k", "v")
		return ctx, nil
	}))
	result := cli.Get(testURL).Send(nil)
	val := result.Context().Value("k")
	assert.Equal(t, "v", val)
}

func TestWithAfter(t *testing.T) {
	cli, _ := NewClient(WithAfter(func(ctx context.Context, response *http.Response) (context.Context, error) {
		ctx = context.WithValue(ctx, "k", "v")
		return ctx, nil
	}))
	result := cli.Get(testURL).Send(nil)
	val := result.Context().Value("k")
	assert.Equal(t, "v", val)
}
