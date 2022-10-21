package hasaki

import (
	"github.com/pkg/errors"
	"net/http"
	"time"
)

var defaultHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
	},
}

func SetGlobalClient(client *http.Client) {
	defaultHTTPClient = client
}

var defaultErrorChecker ErrorChecker = func(resp *http.Response) error {
	if resp.StatusCode != 200 {
		return errors.New("unexpected status_code")
	}
	return nil
}

func SetGlobalErrorChecker(fn ErrorChecker) {
	defaultErrorChecker = fn
}
