package hasaki

import (
	"github.com/pkg/errors"
	"net/http"
	"time"
)

var defaultHTTPClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 32,
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
