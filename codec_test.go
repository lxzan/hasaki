package hasaki

import (
	"bytes"
	"net"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForm_encoder_Encode(t *testing.T) {
	_, err1 := FormEncoder.Encode(struct {
		Name string
	}{Name: "caster"})
	assert.NoError(t, err1)

	_, err2 := FormEncoder.Encode(url.Values{
		"name": []string{"caster"},
	})
	assert.NoError(t, err2)

	_, err3 := FormEncoder.Encode(nil)
	assert.NoError(t, err3)

	var netConn *net.TCPConn
	_, err4 := FormEncoder.Encode(net.Conn(netConn))
	assert.Error(t, err4)
}

func TestStreamEncoder(t *testing.T) {
	encoder := NewStreamEncoder("text/plain")
	assert.Equal(t, encoder.ContentType(), "text/plain")

	_, err1 := encoder.Encode("aha")
	assert.NoError(t, err1)

	_, err2 := encoder.Encode([]byte("aha"))
	assert.NoError(t, err2)

	_, err3 := encoder.Encode(bytes.NewBufferString("oh"))
	assert.NoError(t, err3)

	_, err4 := encoder.Encode(123)
	assert.Error(t, err4)
}

func TestXmlDecode(t *testing.T) {
	var v = struct {
	}{}
	var r = bytes.NewBufferString("")
	var err = XmlDecode(r, &v)
	assert.Error(t, err)
}
