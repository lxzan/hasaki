package hasaki

import "time"

const DefaultTimeout = 10 * time.Second

type Options struct {
	TimeOut            time.Duration
	ProxyURL           string
	InsecureSkipVerify bool // skip verify certificate
}

func (this *Options) SetTimeOut(d time.Duration) *Options {
	this.TimeOut = d
	return this
}

func (this *Options) SetProxyURL(addr string) *Options {
	this.ProxyURL = addr
	return this
}

func (this *Options) SetInsecureSkipVerify(skip bool) *Options {
	this.InsecureSkipVerify = skip
	return this
}
