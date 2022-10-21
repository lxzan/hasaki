package hasaki

import (
	"github.com/pkg/errors"
	"time"
)

const (
	DefaultTimeout             = 10 * time.Second
	DefaultMaxIdleConnsPerHost = 32
)

type ContentType string

func (c ContentType) String() string {
	return string(c)
}

const (
	ContentType_JSON   ContentType = "application/json;charset=utf-8"
	ContentType_FORM   ContentType = "application/x-www-form-urlencoded"
	ContentType_STREAM ContentType = "application/octet-stream"
	ContentType_JPEG   ContentType = "image/jpeg"
	ContentType_GIF    ContentType = "image/gif"
	ContentType_PNG    ContentType = "image/png"
	ContentType_MP4    ContentType = "video/mpeg4"
)

var (
	ErrDataNotSupported = errors.New("data type is not supported")
)
