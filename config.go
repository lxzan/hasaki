package hasaki

import (
	"context"
	"github.com/pkg/errors"
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

	defaultErrorChecker ErrorChecker = func(resp *http.Response) error {
		if resp.StatusCode != 200 {
			return errors.New("unexpected status_code")
		}
		return nil
	}

	defaultBeforeFunc = func(ctx context.Context, request *http.Request) (context.Context, error) {
		return ctx, nil
	}

	defaultAfterFunc = func(ctx context.Context, request *http.Response) (context.Context, error) {
		return ctx, nil
	}
)

func SetGlobalClient(client *http.Client) {
	defaultHTTPClient = client
}

func SetGlobalErrorChecker(fn ErrorChecker) {
	defaultErrorChecker = fn
}

func SetBefore(fn func(ctx context.Context, request *http.Request) (context.Context, error)) {
	defaultBeforeFunc = fn
}

func SetAfter(fn func(ctx context.Context, request *http.Response) (context.Context, error)) {
	defaultAfterFunc = fn
}
