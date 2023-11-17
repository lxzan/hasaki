package internal

import (
	"bytes"
	"github.com/valyala/bytebufferpool"
	"io"
	"sync"
)

const BufferSize = 4 * 1024

var _pool = &sync.Pool{New: func() any {
	return bytes.NewBuffer(make([]byte, 0, BufferSize))
}}

func GetBuffer() *bytes.Buffer {
	b := _pool.Get().(*bytes.Buffer)
	if b.Cap() < BufferSize {
		b.Grow(BufferSize)
	}
	b.Reset()
	return b
}

func PutBuffer(b *bytes.Buffer) {
	_pool.Put(b)
}

type CloserWrapper struct {
	B *bytebufferpool.ByteBuffer
	R io.Reader
}

func (c *CloserWrapper) Bytes() []byte {
	return c.B.Bytes()
}

func (c *CloserWrapper) Read(p []byte) (n int, err error) {
	return c.R.Read(p)
}

func (c *CloserWrapper) Close() error {
	c.B.Reset()
	bytebufferpool.Put(c.B)
	c.B, c.R = nil, nil
	return nil
}
