package hasaki

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonDecode(t *testing.T) {
	var v = struct {
	}{}
	var r = bytes.NewBufferString("")
	var err = JSONDecode(r, &v)
	assert.Error(t, err)
}

func TestXmlDecode(t *testing.T) {
	var v = struct {
	}{}
	var r = bytes.NewBufferString("")
	var err = XMLDecode(r, &v)
	assert.Error(t, err)
}
