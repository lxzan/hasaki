package hasaki

import "time"

type Options struct {
	timeOut            time.Duration // io timeout
	proxyURL           string        // set http proxy
	insecureSkipVerify bool          // skip verify certificate
}

func (this *Options) SetTimeOut(d time.Duration) *Options {
	this.timeOut = d
	return this
}

func (this *Options) SetProxyURL(addr string) *Options {
	this.proxyURL = addr
	return this
}

func (this *Options) SetInsecureSkipVerify(skip bool) *Options {
	this.insecureSkipVerify = skip
	return this
}
