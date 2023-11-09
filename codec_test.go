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

func TestJson_encoder_Encode(t *testing.T) {
	_, err1 := JsonEncoder.Encode(struct {
		Name string
	}{Name: "caster"})
	assert.NoError(t, err1)

	_, err2 := JsonEncoder.Encode(map[string]interface{}{
		"name": "caster",
	})
	assert.NoError(t, err2)

	_, err3 := JsonEncoder.Encode(nil)
	assert.NoError(t, err3)
}

func TestYaml_encoder_Encode(t *testing.T) {
	_, err1 := YamlEncoder.Encode(struct {
		Name string
	}{Name: "caster"})
	assert.NoError(t, err1)

	_, err2 := YamlEncoder.Encode(map[string]interface{}{
		"name": "caster",
	})
	assert.NoError(t, err2)

	_, err3 := YamlEncoder.Encode(nil)
	assert.NoError(t, err3)
}

func TestXml_encoder_Encode(t *testing.T) {
	type A struct {
		Name string `xml:"name"`
	}

	_, err1 := XmlEncoder.Encode(&A{Name: "caster"})
	assert.NoError(t, err1)

	_, err2 := XmlEncoder.Encode(map[string]interface{}{
		"name": "caster",
	})
	assert.Error(t, err2)

	_, err3 := XmlEncoder.Encode(nil)
	assert.NoError(t, err3)
}

func TestProto_encoder_Encode(t *testing.T) {
	_, err1 := ProtoEncoder.Encode(&Test{Name: "caster"})
	assert.NoError(t, err1)

	_, err2 := ProtoEncoder.Encode(&map[string]interface{}{
		"name": "caster",
	})
	assert.Error(t, err2)

	_, err3 := ProtoEncoder.Encode(nil)
	assert.NoError(t, err3)
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
