package hasaki

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

var (
	defaultHTTPClient = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 32,
		},
	}

	defaultBeforeFunc = func(ctx context.Context, request *http.Request) (context.Context, error) {
		return ctx, nil
	}

	defaultAfterFunc = func(ctx context.Context, response *http.Response) (context.Context, error) {
		if response.StatusCode != http.StatusOK {
			return ctx, ErrUnexpectedStatusCode
		}
		return ctx, nil
	}
)

func SetGlobalClient(client *http.Client) {
	defaultHTTPClient = client
}

func SetBefore(fn func(ctx context.Context, request *http.Request) (context.Context, error)) {
	defaultBeforeFunc = fn
}

func SetAfter(fn func(ctx context.Context, response *http.Response) (context.Context, error)) {
	defaultAfterFunc = fn
}

func NewReadCloser(body io.ReadCloser) (*ReadCloser, error) {
	var rc = &ReadCloser{Buffer: bytes.NewBuffer(nil)}
	_, err := io.Copy(rc, body)
	_ = body.Close()
	return rc, err
}

type ReadCloser struct {
	*bytes.Buffer
}

func (r *ReadCloser) Read(p []byte) (n int, err error) {
	return r.Buffer.Read(p)
}

func (r *ReadCloser) Close() error {
	return nil
}
