package hasaki

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const testURL = "https://api.github.com"

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
