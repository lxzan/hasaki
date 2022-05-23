package hasaki

import "time"

const (
	DefaultTimeout             = 10 * time.Second
	DefaultMaxIdleConnsPerHost = 32
)

const (
	Method_GET    = "GET"
	Method_POST   = "POST"
	Method_PUT    = "PUT"
	Method_DELETE = "DELETE"
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