package internal

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/bytebufferpool"
	"testing"
)

func TestGetBuffer(t *testing.T) {
	assert.GreaterOrEqual(t, GetBuffer().Cap(), BufferSize)
	assert.Equal(t, GetBuffer().Len(), 0)
}

func TestPutBuffer(t *testing.T) {
	b0 := bytes.NewBuffer(make([]byte, 100))
	PutBuffer(b0)
	GetBuffer()
}

func TestCloserWrapper(t *testing.T) {
	var cw = &CloserWrapper{
		B: bytebufferpool.Get(),
		R: bytes.NewReader(nil),
	}
	cw.Read(make([]byte, 10))
	cw.Close()
	assert.Nil(t, cw.B)
	assert.Nil(t, cw.R)
}
