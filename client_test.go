package hasaki

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

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
	resp := cli.Get("https://api.github.com/").Send(nil)
	assert.NoError(t, resp.Err())
}
