package internal

import (
	"bytes"
	"github.com/valyala/bytebufferpool"
	"io"
	"sync"
)

// BufferSize 缓冲区大小
// Buffer size
const BufferSize = 4 * 1024

var _pool = &sync.Pool{New: func() any {
	return bytes.NewBuffer(make([]byte, 0, BufferSize))
}}

// GetBuffer 从池中获取缓冲区
// Get buffer from pool
func GetBuffer() *bytes.Buffer {
	b := _pool.Get().(*bytes.Buffer)
	if b.Cap() < BufferSize {
		b.Grow(BufferSize)
	}
	b.Reset()
	return b
}

// PutBuffer 将缓冲区归还到池中
// Put buffer back to pool
func PutBuffer(b *bytes.Buffer) {
	_pool.Put(b)
}

// CloserWrapper 包装字节缓冲区和读取器，实现可关闭的读取接口
// Wrap byte buffer and reader, implement closable read interface
type CloserWrapper struct {
	B *bytebufferpool.ByteBuffer
	R io.Reader
}

// Bytes 返回字节缓冲区的内容
// Return byte buffer content
func (c *CloserWrapper) Bytes() []byte {
	return c.B.Bytes()
}

// Read 从读取器中读取数据
// Read data from reader
func (c *CloserWrapper) Read(p []byte) (n int, err error) {
	return c.R.Read(p)
}

// Close 关闭包装器并释放资源
// Close wrapper and release resources
func (c *CloserWrapper) Close() error {
	// 避免重复关闭, 引发panic
	if c.B != nil {
		bytebufferpool.Put(c.B)
		c.B, c.R = nil, nil
	}
	return nil
}
